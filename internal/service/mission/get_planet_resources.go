package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *Service) getResourcesForUpdate(ctx context.Context, planetID uuid.UUID, updatedAt time.Time, storage TxStorages) (models.Resources, error) {
	planetBuildingsInfo, err := s.getMinesInfo(ctx, planetID, storage)
	if err != nil {
		return models.Resources{}, fmt.Errorf("GetMinesInfo(): %w", err)
	}

	resources, err := storage.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return models.Resources{}, fmt.Errorf("storage.GetResourcesForUpdate(): %w", err)
	}

	millisecondsSinceLastUpdate := updatedAt.Sub(resources.UpdatedAt).Milliseconds()
	if millisecondsSinceLastUpdate <= 0 {
		return models.Resources{}, nil
	}

	updatedResources := models.Resources{
		Metal:     resources.Metal + uint64(millisecondsSinceLastUpdate)*planetBuildingsInfo[consts.BuildingTypeMetalMine].ProductionS/1000,
		Crystal:   resources.Crystal + uint64(millisecondsSinceLastUpdate)*planetBuildingsInfo[consts.BuildingTypeCrystalMine].ProductionS/1000,
		Gas:       resources.Gas + uint64(millisecondsSinceLastUpdate)*planetBuildingsInfo[consts.BuildingTypeGasMine].ProductionS/1000,
		UpdatedAt: updatedAt,
	}

	return updatedResources, nil
}

func (s *Service) getMinesInfo(ctx context.Context, planetID uuid.UUID, storage TxStorages) (map[consts.BuildingType]models.BuildingInfo, error) {
	mines, err := storage.GetBuildingsInfoByTypes(ctx, planetID, consts.GetMineTypes())
	if err != nil {
		return nil, fmt.Errorf("storage.GetBuildingsInfoByTypes(): %w", err)
	}

	// If mines are not build yet, initialize them with default values
	for _, mineType := range consts.GetMineTypes() {
		if _, exists := mines[mineType]; !exists {
			stat, err := s.registry.GetBuildingZeroLvlStats(mineType)
			if err != nil {
				return nil, fmt.Errorf("registry.GetBuildingZeroLvlStats(): %w", err)
			}

			mines[mineType] = models.BuildingInfo{
				ID:          stat.ID,
				Type:        stat.Type,
				ProductionS: stat.ProductionS,
			}
		}
	}

	return mines, nil
}
