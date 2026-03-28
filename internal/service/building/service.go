package building

import (
	"context"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

type BuildingStorage interface {
	GetBuildEvents(ctx context.Context, buildEventsCount uint16) ([]models.BuildEvent, error)
	DeleteBuildEvents(ctx context.Context, events []models.BuildEvent) error
	SetBuildingID(ctx context.Context, building models.BuildingUpgrade) error
}

type txManager interface {
	ExecBuildingTx(ctx context.Context, fn func(ctx context.Context, buildingStorage BuildingStorage) error) error
}

type registryProvider interface {
	GetBuildingNextLvlID(buildingID consts.BuildingID) (consts.BuildingID, error)
}

type Service struct {
	txManager txManager
	registry  registryProvider
	logger    *zap.Logger
}

func New(txManager txManager, registry registryProvider, logger *zap.Logger) *Service {
	return &Service{
		txManager: txManager,
		registry:  registry,
		logger:    logger,
	}
}
