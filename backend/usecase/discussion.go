package usecase

import (
	"backend/entity"
	"backend/internal_constant"
	"backend/internal_error"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDiscussion interface {
	SendMessage(ctx context.Context, text string) error
	EndDiscussion(ctx context.Context) error
}

type Discussion struct {
	db *gorm.DB
}

func NewDiscussion(db *gorm.DB) IDiscussion {
	return &Discussion{
		db: db,
	}
}

func (d *Discussion) SendMessage(ctx context.Context, text string) error {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)

	err := d.db.Preload("Profile").Take(&currentUser).Error
	if err != nil {
		return err
	}

	if currentUser.Profile.CurrentDiscussionID == nil {
		return internal_error.ErrNoDiscussionInProgress
	}

	ID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	discussionMessage := entity.DiscussionMessage{
		ID:           ID.String(),
		UserID:       currentUser.ID,
		DiscussionID: *currentUser.Profile.CurrentDiscussionID,
		Text:         text,
	}
	if err := d.db.Create(&discussionMessage).Error; err != nil {
		return err
	}

	return nil
}

func (d *Discussion) EndDiscussion(ctx context.Context) error {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)

	err := d.db.Preload("Profile").Take(&currentUser).Error
	if err != nil {
		return err
	}

	if currentUser.Profile.CurrentDiscussionID == nil {
		return internal_error.ErrNoDiscussionInProgress
	}

	var discussion entity.Discussion
	if err := d.db.Model(&entity.Discussion{}).Where("id = ?", *currentUser.Profile.CurrentDiscussionID).Select("male_user_id", "female_user_id").Take(&discussion).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"ended_by_user_id": currentUser.ID,
		"ended_at":         time.Now(),
	}
	if err := d.db.Model(&entity.Discussion{}).Where("id = ?", *currentUser.Profile.CurrentDiscussionID).Updates(updates).Error; err != nil {
		return err
	}
	if err := d.db.Model(&entity.UserProfile{}).Where("user_id IN ?", []string{discussion.FemaleUserID, discussion.MaleUserID}).Update("current_discussion_id", nil).Error; err != nil {
		return err
	}

	return nil
}
