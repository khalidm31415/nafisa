package usecase

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Usecases struct {
	User           IUser
	Admin          IAdmin
	Verification   IVerification
	Recommendation IRecommendation
}

func NewUsecases(db *gorm.DB, rdb *redis.Client) Usecases {
	user := NewUser(db)
	admin := NewAdmin(db)
	verification := NewVerification(db)
	recommendation := NewRecommendation(db, rdb)
	return Usecases{
		User:           user,
		Admin:          admin,
		Verification:   verification,
		Recommendation: recommendation,
	}
}
