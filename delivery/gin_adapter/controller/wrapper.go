package controller

import "nafisah/usecase"

type Controllers struct {
	Auth           IAuth
	Verification   IVerification
	Recommendation IRecommendation
}

func NewControllers(usecases usecase.Usecases) *Controllers {
	identityKey := "id"
	auth := NewAuth(usecases, identityKey)
	verification := NewVerification(usecases, identityKey)
	recommendation := NewRecommendation(usecases)
	return &Controllers{
		Auth:           auth,
		Verification:   verification,
		Recommendation: recommendation,
	}
}
