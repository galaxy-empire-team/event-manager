package mission

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) HandleMissions(ctx context.Context) error {
	err := s.txManager.ExecMissionTx(ctx, func(ctx context.Context, txStorages TxStorages) error {
		missionEvents, err := txStorages.GetMissionEventsForUpdate(ctx)
		if err != nil {
			return fmt.Errorf("txStorages.GetMissionEventsForUpdate(): %w", err)
		}

		s.logger.Info("Fetched mission events", zap.Int("count", len(missionEvents)))

		if len(missionEvents) == 0 {
			return nil
		}

		for _, missionEvent := range missionEvents {
			mType, err := s.registryProvider.GetMissionTypeByID(missionEvent.MissionID)
			if err != nil {
				s.logger.Error("get mission type from registry", zap.Any("missionEvent", missionEvent), zap.Error(err))
				continue
			}

			// I don't want to create return mission type cause I need to store previous type somehow
			if missionEvent.IsReturning == true {
				err := s.returnMission(ctx, missionEvent, txStorages)
				if err != nil {
					return fmt.Errorf("s.returnMission(): %w", err)
				}
				break
			}

			switch mType {
			case consts.MissionTypeColonize:
				err := s.handleColonization(ctx, missionEvent, txStorages)
				if err != nil {
					return fmt.Errorf("s.handleColonization(): %w", err)
				}
			case consts.MissionTypeAttack:
				err := s.handleAttack(ctx, missionEvent, txStorages)
				if err != nil {
					return fmt.Errorf("s.handleAttack(): %w", err)
				}
			default:
				s.logger.Warn("Unknown mission type", zap.Any("missionEvent", missionEvent))
			}
		}

		err = txStorages.DeleteMissionEvents(ctx, missionEvents)
		if err != nil {
			return fmt.Errorf("txStorages.DeleteMissionEvents(): %w", err)
		}

		s.logger.Info("Completed handling missions", zap.Int("count", len(missionEvents)))

		return nil
	})
	if err != nil {
		return fmt.Errorf("txManager.ExecMissionTx(): %w", err)
	}

	return nil
}
