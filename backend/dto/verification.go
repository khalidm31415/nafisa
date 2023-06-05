package dto

type VerificationInput struct {
	UserID string `json:"userId" binding:"required"`
}

type ProfileToVerify struct {
	UserID              string `json:"userId"`
	SelfieWithIDCardURL string `json:"selfieWithIDCardURL"`
}
