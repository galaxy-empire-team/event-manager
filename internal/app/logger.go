package app

import (
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(cfg config.Logger) (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("zap.NewAtomicLevelAt(): %w", err)
	}

	config := zap.Config{
		Level:             level,
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          cfg.Format,
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stdout"},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("config.Build(): %w", err)
	}

	return logger, nil
}
