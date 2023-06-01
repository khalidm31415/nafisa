package dto

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupInput struct {
	Username            string `json:"username" binding:"required"`
	Password            string `json:"password" binding:"required"`
	SelfieWithIDCardURL string `json:"selfieWithIDCardURL" binding:"required"`

	YearBorn      int      `json:"yearBorn" binding:"required"`
	Sex           string   `json:"sex" binding:"required,oneof=m f"`
	LastEducation string   `json:"lastEducation" binding:"required"`
	Summary       string   `json:"summary" binding:"required"`
	PhotoURLs     []string `json:"photoUrls" binding:"required"`

	PreferencePartnerCriteria  string `json:"preferencePartnerCriteria" binding:"required"`
	PreferenceMinLastEducation string `json:"preferenceMinLastEducation" binding:"required"`
	PreferenceMaxAge           int    `json:"preferenceMaxAge" binding:"required"`
	PreferenceMinAge           int    `json:"preferenceMinAge" binding:"required"`
}
