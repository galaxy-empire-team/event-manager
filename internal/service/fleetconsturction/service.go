package fleetconsturction

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

type TxStorages interface {
	GetFleetConstructionEvents(ctx context.Context, fleetConstructionEventsCount uint16) ([]models.FleetConstructionEvent, error)
	DeleteFleetConstructionEvents(ctx context.Context, events []models.FleetConstructionEvent) error
	AddFleet(ctx context.Context, planetID uuid.UUID, fleetConstruction []models.FleetUnit) error
}

type txManager interface {
	ExecFleetConstructionTx(ctx context.Context, fn func(ctx context.Context, fleetConstructionStorage TxStorages) error) error
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
