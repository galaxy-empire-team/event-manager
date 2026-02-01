package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/config"
)

type App struct {
	cancelFn context.CancelFunc

	logger *zap.Logger
}

func New(cfg config.App) (context.Context, *App, error) {
	ctx, cancelFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger, err := newLogger(cfg.Logger)
	if err != nil {
		return context.Background(), nil, fmt.Errorf("newLogger(): %w", err)
	}

	go func() {
		<-ctx.Done()
		logger.Sync() //nolint:errcheck,gosec
	}()

	return ctx, &App{
		cancelFn: cancelFn,
		logger:   logger,
	}, nil
}

func (a *App) ComponentLogger(component string) *zap.Logger {
	return a.logger.With(zap.String("component", component))
}

func (a *App) Logger() *zap.Logger {
	return a.logger
}

func (a *App) WaitShutdown(ctx context.Context) {
	<-ctx.Done()
}
