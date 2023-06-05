package controller

import "backend/usecase"

type Controllers struct {
	Auth           IAuth
	Verification   IVerification
	Recommendation IRecommendation
}

func NewControllers(usecases usecase.Usecases) *Controllers {
	auth := NewAuth(usecases)
	verification := NewVerification(usecases)
	recommendation := NewRecommendation(usecases)
	return &Controllers{
		Auth:           auth,
		Verification:   verification,
		Recommendation: recommendation,
	}
}
