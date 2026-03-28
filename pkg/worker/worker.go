// Package worker provides a scalable wrapper around services.
// Unlike Kubernetes deployment scaling, it creates light workers
// that process transactions with configurable event limits.
package worker

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/config"
)

type service interface {
	Process(ctx context.Context, missionCount uint16) error
}

func StartWorker(ctx context.Context, cfg config.Worker, service service, logger *zap.Logger) {
	if cfg.ThreadCount == 0 {
		logger.Info("Worker thread count is set to 0, no workers will be started")
		return
	}

	// Distribute jobs within a second to make event processing more consistent.
	const millisecondInSecond = 1000
	delay := millisecondInSecond / int(cfg.ThreadCount)

	for range cfg.ThreadCount {
		time.Sleep(time.Millisecond * time.Duration(delay))

		go func() {
			ticker := time.NewTicker(cfg.TimeInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := service.Process(ctx, cfg.EventCount); err != nil {
						logger.Error("worker process error", zap.Error(err))
					}
				}
			}
		}()
	}
}
