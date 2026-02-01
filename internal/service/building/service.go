package building

import (
	"context"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

type BuildingStorage interface {
	GetBuildEvents(ctx context.Context) ([]models.BuildEvent, error)
	DeleteBuildEvents(ctx context.Context, events []models.BuildEvent) error
	GetCurrentBuilding(ctx context.Context, building models.BuildEvent) (models.PlanetBuilding, error)
	CreateBuilding(ctx context.Context, building models.PlanetBuilding) error
	UpgradeBuildingLevel(ctx context.Context, building models.PlanetBuilding) error
}

type txManager interface {
	ExecBuildingTx(ctx context.Context, fn func(ctx context.Context, buildingStorage BuildingStorage) error) error
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
