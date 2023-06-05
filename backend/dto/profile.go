package dto

import (
	"backend/entity"
	"backend/package_helper/embeddings_helper"
	"context"
	"fmt"
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
	summary_dense_vector, err := embeddings_helper.Embed(context.Background(), profile.Summary)
	if err != nil {
		fmt.Println("Error embedding: ", err.Error())
		return nil, err
	}
	return &ProfileIndex{
		UserID:             profile.UserID,
		YearBorn:           profile.YearBorn,
		Sex:                profile.Sex,
		LastEducation:      profile.LastEducation,
		Summary:            profile.Summary,
		SummaryDenseVector: summary_dense_vector,
	}, nil
}
