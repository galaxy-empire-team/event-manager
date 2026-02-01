package building

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

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
			upgradedBuilding, err := buildingStorage.GetCurrentBuilding(ctx, buildEvent)
			if err != nil {
				if !errors.Is(err, models.ErrBuildingNotFound) {
					return fmt.Errorf("buildingStorage.GetCurrentBuilding(): %w", err)
				}

				upgradedBuilding = models.PlanetBuilding{
					PlanetID:  buildEvent.PlanetID,
					BuildType: buildEvent.BuildType,
					Level:     0,
				}

				err = buildingStorage.CreateBuilding(ctx, upgradedBuilding)
				if err != nil {
					return fmt.Errorf("buildingStorage.CreateBuilding(): %w", err)
				}

				continue
			}

			upgradedBuilding.Level++

			err = buildingStorage.UpgradeBuildingLevel(ctx, upgradedBuilding)
			if err != nil {
				return fmt.Errorf("buildingStorage.UpgradeBuildingLevel(): %w", err)
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
