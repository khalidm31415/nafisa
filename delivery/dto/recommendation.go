package dto

import (
	"nafisah/entity"
	"time"
)

type RecommendedProfile struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`

	Age           int    `json:"age"`
	Sex           string `json:"sex"`
	LastEducation string `json:"last_education"`
	Summary       string `json:"summary"`
}

func NewRecommendedProfile(user entity.User) RecommendedProfile {
	year, _, _ := time.Now().Date()
	age := year - user.Profile.YearBorn
	return RecommendedProfile{
		UserID:   user.ID,
		Username: user.Username,

		Age:           age,
		Sex:           user.Profile.Sex,
		LastEducation: user.Profile.LastEducation,
		Summary:       user.Profile.Summary,
	}
}
