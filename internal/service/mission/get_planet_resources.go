package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *Service) getResourcesForUpdate(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedAt time.Time, storage TxStorages) (models.Resources, error) {
	multiplier, err := s.getResearchResourceMultiplier(ctx, userID, storage)
	if err != nil {
		return models.Resources{}, fmt.Errorf("getResearchResourceMultiplier(): %w", err)
	}

	mines, err := storage.GetPlanetMinesProduction(ctx, planetID)
	if err != nil {
		return models.Resources{}, fmt.Errorf("storage.GetPlanetMinesProduction(): %w", err)
	}

	resources, err := storage.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return models.Resources{}, fmt.Errorf("storage.GetResourcesForUpdate(): %w", err)
	}

	millisecondsSinceLastUpdate := updatedAt.Sub(resources.UpdatedAt).Milliseconds()
	if millisecondsSinceLastUpdate <= 0 {
		return models.Resources{}, nil
	}

	metalProductionPerSecond := float32(mines[consts.BuildingTypeMetalMine]) * multiplier
	crystalProductionPerSecond := float32(mines[consts.BuildingTypeCrystalMine]) * multiplier
	gasProductionPerSecond := float32(mines[consts.BuildingTypeGasMine]) * multiplier

	updatedResources := models.Resources{
		Metal:     resources.Metal + uint64(millisecondsSinceLastUpdate)*uint64(metalProductionPerSecond)/1000,
		Crystal:   resources.Crystal + uint64(millisecondsSinceLastUpdate)*uint64(crystalProductionPerSecond)/1000,
		Gas:       resources.Gas + uint64(millisecondsSinceLastUpdate)*uint64(gasProductionPerSecond)/1000,
		UpdatedAt: updatedAt,
	}

	return updatedResources, nil
}

func (s *Service) getResearchResourceMultiplier(ctx context.Context, userID uuid.UUID, storage TxStorages) (float32, error) {
	researchIDs, err := storage.GetUserResearches(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("storage.GetUserResearches(): %w", err)
	}

	for _, researchID := range researchIDs {
		research, err := s.registry.GetResearchStatsByID(researchID)
		if err != nil {
			return 0, fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
		}

		if research.Type != consts.ResearchTypeIndustrialTechnology {
			continue
		}

		return research.Bonuses.ProductionSpeedImprove, nil
	}

	// If user has no industrial technology research, return 1
	return 1, nil
}
