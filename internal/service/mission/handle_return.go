package mission

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const StatusFinished = "Finished"

func (s *Service) returnMission(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) error {
	planetInfo, err := storage.GetPlanetInfoByCoordinates(ctx, missionEvent.PlanetTo)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetInfoByCoordinates(): %w", err)
	}

	err = storage.UpsertFleet(ctx, planetInfo.ID, missionEvent.Fleet)
	if err != nil {
		return fmt.Errorf("storage.UpsertFleet(): %w", err)
	}

	// --- create return notification ---
	mType, err := s.registryProvider.GetMissionTypeByID(missionEvent.MissionID)
	if err != nil {
		return fmt.Errorf("s.registryProvider.GetMissionTypeByID(): %w", err)
	}

	notificationMsg := returnNotification{
		MissionType: mType,
		Msg:         StatusFinished,
	}

	err = s.createReturnNotification(ctx, missionEvent.UserID, notificationMsg, storage)
	if err != nil {
		return fmt.Errorf("s.createReturnNotification(): %w", err)
	}

	return nil
}

type returnNotification struct {
	MissionType consts.MissionType `json:"mission_type"`
	Msg         string             `json:"status"`
}

func (s *Service) createReturnNotification(ctx context.Context, userID uuid.UUID, returnNotification returnNotification, storage TxStorages) error {
	nID, err := s.registryProvider.GetNotificationIDByType(consts.NotificationTypeReturn)
	if err != nil {
		return fmt.Errorf("s.registryProvider.GetNotificationIDByType(): %w", err)
	}

	msg, err := json.Marshal(returnNotification)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	err = storage.SaveNotificationEvents(ctx, []models.NotificationEvent{{
		UserID:         userID,
		NotificationID: nID,
		Data:           msg,
	}})
	if err != nil {
		return fmt.Errorf("storage.SaveNotification(): %w", err)
	}

	return nil
}
