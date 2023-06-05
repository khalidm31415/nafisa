package entity

type User struct {
	Base
	ID                  string `gorm:"type:varchar(36);primaryKey"`
	Username            string
	PasswordHash        string
	IsVerificationAdmin bool
	IsDiscussionAdmin   bool

	Profile          UserProfile
	Photos           []UserPhoto
	MatchingProfiles []UserMatchingProfile
	Likers           []UserLiker
	PendingMatches   []UserPendingMatch
}
