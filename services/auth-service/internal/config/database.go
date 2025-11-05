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
	if strings.Contains(dsn, "?sslmode=") && strings.Count(dsn, "?sslmode=") > 1 {
		// Find all SSL mode positions
		sslIndices := []int{}
		start := 0
		for {
			idx := strings.Index(dsn[start:], "?sslmode=")
			if idx == -1 {
				break
			}
			sslIndices = append(sslIndices, start+idx)
			start += idx + 1
		}

		// If we found more than one, remove the second one onwards
		if len(sslIndices) > 1 {
			// Keep everything up to the second SSL mode
			dsn = dsn[:sslIndices[1]]
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
