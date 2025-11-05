package config

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL

	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.Environment == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Configure connection pool
	maxIdleConns := 10
	maxOpenConns := 100
	connMaxLifetime := time.Hour

	if cfg.Environment == "production" {
		maxIdleConns = 25
		maxOpenConns = 200
		connMaxLifetime = 2 * time.Hour
	}

	// Open database connection with pooling
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	log.Printf("Database connected successfully with pool: maxIdle=%d, maxOpen=%d, lifetime=%v",
		maxIdleConns, maxOpenConns, connMaxLifetime)

	return db, nil
}
