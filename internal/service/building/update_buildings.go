package planet

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (s *Service) UpdateBuildings(ctx context.Context) error {
	s.logger.Info("Starting UpdateBuildings process")

	err := s.txManager.ExecBuildingTx(ctx, func(ctx context.Context, buildingStorage BuildingStorage) error {
		buildEvents, err := buildingStorage.GetBuildingEvents(ctx)
		if err != nil {
			return fmt.Errorf("buildingStorage.GetBuildingEvents(): %w", err)
		}

		s.logger.Info("Fetched building events", zap.Int("count", len(buildEvents)))

		for _, buildEvent := range buildEvents {
			err := buildingStorage.UpgradeBuilding(ctx, buildEvent)
			if err != nil {
				return fmt.Errorf("buildingStorage.UpgradeBuilding(): %w", err)
			}

			err = buildingStorage.DeleteBuildingEvent(ctx, buildEvent)
			if err != nil {
				return fmt.Errorf("buildingStorage.DeleteBuildingEvent(): %w", err)
			}
		}

		s.logger.Info("Completed upgrading buildings", zap.Int("count", len(buildEvents)))

		return nil
	})
	if err != nil {
		return fmt.Errorf("txManager.ExecBuildingsTx(): %w", err)
	}

	return nil
}
