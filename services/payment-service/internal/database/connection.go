package database

import (
	"fmt"

	"github.com/blytz/payment-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Payment{},
	)
}