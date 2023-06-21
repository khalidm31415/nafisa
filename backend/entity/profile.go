package entity

type UserProfile struct {
	Base
	UserID              string `gorm:"type:varchar(36);primaryKey"`
	SelfieWithIDCardURL string
	IsVerified          bool
	IsPremium           bool
	IsProfileComplete   bool

	YearBorn      int
	Sex           string
	LastEducation string
	Summary       string

	PreferencePartnerCriteria  string
	PreferenceMinLastEducation string
	PreferenceMaxAge           int
	PreferenceMinAge           int

	CurrentRecommendationID   *string `gorm:"type:varchar(36)"`
	CurrentRecommendationType *string

	InDiscussionWithUserID *string
}

type UserPhoto struct {
	Base
	ID       string `gorm:"type:varchar(36);primaryKey"`
	UserID   string `gorm:"type:varchar(36)"`
	PhotoURL string
}
