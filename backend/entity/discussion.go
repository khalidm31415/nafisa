package entity

type Discussion struct {
	Base
	ID           string `gorm:"type:varchar(36);primaryKey"`
	MaleUserID   string `gorm:"type:varchar(36)"`
	MaleUser     User   `gorm:"foreignKey:MaleUserID"`
	FemaleUserID string `gorm:"type:varchar(36)"`
	FemaleUser   User   `gorm:"foreignKey:FemaleUserID"`
	IsCompleted  bool
}

type DiscussionMessage struct {
	Base
	ID           string `gorm:"type:varchar(36);primaryKey"`
	UserID       string `gorm:"type:varchar(36)"`
	User         User
	DiscussionID string `gorm:"type:varchar(36)"`
	Discussion   Discussion
	Text         string
}
