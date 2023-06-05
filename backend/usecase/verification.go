package usecase

import (
	"backend/dto"
	"backend/entity"
	"backend/internal_constant"
	"backend/internal_error"
	"context"

	"gorm.io/gorm"
)

type IVerification interface {
	GetUnverifiedUsers(ctx context.Context) ([]dto.ProfileToVerify, error)
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

func (v *Verification) GetUnverifiedUsers(ctx context.Context) ([]dto.ProfileToVerify, error) {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)
	if !currentUser.IsVerificationAdmin {
		return nil, internal_error.ErrUnauthorized
	}
	var fullProfiles []entity.UserProfile
	if err := v.db.Where("is_verified = ?", false).Find(&fullProfiles).Error; err != nil {
		return nil, err
	}

	profilesToVerify := []dto.ProfileToVerify{}
	for _, profile := range fullProfiles {
		profilesToVerify = append(profilesToVerify, dto.ProfileToVerify{
			UserID:              profile.UserID,
			SelfieWithIDCardURL: profile.SelfieWithIDCardURL,
		})
	}
	return profilesToVerify, nil
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
