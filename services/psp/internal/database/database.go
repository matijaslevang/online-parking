package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"psp/internal/models"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Merchant{})
}
