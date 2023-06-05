package usecase

import (
	elasticsarch_helper "backend/package_helper/elasticsearch_helper"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Usecases struct {
	User           IUser
	Admin          IAdmin
	Verification   IVerification
	Recommendation IRecommendation
}

func NewUsecases(db *gorm.DB, rdb *redis.Client, profileIndex elasticsarch_helper.IElasticsearchProfileIndex) Usecases {
	user := NewUser(db)
	admin := NewAdmin(db, profileIndex)
	verification := NewVerification(db)
	recommendation := NewRecommendation(db, rdb, profileIndex)
	return Usecases{
		User:           user,
		Admin:          admin,
		Verification:   verification,
		Recommendation: recommendation,
	}
}
