package dto

import (
	"backend/entity"
)

type Profile struct {
	UserID        string  `json:"user_id"`
	Username      *string `json:"username"`
	OauthGmail    *string `json:"oauth_gmail"`
	IsVerified    bool    `json:"is_verified"`
	IsPremium     bool    `json:"is_premium"`
	YearBorn      int     `json:"year_born"`
	Gender        string  `json:"gender"`
	LastEducation string  `json:"last_education"`
	Summary       string  `json:"summary"`

	PreferencePartnerCriteria  string `json:"preference_partner_criteria"`
	PreferenceMinLastEducation string `json:"preference_min_last_education"`
	PreferenceMaxAge           int    `json:"preference_max_age"`
	PreferenceMinAge           int    `json:"preference_min_age"`
}

func NewProfile(user entity.User) Profile {
	profile := Profile{
		UserID:     user.ID,
		Username:   user.Username,
		OauthGmail: user.OauthGmail,

		IsVerified:    user.Profile.IsVerified,
		IsPremium:     user.Profile.IsPremium,
		YearBorn:      user.Profile.YearBorn,
		Gender:        user.Profile.Gender,
		LastEducation: user.Profile.LastEducation,
		Summary:       user.Profile.Summary,

		PreferencePartnerCriteria:  user.Profile.PreferencePartnerCriteria,
		PreferenceMinLastEducation: user.Profile.PreferenceMinLastEducation,
		PreferenceMaxAge:           user.Profile.PreferenceMaxAge,
		PreferenceMinAge:           user.Profile.PreferenceMinAge,
	}
	return profile
}

type ProfileIndex struct {
	UserID             string    `json:"user_id"`
	YearBorn           int       `json:"year_born"`
	Gender             string    `json:"gender"`
	LastEducation      string    `json:"last_education"`
	Summary            string    `json:"summary"`
	SummaryDenseVector []float32 `json:"summary_dense_vector"`
}

func NewProfileIndex(profile entity.UserProfile) (*ProfileIndex, error) {
	return &ProfileIndex{
		UserID:        profile.UserID,
		YearBorn:      profile.YearBorn,
		Gender:        profile.Gender,
		LastEducation: profile.LastEducation,
		Summary:       profile.Summary,
	}, nil
}

type CompleteProfileInput struct {
	Username            string `json:"username" binding:"required"`
	SelfieWithIDCardURL string `json:"selfieWithIDCardURL" binding:"required"`

	YearBorn      int      `json:"yearBorn" binding:"required"`
	Gender        string   `json:"gender" binding:"required,oneof=m f"`
	LastEducation string   `json:"lastEducation" binding:"required"`
	Summary       string   `json:"summary" binding:"required"`
	PhotoURLs     []string `json:"photoUrls" binding:"required"`

	PreferencePartnerCriteria  string `json:"preferencePartnerCriteria" binding:"required"`
	PreferenceMinLastEducation string `json:"preferenceMinLastEducation" binding:"required"`
	PreferenceMaxAge           int    `json:"preferenceMaxAge" binding:"required"`
	PreferenceMinAge           int    `json:"preferenceMinAge" binding:"required"`
}
