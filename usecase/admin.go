package usecase

import (
	"context"
	"fmt"
	"nafisah/delivery/dto"
	"nafisah/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAdmin interface {
	NewAdmin(ctx context.Context, input dto.NewAdminInput) error
}

type Admin struct {
	db *gorm.DB
}

func NewAdmin(db *gorm.DB) IAdmin {
	return &Admin{
		db: db,
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

	admin := &entity.User{
		ID:                  userID.String(),
		Username:            input.Username,
		PasswordHash:        string(hashedPassword),
		IsVerificationAdmin: input.IsVerificationAdmin,
		IsDiscussionAdmin:   input.IsDiscussionAdmin,
	}

	if err := a.db.Create(admin).Error; err != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", err))
		return err
	}
	return nil
}
