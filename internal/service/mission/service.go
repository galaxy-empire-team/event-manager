package mission

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

//go:generate mockery --name=TxStorages --filename=tx_storages.go --exported --with-expecter
type TxStorages interface {
	// mission storage
	CreateMissionEvent(ctx context.Context, missionEvent models.MissionEvent) error
	GetMissionEvents(ctx context.Context, missionEventsCount uint16) ([]models.MissionEvent, error)
	DeleteMissionEvents(ctx context.Context, eventsToDelete []models.MissionEvent) error

	// planet storage
	GetPlanetInfoByCoordinates(ctx context.Context, planetFrom models.Coordinates) (models.Planet, error)
	GetPlanetInfoByID(ctx context.Context, planetID uuid.UUID) (models.Planet, error)
	GetResources(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
	GetPlanetFleetForUpdate(ctx context.Context, planetID uuid.UUID) ([]models.FleetUnit, error)
	SetPlanetFleet(ctx context.Context, planetID uuid.UUID, fleet []models.FleetUnit) error
	AddFleet(ctx context.Context, planetID uuid.UUID, fleet []models.FleetUnit) error
	GetBuildings(ctx context.Context, planetID uuid.UUID) ([]consts.BuildingID, error)
	GetPlanetMinesProduction(ctx context.Context, planetID uuid.UUID) (map[consts.BuildingType]uint64, error)

	// notification storage
	SaveNotificationEvents(ctx context.Context, notificationEvents []models.NotificationEvent) error

	// research storage
	GetUserResearches(ctx context.Context, userID uuid.UUID) ([]consts.ResearchID, error)
	GetUserResearchesByTypes(ctx context.Context, userID uuid.UUID, researchTypes []consts.ResearchType) (map[consts.ResearchType]consts.ResearchID, error)
}

type txManager interface {
	ExecMissionTx(ctx context.Context, fn func(ctx context.Context, txStorages TxStorages) error) error
}

type bridgeAPIClient interface {
	ColonizePlanet(ctx context.Context, userID uuid.UUID, colonizeEvent models.MissionEvent) error
	UpdatePlanetResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedTime time.Time) error
}

//go:generate mockery --name=registryProvider --filename=registry_provider.go --exported --with-expecter
type registryProvider interface {
	GetMissionTypeByID(missionID consts.MissionID) (consts.MissionType, error)
	GetMissionIDByType(missionType consts.MissionType) (consts.MissionID, error)
	GetNotificationIDByType(notificationType consts.NotificationType) (consts.NotificationID, error)
	GetResearchStatsByID(researchID consts.ResearchID) (registry.ResearchStats, error)
	GetResearchZeroLvlIDByType(researchType consts.ResearchType) (consts.ResearchID, error)
	GetFleetUnitStatsByID(fleetUnitID consts.FleetUnitID) (registry.FleetUnitStats, error)
}

//go:generate mockery --name=randGenerator --filename=rand_generator.go --exported --with-expecter
type randGenerator interface {
	Intn(n int) int
}

type Service struct {
	txManager       txManager
	bridgeAPIClient bridgeAPIClient
	registry        registryProvider
	randGenerator   randGenerator
	logger          *zap.Logger
}

func New(txManager txManager, bridgeAPIClient bridgeAPIClient, registryProvider registryProvider, logger *zap.Logger) *Service {
	return &Service{
		txManager:       txManager,
		bridgeAPIClient: bridgeAPIClient,
		registry:        registryProvider,
		randGenerator:   rand.New(rand.NewSource(time.Now().UnixNano())),
		logger:          logger,
	}
}
