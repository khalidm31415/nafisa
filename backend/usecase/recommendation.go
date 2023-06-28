package usecase

import (
	"backend/dto"
	"backend/entity"
	"backend/internal_constant"
	"backend/internal_error"
	elasticsarch_helper "backend/package_helper/elasticsearch_helper"
	"backend/package_helper/redis_helper"
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRecommendation interface {
	FindMatchingProfiles(ctx context.Context) error
	View(ctx context.Context) (*dto.RecommendedProfile, error)
	ShiftRecommendation(ctx context.Context) error
	IncrementUserAction(ctx context.Context, userID string) error
	Like(ctx context.Context) error
	Pass(ctx context.Context) error
}

type Recommendation struct {
	db           *gorm.DB
	redisCache   redis_helper.IRedisCache
	profileIndex elasticsarch_helper.IElasticsearchProfileIndex
}

func NewRecommendation(db *gorm.DB, rdb *redis.Client, profileIndex elasticsarch_helper.IElasticsearchProfileIndex) IRecommendation {
	redisCache := redis_helper.NewRedisCache(rdb, internal_constant.RecommendationUserActionCount)
	return &Recommendation{
		db:           db,
		redisCache:   redisCache,
		profileIndex: profileIndex,
	}
}

func (r *Recommendation) FindMatchingProfiles(ctx context.Context) error {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)

	err := r.db.Preload("Profile").Take(&currentUser).Error
	if err != nil {
		return err
	}

	if !currentUser.Profile.IsVerified {
		return internal_error.ErrUnverified
	}

	esProfiles, err := r.profileIndex.GetMatchingProfiles(ctx, currentUser.Profile)
	if err != nil {
		return err
	}

	matchingProfiles := []entity.UserMatchingProfile{}
	for _, p := range esProfiles {
		ID, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		profile := entity.UserMatchingProfile{
			ID:                    ID.String(),
			UserID:                currentUser.ID,
			MatchingProfileUserID: p.Source.UserID,
			Score:                 p.Score,
		}
		matchingProfiles = append(matchingProfiles, profile)
	}

	if len(matchingProfiles) == 0 {
		return internal_error.ErrNoMatchingProfileFound
	}

	if err := r.db.Create(&matchingProfiles).Error; err != nil {
		return err
	}

	return nil
}

func (r *Recommendation) View(ctx context.Context) (*dto.RecommendedProfile, error) {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)

	err := r.db.Preload("Profile").Take(&currentUser).Error
	if err != nil {
		return nil, err
	}
	if !currentUser.Profile.IsProfileComplete {
		return nil, internal_error.ErrIncompleteProfile
	}

	if !currentUser.Profile.IsVerified {
		return nil, internal_error.ErrUnverified
	}

	if currentUser.Profile.CurrentDiscussionID != nil {
		return nil, internal_error.ErrStillInDiscussion
	}

	userActionCountString, err := r.redisCache.Get(ctx, currentUser.ID)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if userActionCountString != nil {
		userActionCount, err := strconv.Atoi(*userActionCountString)
		if err != nil {
			return nil, err
		}
		if userActionCount >= internal_constant.MaxProfileRecommendationView {
			return nil, internal_error.ErrRecommendationLimitReached
		}
	}

	if currentUser.Profile.CurrentRecommendationID == nil || currentUser.Profile.CurrentRecommendationType == nil {
		if err := r.ShiftRecommendation(ctx); err != nil {
			return nil, err
		}
	}

	var recommendedUser entity.User
	if *currentUser.Profile.CurrentRecommendationType == "matching_profile" {
		userMatchingProfile := entity.UserMatchingProfile{
			ID: *currentUser.Profile.CurrentRecommendationID,
		}
		if err = r.db.Preload("MatchingProfileUser.Profile").Take(&userMatchingProfile).Error; err != nil {
			return nil, err
		}
		recommendedUser = userMatchingProfile.MatchingProfileUser
	}
	if *currentUser.Profile.CurrentRecommendationType == "liker" {
		userLiker := entity.Like{
			ID: *currentUser.Profile.CurrentRecommendationID,
		}
		if err = r.db.Preload("User.Profile").Take(&userLiker).Error; err != nil {
			return nil, err
		}
		recommendedUser = userLiker.User
	}
	recommendedProfile := dto.NewRecommendedProfile(recommendedUser)
	return &recommendedProfile, nil
}

func (r *Recommendation) ShiftRecommendation(ctx context.Context) error {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)
	err := r.db.Preload("Profile").Take(&currentUser).Error
	if err != nil {
		return err
	}

	if !currentUser.Profile.IsVerified {
		return internal_error.ErrUnverified
	}

	liker_exists := true
	var nextRecommendationID string

	// case 1: current recommendation is from matching profile, next recommendation is from liker (if exists)
	if currentUser.Profile.CurrentRecommendationType != nil && *currentUser.Profile.CurrentRecommendationType == "matching_profile" {
		where := map[string]interface{}{
			"liked_user_id": currentUser.ID,
			"is_liked_back": nil,
		}
		err := r.db.
			Model(&entity.Like{}).
			Select("id").
			Where(where).
			Order("created_at ASC").
			First(&nextRecommendationID).
			Error
		if err == gorm.ErrRecordNotFound {
			liker_exists = false
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		update := map[string]interface{}{
			"current_recommendation_id":   nextRecommendationID,
			"current_recommendation_type": "liker",
		}
		if err = r.db.Model(&currentUser.Profile).Updates(update).Error; err != nil {
			return nil
		}
	}

	// case 2: current recommendation is from liker, next recommendation is from matching profile
	if currentUser.Profile.CurrentRecommendationType == nil || *currentUser.Profile.CurrentRecommendationType == "liker" || !liker_exists {
		where := map[string]interface{}{
			"user_id": currentUser.ID,
			"action":  nil,
		}
		err := r.db.
			Model(&entity.UserMatchingProfile{}).
			Select("id").
			Where(where).
			Order("score DESC").
			First(&nextRecommendationID).
			Error
		if err == gorm.ErrRecordNotFound {
			return internal_error.ErrRecommendationNotReady
		}
		if err != nil {
			return nil
		}
		updates := map[string]interface{}{
			"current_recommendation_id":   nextRecommendationID,
			"current_recommendation_type": "matching_profile",
		}
		if err = r.db.Model(&currentUser.Profile).Updates(updates).Error; err != nil {
			return nil
		}
	}

	return nil
}

