package utils

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger(env string) (*zap.Logger, error) {
	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	var err error
	logger, err = config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func GetLogger() *zap.Logger {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return logger
}
