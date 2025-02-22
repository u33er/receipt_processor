package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

const (
	localEnv       = "local"
	productionEnv  = "prod"
	developmentEnv = "dev"
	uatEnv         = "uat"
)

func NewLogger(env string) *zap.Logger {
	var config zap.Config

	switch env {
	case localEnv, developmentEnv, uatEnv:
		config = zap.NewDevelopmentConfig()
	case productionEnv:
		config.EncoderConfig.CallerKey = ""
		config.DisableStacktrace = true
		config.DisableCaller = true
		config = zap.NewProductionConfig()

	default:
		config.EncoderConfig.CallerKey = ""
		config.DisableStacktrace = true
		config.DisableCaller = true
		config = zap.NewProductionConfig()
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	return logger
}

func SetupLogger(env string) *zap.Logger {
	logger := NewLogger(env)
	zap.ReplaceGlobals(logger)
	return logger
}
