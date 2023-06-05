package usecase

import (
	"backend/dto"
	"backend/entity"
	"backend/internal_constant"
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUser interface {
	Signup(ctx context.Context, input dto.SignupInput) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByUserID(ctx context.Context, userID string) (*entity.User, error)
	CurrentUserProfile(ctx context.Context) (*dto.Profile, error)
}

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) IUser {
	return &User{
		db: db,
	}
}

func (u *User) Signup(ctx context.Context, input dto.SignupInput) error {
	userID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	photos := []entity.UserPhoto{}
	for _, p := range input.PhotoURLs {
		photoID, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		photos = append(photos, entity.UserPhoto{
			ID:       photoID.String(),
			UserID:   userID.String(),
			PhotoURL: p,
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", err))
		return err
	}

	profile := entity.UserProfile{
		UserID:              userID.String(),
		SelfieWithIDCardURL: input.SelfieWithIDCardURL,

		YearBorn:      input.YearBorn,
		Sex:           input.Sex,
		LastEducation: input.LastEducation,
		Summary:       input.Summary,

		PreferencePartnerCriteria:  input.PreferencePartnerCriteria,
		PreferenceMinLastEducation: input.PreferenceMinLastEducation,
		PreferenceMaxAge:           input.PreferenceMaxAge,
		PreferenceMinAge:           input.PreferenceMinAge,
	}

	user := entity.User{
		ID:           userID.String(),
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		Profile:      profile,
		Photos:       photos,
	}

	if err := u.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user *entity.User
	if err := u.db.Take(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) FindByUserID(ctx context.Context, userID string) (*entity.User, error) {
	var user *entity.User
	if err := u.db.Take(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) CurrentUserProfile(ctx context.Context) (*dto.Profile, error) {
	currentUser := ctx.Value(internal_constant.ContextUserKey).(*entity.User)

	err := u.db.Preload("Profile").Take(&currentUser).Error
	if err != nil {
		return nil, err
	}

	profile := dto.NewProfile(*currentUser)
	return &profile, nil
}
