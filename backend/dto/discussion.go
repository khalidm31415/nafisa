package dto

type SendMessageInput struct {
	Text string `json:"text" binding:"required"`
}
