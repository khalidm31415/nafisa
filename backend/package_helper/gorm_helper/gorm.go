package gorm_helper

import (
	"backend/entity"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := os.Getenv("BACKEND_MYSQL_DSN")
	db, err := gorm.Open(mysql.Open(dsn))
	// db = db.Debug()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.UserProfile{})
	db.AutoMigrate(&entity.UserPhoto{})
	db.AutoMigrate(&entity.UserMatchingProfile{})
	db.AutoMigrate(&entity.Like{})
	db.AutoMigrate(&entity.Match{})
	db.AutoMigrate(&entity.Discussion{})
	db.AutoMigrate(&entity.DiscussionMessage{})
	return db
}
