package dto

import (
	"backend/entity"
	"time"
)

type RecommendedProfile struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`

	Age           int    `json:"age"`
	Gender        string `json:"gender"`
	LastEducation string `json:"last_education"`
	Summary       string `json:"summary"`
}

func NewRecommendedProfile(user entity.User) RecommendedProfile {
	year, _, _ := time.Now().Date()
	age := year - user.Profile.YearBorn
	return RecommendedProfile{
		UserID:   user.ID,
		Username: *user.Username,

		Age:           age,
		Gender:        user.Profile.Gender,
		LastEducation: user.Profile.LastEducation,
		Summary:       user.Profile.Summary,
	}
}
