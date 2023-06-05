package dto

import (
	"backend/entity"
)

type Profile struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	IsVerified    bool   `json:"is_verified"`
	IsPremium     bool   `json:"is_premium"`
	YearBorn      int    `json:"year_born"`
	Sex           string `json:"sex"`
	LastEducation string `json:"last_education"`
	Summary       string `json:"summary"`

	PreferencePartnerCriteria  string `json:"preference_partner_criteria"`
	PreferenceMinLastEducation string `json:"preference_min_last_education"`
	PreferenceMaxAge           int    `json:"preference_max_age"`
	PreferenceMinAge           int    `json:"preference_min_age"`
}

func NewProfile(user entity.User) Profile {
	return Profile{
		UserID:        user.ID,
		Username:      user.Username,
		IsVerified:    user.Profile.IsVerified,
		IsPremium:     user.Profile.IsPremium,
		YearBorn:      user.Profile.YearBorn,
		Sex:           user.Profile.Sex,
		LastEducation: user.Profile.LastEducation,
		Summary:       user.Profile.Summary,

		PreferencePartnerCriteria:  user.Profile.PreferencePartnerCriteria,
		PreferenceMinLastEducation: user.Profile.PreferenceMinLastEducation,
		PreferenceMaxAge:           user.Profile.PreferenceMaxAge,
		PreferenceMinAge:           user.Profile.PreferenceMinAge,
	}
}

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
