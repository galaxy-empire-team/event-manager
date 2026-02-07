package building

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *Service) HandleBuilds(ctx context.Context) error {
	err := s.txManager.ExecBuildingTx(ctx, func(ctx context.Context, buildingStorage BuildingStorage) error {
		buildEvents, err := buildingStorage.GetBuildEvents(ctx)
		if err != nil {
			return fmt.Errorf("buildingStorage.GetBuildEvents(): %w", err)
		}

		s.logger.Info("Fetched building events", zap.Int("count", len(buildEvents)))

		if len(buildEvents) == 0 {
			return nil
		}

		for _, buildEvent := range buildEvents {
			nextLvlBuilding, err := s.registry.GetBuildingNextLvlStats(consts.BuildingID(buildEvent.BuildingID))
			if err != nil {
				return fmt.Errorf("s.registry.GetNextLevelBuildingID(): %w", err)
			}

			updatedBuilding := models.BuildingUpgrade{
				PlanetID:          buildEvent.PlanetID,
				CurrentBuildingID: buildEvent.BuildingID,
				UpdatedBuildingID: nextLvlBuilding.ID,
			}

			err = buildingStorage.SetBuildingID(ctx, updatedBuilding)
			if err != nil {
				return fmt.Errorf("buildingStorage.SetBuildingID(): %w", err)
			}
		}

		err = buildingStorage.DeleteBuildEvents(ctx, buildEvents)
		if err != nil {
			return fmt.Errorf("buildingStorage.DeleteBuildEvents(): %w", err)
		}

		s.logger.Info("Completed upgrading buildings", zap.Int("count", len(buildEvents)))

		return nil
	})
	if err != nil {
		return fmt.Errorf("txManager.ExecBuildingsTx(): %w", err)
	}

	return nil
}