func (r *Recommendation) IncrementUserAction(ctx context.Context, userID string) error {
	userActionCountString, err := r.redisCache.Get(ctx, userID)
	if err != nil && err != redis.Nil {
		return err
	}
	userActionCount := 0
	if userActionCountString != nil {
		userActionCount, err = strconv.Atoi(*userActionCountString)
		if err != nil {
			return err
		}
	}
	now := time.Now()
	year, month, day := now.Date()
	midnight := time.Date(year, month, day+1, 0, 0, 0, 0, now.Location())
	durationUntilMidnight := time.Until(midnight)
	if err := r.redisCache.Set(ctx, userID, userActionCount+1, durationUntilMidnight); err != nil {
		return err
	}
	return nil
}

func (r *Recommendation) Like(ctx context.Context) error {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)
	err := r.db.Preload("Profile").Take(&currentUser).Error
	if err != nil {
		return err
	}

	if !currentUser.Profile.IsVerified {
		return internal_error.ErrUnverified
	}
	if currentUser.Profile.CurrentRecommendationID == nil || currentUser.Profile.CurrentRecommendationType == nil {
		return internal_error.ErrRecommendationNotReady
	}
	if currentUser.Profile.CurrentDiscussionID != nil {
		return internal_error.ErrStillInDiscussion
	}

	isMatch := false
	var recommendedUserID string

	if *currentUser.Profile.CurrentRecommendationType == "matching_profile" {
		if err := r.db.Model(&entity.UserMatchingProfile{}).Where("id = ?", currentUser.Profile.CurrentRecommendationID).Update("action", "like").Error; err != nil {
			return err
		}
		if err := r.db.Model(&entity.UserMatchingProfile{}).Where("id = ?", currentUser.Profile.CurrentRecommendationID).Select("matching_profile_user_id").Take(&recommendedUserID).Error; err != nil {
			return err
		}
		var like entity.Like
		err := r.db.Where("user_id = ? AND liked_user_id = ?", recommendedUserID, currentUser.ID).Take(&like).Error
		if err == nil {
			isMatch = true
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if *currentUser.Profile.CurrentRecommendationType == "liker" {
		isMatch = true
		if err := r.db.Model(&entity.Like{}).Where("id = ?", currentUser.Profile.CurrentRecommendationID).Update("is_liked_back", true).Error; err != nil {
			return err
		}
		if err := r.db.Model(&entity.Like{}).Where("id = ?", currentUser.Profile.CurrentRecommendationID).Select("user_id").Take(&recommendedUserID).Error; err != nil {
			return err
		}
	}

	ID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	like := entity.Like{
		ID:          ID.String(),
		UserID:      currentUser.ID,
		LikedUserID: recommendedUserID,
	}
	if err := r.db.Create(&like).Error; err != nil {
		return err
	}
	if err := r.ShiftRecommendation(ctx); err != nil {
		return err
	}
	if err := r.IncrementUserAction(ctx, currentUser.ID); err != nil {
		return err
	}

	if isMatch {
		discussionID, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		discussion := entity.Discussion{
			ID: discussionID.String(),
		}

		isUserMale := true
		if currentUser.Profile.Gender == "f" {
			isUserMale = false
		}
		if isUserMale {
			discussion.MaleUserID = currentUser.ID
			discussion.FemaleUserID = recommendedUserID
		} else {
			discussion.MaleUserID = recommendedUserID
			discussion.FemaleUserID = currentUser.ID
		}

		if err := r.db.Create(&discussion).Error; err != nil {
			return err
		}
		if err := r.db.Model(&entity.UserProfile{}).Where("user_id IN ?", []string{currentUser.ID, recommendedUserID}).Update("current_discussion_id", discussionID.String()).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Recommendation) Pass(ctx context.Context) error {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)
	err := r.db.Preload("Profile").Take(&currentUser).Error
	if err != nil {
		return err
	}

	if !currentUser.Profile.IsVerified {
		return internal_error.ErrUnverified
	}
	if currentUser.Profile.CurrentRecommendationID == nil || currentUser.Profile.CurrentRecommendationType == nil {
		return internal_error.ErrRecommendationNotReady
	}
	if currentUser.Profile.CurrentDiscussionID != nil {
		return internal_error.ErrStillInDiscussion
	}

	if *currentUser.Profile.CurrentRecommendationType == "matching_profile" {
		if err := r.db.Model(&entity.UserMatchingProfile{ID: *currentUser.Profile.CurrentRecommendationID}).Update("action", "pass").Error; err != nil {
			return err
		}
	}
	if *currentUser.Profile.CurrentRecommendationType == "liker" {
		if err := r.db.Model(&entity.Like{ID: *currentUser.Profile.CurrentRecommendationID}).Update("is_liked_back", false).Error; err != nil {
			return err
		}
	}

	if err := r.IncrementUserAction(ctx, currentUser.ID); err != nil {
		return err
	}
	if err := r.ShiftRecommendation(ctx); err != nil {
		return err
	}
	return nil
}
