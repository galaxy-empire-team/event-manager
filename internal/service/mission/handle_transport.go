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

func (s *Service) handleTransport(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) error {
	planetInfo, err := storage.GetPlanetInfoByCoordinates(ctx, missionEvent.PlanetTo)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetInfoByCoordinates(): %w", err)
	}

	planetFromCoordinates, err := storage.GetPlanetCoordinatesByID(ctx, missionEvent.PlanetFrom)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetCoordinatesByID(): %w", err)
	}

	err = storage.AddFleet(ctx, planetInfo.ID, missionEvent.Fleet)
	if err != nil {
		return fmt.Errorf("storage.AddFleet(): %w", err)
	}

	err = storage.AddResources(ctx, planetInfo.ID, missionEvent.Cargo)
	if err != nil {
		return fmt.Errorf("storage.AddResources(): %w", err)
	}

	// --- create return notification ---
	notificationMsg := notifications.TransportV1{
		PlanetFrom: notifications.Coordinates{
			X: planetFromCoordinates.X,
			Y: planetFromCoordinates.Y,
			Z: planetFromCoordinates.Z,
		},
		PlanetTo: notifications.Coordinates{
			X: planetInfo.Coordinates.X,
			Y: planetInfo.Coordinates.Y,
			Z: planetInfo.Coordinates.Z,
		},
		Status: StatusFinished,
	}

	err = s.createTransportNotification(ctx, missionEvent.UserID, notificationMsg, storage)
	if err != nil {
		return fmt.Errorf("s.createTransportNotification(): %w", err)
	}

	return nil
}

func (s *Service) createTransportNotification(ctx context.Context, userID uuid.UUID, transportNotification notifications.TransportV1, storage TxStorages) error {
	nID, err := s.registry.GetNotificationIDByType(consts.NotificationTypeTransport)
	if err != nil {
		return fmt.Errorf("s.registry.GetNotificationIDByType(): %w", err)
	}

	msg, err := json.Marshal(transportNotification)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	const transportNotificationVersion = 1
	err = storage.SaveNotificationEvents(ctx, []models.NotificationEvent{{
		UserID:         userID,
		Version:        transportNotificationVersion,
		NotificationID: nID,
		Data:           msg,
	}})
	if err != nil {
		return fmt.Errorf("storage.SaveNotification(): %w", err)
	}

	return nil
}
