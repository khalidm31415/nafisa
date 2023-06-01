package usecase

import (
	"context"
	"nafisah/entity"
	"nafisah/internal_constant"
	"nafisah/internal_error"

	"gorm.io/gorm"
)

type IVerification interface {
	Verify(ctx context.Context, userID string) error
}

type Verification struct {
	db *gorm.DB
}

func NewVerification(db *gorm.DB) IVerification {
	return &Verification{
		db: db,
	}
}

func (v *Verification) Verify(ctx context.Context, userID string) error {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)
	if !currentUser.IsVerificationAdmin {
		return internal_error.ErrUnauthorized
	}
	if err := v.db.Model(&entity.UserProfile{}).Where("user_id = ?", userID).Update("is_verified", true).Error; err != nil {
		return err
	}
	return nil
}
