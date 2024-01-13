package inits

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/mahalichev/WB-L0/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.DBOrder{})
	return db, err
}
