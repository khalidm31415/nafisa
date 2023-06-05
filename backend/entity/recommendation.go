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

type UserLiker struct {
	Base
	ID          string `gorm:"type:varchar(36);primaryKey"`
	UserID      string `gorm:"type:varchar(36)"`
	LikerUserID string `gorm:"type:varchar(36)"`
	LikerUser   User   `gorm:"foreignKey:LikerUserID"`
	Action      *string
}

type UserPendingMatch struct {
	Base
	ID                 string `gorm:"type:varchar(36);primaryKey"`
	UserID             string `gorm:"type:varchar(36)"`
	PendingMatchUserID string `gorm:"type:varchar(36)"`
	PendingMatchUser   User   `gorm:"foreignKey:PendingMatchUserID"`
}
