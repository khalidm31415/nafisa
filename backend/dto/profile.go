package dto

import (
	"backend/entity"
)

type ProfileIndex struct {
	UserID             string    `json:"user_id"`
	YearBorn           int       `json:"year_born"`
	Sex                string    `json:"sex"`
	LastEducation      string    `json:"last_education"`
	Summary            string    `json:"summary"`
	SummaryDenseVector []float32 `json:"summary_dense_vector"`
}

func NewProfileIndex(profile entity.UserProfile) (*ProfileIndex, error) {
	return &ProfileIndex{
		UserID:        profile.UserID,
		YearBorn:      profile.YearBorn,
		Sex:           profile.Sex,
		LastEducation: profile.LastEducation,
		Summary:       profile.Summary,
	}, nil
}
