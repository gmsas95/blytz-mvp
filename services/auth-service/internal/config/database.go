package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL

	// Remove duplicate SSL mode parameters if they exist
	if strings.Contains(dsn, "?sslmode=disable?sslmode=") || strings.Contains(dsn, "?sslmode=require?sslmode=") {
		// Find and remove the duplicate SSL mode parameter
		if idx := strings.Index(dsn, "?sslmode="); idx != -1 {
			endIdx := strings.Index(dsn[idx+11:], "?")
			if endIdx != -1 {
				dsn = dsn[:idx] + dsn[idx+11+endIdx+1:]
			}
		}
	}

	gormConfig := &gorm.Config{}

	if cfg.Environment == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established with connection pooling")
	return db, nil
}
