package mission

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

type TxStorages interface {
	// mission storage
	CreateMissionEvent(ctx context.Context, missionEvent models.MissionEvent) error
	GetMissionEventsForUpdate(ctx context.Context) ([]models.MissionEvent, error)
	DeleteMissionEvents(ctx context.Context, eventsToDelete []models.MissionEvent) error

	// planet storage
	ColonizePlanet(ctx context.Context, colonizeEvents models.MissionEvent) (colonized bool, err error)
	GetPlanetInfoByCoordinates(ctx context.Context, planetFrom models.Coordinates) (models.Planet, error)
	GetPlanetInfoByID(ctx context.Context, planetID uuid.UUID) (models.Planet, error)
	GetPlanetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	GetPlanetFleetForUpdate(ctx context.Context, planetID uuid.UUID) ([]models.FleetUnit, error)
	SetPlanetFleet(ctx context.Context, planetID uuid.UUID, fleet []models.FleetUnit) error
	UpsertFleet(ctx context.Context, planetID uuid.UUID, fleet []models.FleetUnit) error

	// notification storage
	SaveNotificationEvents(ctx context.Context, notificationEvents []models.NotificationEvent) error
}

type txManager interface {
	ExecMissionTx(ctx context.Context, fn func(ctx context.Context, txStorages TxStorages) error) error
}

type registryProvider interface {
	GetMissionTypeByID(missionID consts.MissionID) (consts.MissionType, error)
	GetMissionIDByType(missionType consts.MissionType) (consts.MissionID, error)
	GetNotificationIDByType(notificationType consts.NotificationType) (consts.NotificationID, error)
}

type Service struct {
	txManager        txManager
	registryProvider registryProvider
	logger           *zap.Logger
}

func New(txManager txManager, registryProvider registryProvider, logger *zap.Logger) *Service {
	return &Service{
		txManager:        txManager,
		registryProvider: registryProvider,
		logger:           logger,
	}
}
