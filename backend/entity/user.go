package entity

type User struct {
	Base
	ID                  string  `gorm:"type:varchar(36);primaryKey"`
	OauthGmail          *string `gorm:"type:varchar(320);unique"`
	Username            *string `gorm:"type:varchar(50);unique"`
	Password            *string
	IsVerificationAdmin bool
	IsDiscussionAdmin   bool

	Profile          UserProfile
	Photos           []UserPhoto
	MatchingProfiles []UserMatchingProfile
	Likers           []UserLiker
	PendingMatches   []UserPendingMatch
}
