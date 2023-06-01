package pkg

import (
	"log"
	"nafisah/entity"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := os.Getenv("API_MYSQL_DSN")
	db, err := gorm.Open(mysql.Open(dsn))
	db = db.Debug()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.UserProfile{})
	db.AutoMigrate(&entity.UserPhoto{})
	db.AutoMigrate(&entity.UserMatchingProfile{})
	db.AutoMigrate(&entity.UserLiker{})
	db.AutoMigrate(&entity.UserPendingMatch{})
	return db
}
