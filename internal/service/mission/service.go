package mission

import (
	"context"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

type MissionStorage interface {
	GetMissionEventsForUpdate(ctx context.Context) ([]models.MissionEvent, error)
	DeleteMissionEvents(ctx context.Context, eventsToDelete []models.MissionEvent) error
	ColonizePlanet(ctx context.Context, colonizeEvents models.MissionEvent) error
}

type txManager interface {
	ExecMissionTx(ctx context.Context, fn func(ctx context.Context, missionStorage MissionStorage) error) error
}

type Service struct {
	txManager txManager
	logger    *zap.Logger
}

func New(txManager txManager, logger *zap.Logger) *Service {
	return &Service{
		txManager: txManager,
		logger:    logger,
	}
}
