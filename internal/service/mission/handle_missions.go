package mission

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *Service) HandleMissions(ctx context.Context) error {
	err := s.txManager.ExecMissionTx(ctx, func(ctx context.Context, missionStorage MissionStorage) error {
		missionEvents, err := missionStorage.GetMissionEventsForUpdate(ctx)
		if err != nil {
			return fmt.Errorf("missionStorage.GetMissionEvents(): %w", err)
		}

		s.logger.Info("Fetched mission events", zap.Int("count", len(missionEvents)))
		
		if len(missionEvents) == 0 {
			return nil
		}

		for _, missionEvent := range missionEvents {
			if missionEvent.Type == models.MissionTypeColonize {
				err := missionStorage.ColonizePlanet(ctx, missionEvent)
				if err != nil {
					return fmt.Errorf("missionStorage.ColonizePlanet(): %w", err)
				}
			}
		}

		err = missionStorage.DeleteMissionEvents(ctx, missionEvents)
		if err != nil {
			return fmt.Errorf("missionStorage.DeleteMissionEvents(): %w", err)
		}

		s.logger.Info("Completed handling missions", zap.Int("count", len(missionEvents)))

		return nil
	})
	if err != nil {
		return fmt.Errorf("txManager.ExecMissionTx(): %w", err)
	}

	return nil
}
