package entity

type UserMatchingProfile struct {
	Base
	ID                    string `gorm:"type:varchar(36);primaryKey"`
	UserID                string `gorm:"type:varchar(36)"`
	MatchingProfileUserID string `gorm:"type:varchar(36)"`
	MatchingProfileUser   User   `gorm:"foreignKey:MatchingProfileUserID"`
	Score                 float32
	Action                *string
}

type Like struct {
	Base
	ID          string `gorm:"type:varchar(36);primaryKey"`
	UserID      string `gorm:"type:varchar(36)"`
	User        User
	LikedUserID string `gorm:"type:varchar(36)"`
	IsLikedBack *bool
}

// Used for pending matches
type Match struct {
	Base
	ID           string `gorm:"type:varchar(36);primaryKey"`
	MaleUserID   string `gorm:"type:varchar(36)"`
	MaleUser     User   `gorm:"foreignKey:MaleUserID"`
	FemaleUserID string `gorm:"type:varchar(36)"`
	FemaleUser   User   `gorm:"foreignKey:FemaleUserID"`
}
