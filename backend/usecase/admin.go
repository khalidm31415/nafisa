package usecase

import (
	"backend/dto"
	"backend/entity"
	"backend/internal_constant"
	"backend/internal_error"
	elasticsarch_helper "backend/package_helper/elasticsearch_helper"
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAdmin interface {
	NewAdmin(ctx context.Context, input dto.NewAdminInput) error
	IndexProfile(ctx context.Context, userID string) error
}

type Admin struct {
	db           *gorm.DB
	profileIndex elasticsarch_helper.IElasticsearchProfileIndex
}

func NewAdmin(db *gorm.DB, profileIndex elasticsarch_helper.IElasticsearchProfileIndex) IAdmin {
	return &Admin{
		db:           db,
		profileIndex: profileIndex,
	}
}

func (a *Admin) NewAdmin(ctx context.Context, input dto.NewAdminInput) error {
	userID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", err))
		return err
	}
	hashedPasswordString := string(hashedPassword)

	admin := &entity.User{
		ID:                  userID.String(),
		OauthGmail:          &input.OauthGmail,
		Username:            &input.Username,
		Password:            &hashedPasswordString,
		IsVerificationAdmin: input.IsVerificationAdmin,
		IsDiscussionAdmin:   input.IsDiscussionAdmin,
	}

	if err := a.db.Create(admin).Error; err != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", err))
		return err
	}
	return nil
}

func (a *Admin) IndexProfile(ctx context.Context, userID string) error {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)
	if !currentUser.IsVerificationAdmin {
		return internal_error.ErrUnauthorized
	}

	var userProfile entity.UserProfile
	err := a.db.Take(&userProfile, "user_id = ?", userID).Error
	if err != nil {
		return err
	}

	if !userProfile.IsVerified {
		return internal_error.ErrUnverified
	}

	if err := a.profileIndex.Index(ctx, userProfile); err != nil {
		return err
	}

	return nil
}
