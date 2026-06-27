package mission

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
	"github.com/galaxy-empire-team/event-manager/pkg/notifications"
)

const StatusFinished = "finished"

func (s *Service) handleReturn(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) error {
	planetInfo, err := storage.GetPlanetInfoByCoordinates(ctx, missionEvent.PlanetTo)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetInfoByCoordinates(): %w", err)
	}

	err = storage.AddFleet(ctx, planetInfo.ID, missionEvent.Fleet)
	if err != nil {
		return fmt.Errorf("storage.AddFleet(): %w", err)
	}

	if !missionEvent.Cargo.IsEmptyBase() {
		err = storage.AddResources(ctx, planetInfo.ID, missionEvent.Cargo)
		if err != nil {
			return fmt.Errorf("storage.AddResources(): %w", err)
		}
	}

	if missionEvent.Cargo.Matter > 0 {
		err = storage.AddMatter(ctx, planetInfo.UserID, missionEvent.Cargo.Matter)
		if err != nil {
			return fmt.Errorf("storage.AddMatter(): %w", err)
		}
	}

	if missionEvent.Cargo.Boost.ID != consts.BoostID(0) {
		err = storage.AddBoost(ctx, planetInfo.UserID, missionEvent.Cargo.Boost)
		if err != nil {
			return fmt.Errorf("storage.AddBoost(): %w", err)
		}
	}

	// --- create return notification ---
	notificationMsg := notifications.ReturnV1{
		MissionType: missionEvent.MissionID,
		Status:      StatusFinished,
	}

	err = s.createReturnNotification(ctx, missionEvent.UserID, notificationMsg, storage)
	if err != nil {
		return fmt.Errorf("s.createReturnNotification(): %w", err)
	}

	return nil
}

func (s *Service) createReturnNotification(ctx context.Context, userID uuid.UUID, returnNotification notifications.ReturnV1, storage TxStorages) error {
	nID, err := s.registry.GetNotificationIDByType(consts.NotificationTypeReturn)
	if err != nil {
		return fmt.Errorf("s.registry.GetNotificationIDByType(): %w", err)
	}

	msg, err := json.Marshal(returnNotification)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	const returnNotificationVersion = 1
	err = storage.SaveNotificationEvents(ctx, []models.NotificationEvent{{
		UserID:         userID,
		Version:        returnNotificationVersion,
		NotificationID: nID,
		Data:           msg,
	}})
	if err != nil {
		return fmt.Errorf("storage.SaveNotification(): %w", err)
	}

	return nil
}
