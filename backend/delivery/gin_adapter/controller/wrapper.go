package controller

import (
	"backend/usecase"

	"golang.org/x/oauth2"
)

type Controllers struct {
	Auth           IAuth
	Verification   IVerification
	Recommendation IRecommendation
}

func NewControllers(googleOauthConfig oauth2.Config, usecases usecase.Usecases) *Controllers {
	auth := NewAuth(googleOauthConfig, usecases)
	verification := NewVerification(usecases)
	recommendation := NewRecommendation(usecases)
	return &Controllers{
		Auth:           auth,
		Verification:   verification,
		Recommendation: recommendation,
	}
}
