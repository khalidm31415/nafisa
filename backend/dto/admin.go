package dto

type NewAdminInput struct {
	OauthGmail          string `json:"oauthGmail" binding:"required"`
	Username            string `json:"username" binding:"required"`
	Password            string `json:"password" binding:"required"`
	IsVerificationAdmin bool   `json:"isVerificationAdmin" binding:"required"`
	IsDiscussionAdmin   bool   `json:"isDiscussionAdmin" binding:"required"`
}
