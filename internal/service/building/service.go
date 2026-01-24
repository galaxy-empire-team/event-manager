package planet

import (
	"context"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const (
	galaxyCount          = 1
	systemInGalaxyCount  = 3
	planetsInSystemCount = 16

	defaultLvl = 0
)

type BuildingStorage interface {
	GetBuildingEvents(ctx context.Context) ([]models.BuildEvent, error)
	UpgradeBuilding(ctx context.Context, building models.BuildEvent) error
	DeleteBuildingEvent(ctx context.Context, event models.BuildEvent) error
}

type txManager interface {
	ExecBuildingTx(ctx context.Context, fn func(ctx context.Context, buildingStorage BuildingStorage) error) error
}

type Service struct {
	txManager txManager

	logger *zap.Logger
}

func New(txManager txManager, logger *zap.Logger) *Service {
	return &Service{
		txManager: txManager,
		logger:    logger,
	}
}
