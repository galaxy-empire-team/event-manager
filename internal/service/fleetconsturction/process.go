package fleetconsturction

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *Service) Process(ctx context.Context, fleetConstructionEventsCount uint16) error {
	err := s.txManager.ExecFleetConstructionTx(ctx, func(ctx context.Context, fleetConstructionStorages TxStorages) error {
		fleetConstructionEvents, err := fleetConstructionStorages.GetFleetConstructionEvents(ctx, fleetConstructionEventsCount)
		if err != nil {
			return fmt.Errorf("fleetConstructionStorages.GetFleetConstructionEvents(): %w", err)
		}

		s.logger.Info("Fetched fleet construction events", zap.Int("count", len(fleetConstructionEvents)))

		if len(fleetConstructionEvents) == 0 {
			return nil
		}

		for _, fleetConstructionEvent := range fleetConstructionEvents {
			err := fleetConstructionStorages.AddFleet(ctx, fleetConstructionEvent.PlanetID, []models.FleetUnit{
				{
					ID:    fleetConstructionEvent.FleetID,
					Count: fleetConstructionEvent.Count,
				},
			})
			if err != nil {
				return fmt.Errorf("fleetConstructionStorages.AddFleetToPlanet(): %w", err)
			}
		}

		err = fleetConstructionStorages.DeleteFleetConstructionEvents(ctx, fleetConstructionEvents)
		if err != nil {
			return fmt.Errorf("fleetConstructionStorages.DeleteFleetConstructionEvents(): %w", err)
		}

		s.logger.Info("Completed upgrading fleet construction", zap.Int("count", len(fleetConstructionEvents)))

		return nil
	})
	if err != nil {
		return fmt.Errorf("txManager.ExecFleetConstructionTx(): %w", err)
	}

	return nil
}
