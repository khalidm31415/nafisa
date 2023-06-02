package usecase

import (
	"context"
	"errors"
	"nafisah/delivery/dto"
	"nafisah/entity"
	"nafisah/internal_constant"
	"nafisah/internal_error"
	"nafisah/pkg"
	"time"

	"github.com/go-redis/redis/v8"
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
	db         *gorm.DB
	redisCache pkg.IRedisCache
}

func NewRecommendation(db *gorm.DB, rdb *redis.Client) IRecommendation {
	redisCache := pkg.NewRedisCache(rdb, internal_constant.RecommendationUserActionCount)
	return &Recommendation{
		db:         db,
		redisCache: redisCache,
	}
}

func (r *Recommendation) FindMatchingProfiles(ctx context.Context) error {
	return errors.New("not implemented")
}

func (r *Recommendation) View(ctx context.Context) (*dto.RecommendedProfile, error) {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)

	err := r.db.Preload("Profile").Take(&currentUser).Error
	if err != nil {
		return nil, err
	}

	if !currentUser.Profile.IsVerified {
		return nil, internal_error.ErrUnverified
	}

	var userActionCount int
	r.redisCache.Get(ctx, currentUser.ID, userActionCount)
	if userActionCount >= internal_constant.MaxProfileRecommendationView {
		return nil, internal_error.ErrRecommendationLimitReached
	}

	if currentUser.Profile.CurrentRecommendationID == nil || currentUser.Profile.CurrentRecommendationType == nil {
		return nil, internal_error.ErrRecommendationNotReady
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
		userLiker := entity.UserLiker{
			ID: *currentUser.Profile.CurrentRecommendationID,
		}
		if err = r.db.Preload("LikerUser.Profile").Take(&userLiker).Error; err != nil {
			return nil, err
		}
		recommendedUser = userLiker.LikerUser
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

	if currentUser.Profile.CurrentRecommendationID == nil || currentUser.Profile.CurrentRecommendationType == nil {
		return internal_error.ErrRecommendationNotReady
	}

	liker_exists := true
	var nextRecommendationID string
	// current recommendation is from matching profile, next recommendation is from liker
	if *currentUser.Profile.CurrentRecommendationType == "matching_profile" {
		where := map[string]interface{}{
			"user_id": currentUser.ID,
			"action":  nil,
		}
		err := r.db.
			Model(&entity.UserLiker{}).
			Select("user_likers.id").
			Where(where).
			Order("user_likers.created_at ASC").
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

	// current recommendation is from liker, next recommendation is from matching profile
	if *currentUser.Profile.CurrentRecommendationType == "liker" || !liker_exists {
		where := map[string]interface{}{
			"user_id": currentUser.ID,
			"action":  nil,
		}
		err := r.db.
			Model(&entity.UserMatchingProfile{}).
			Select("user_matching_profiles.matching_profile_user_id").
			Where(where).
			Order("score DESC").
			First(&nextRecommendationID).
			Error
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
	var userActionCount int
	if err := r.redisCache.Get(ctx, userID, userActionCount); err != nil {
		return err
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

	if *currentUser.Profile.CurrentRecommendationType == "matching_profile" {
		if err := r.db.Model(&entity.UserMatchingProfile{ID: *currentUser.Profile.CurrentRecommendationID}).Update("action", "like").Error; err != nil {
			return err
		}
		if err := r.ShiftRecommendation(ctx); err != nil {
			return err
		}
		if err := r.IncrementUserAction(ctx, currentUser.ID); err != nil {
			return err
		}
	}

	// it's a match!
	if *currentUser.Profile.CurrentRecommendationType == "liker" {
		var likerUserID string
		if err := r.db.Model(&entity.UserLiker{ID: *currentUser.Profile.CurrentRecommendationID}).Update("action", "like").Error; err != nil {
			return err
		}
		if err := r.db.Model(&entity.UserLiker{ID: *currentUser.Profile.CurrentRecommendationID}).Select("user_likers.liker_user_id").Take(likerUserID).Error; err != nil {
			return err
		}
		if err := r.db.Model(&entity.UserProfile{}).Update("in_discussion_with_user_id", likerUserID).Error; err != nil {
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

	if *currentUser.Profile.CurrentRecommendationType == "matching_profile" {
		if err := r.db.Model(&entity.UserMatchingProfile{ID: *currentUser.Profile.CurrentRecommendationID}).Update("action", "pass").Error; err != nil {
			return err
		}
	}
	if *currentUser.Profile.CurrentRecommendationType == "liker" {
		if err := r.db.Model(&entity.UserLiker{ID: *currentUser.Profile.CurrentRecommendationID}).Update("action", "pass").Error; err != nil {
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
