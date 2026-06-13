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
	CreateMoon(ctx context.Context, planetID uuid.UUID) error
	GetPlanetIDByCoordinates(ctx context.Context, coordinates models.Coordinates) (uuid.UUID, error)
	GetPlanetCoordinatesByID(ctx context.Context, planetID uuid.UUID) (models.Coordinates, error)
	GetPlanetInfoByCoordinates(ctx context.Context, planetFrom models.Coordinates) (models.Planet, error)
	GetPlanetInfoByID(ctx context.Context, planetID uuid.UUID) (models.Planet, error)
	AddResources(ctx context.Context, planetID uuid.UUID, resources models.Resources) error
	GetResources(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
	GetPlanetFleetForUpdate(ctx context.Context, planetID uuid.UUID) ([]models.FleetUnit, error)
	SetPlanetFleet(ctx context.Context, planetID uuid.UUID, fleet []models.FleetUnit) error
	AddFleet(ctx context.Context, planetID uuid.UUID, fleet []models.FleetUnit) error
	AddDebris(ctx context.Context, planetID uuid.UUID, debris models.Resources) error
	GetDebrisForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SetDebris(ctx context.Context, planetID uuid.UUID, debris models.Resources) error
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
	GetNPCStatsByPosition(positionZ consts.PlanetPositionZ) (registry.NPCStats, error)
}

type repository interface {
	GetResearchByType(ctx context.Context, userID uuid.UUID, researchType consts.ResearchType) (registry.ResearchStats, error)
}

//go:generate mockery --name=randGenerator --filename=rand_generator.go --exported --with-expecter
type randGenerator interface {
	Intn(n int) int
}

type Service struct {
	txManager       txManager
	bridgeAPIClient bridgeAPIClient
	registry        registryProvider
	repository      repository
	randGenerator   randGenerator
	logger          *zap.Logger
}

func New(txManager txManager, bridgeAPIClient bridgeAPIClient, repository repository, registryProvider registryProvider, logger *zap.Logger) *Service {
	return &Service{
		txManager:       txManager,
		bridgeAPIClient: bridgeAPIClient,
		registry:        registryProvider,
		repository:      repository,
		randGenerator:   rand.New(rand.NewSource(time.Now().UnixNano())),
		logger:          logger,
	}
}
