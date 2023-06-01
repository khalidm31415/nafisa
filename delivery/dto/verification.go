package dto

type VerificationInput struct {
	UserID string `json:"userId" binding:"required"`
}
