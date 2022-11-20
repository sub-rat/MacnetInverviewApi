package postgres

import (
	"fmt"
	"os"

	"github.com/sub-rat/machnet_api_assingment/internals/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(
		&models.User{},
		&models.UserFriend{},
		&models.Post{},
		&models.PostLike{},
		&models.PostShare{},
	)
	if err != nil {
		panic(err)
	}
	return db
}
