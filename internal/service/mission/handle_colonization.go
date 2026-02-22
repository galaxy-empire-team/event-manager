package mission

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *Service) handleColonization(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) error {
	var colonizationErr string

	colonized, err := storage.ColonizePlanet(ctx, missionEvent)
	if err != nil {
		return fmt.Errorf("storage.ColonizePlanet(): %w", err)
	}
	if !colonized {
		colonizationErr = "Planet is already colonized"
	}

	// --- create colonization notification ---
	notificationMsg := colonizationNotification{
		UserID: missionEvent.UserID,
		Planet: colonizationCoordinates{
			X: missionEvent.PlanetTo.X,
			Y: missionEvent.PlanetTo.Y,
			Z: missionEvent.PlanetTo.Z,
		},
		Err: colonizationErr,
	}

	err = s.createColonizationNotificationEvent(ctx, notificationMsg, storage)
	if err != nil {
		return fmt.Errorf("storage.CreateColonizationNotification(): %w", err)
	}

	return nil
}

type colonizationNotification struct {
	UserID uuid.UUID               `json:"user_id"`
	Planet colonizationCoordinates `json:"planet"`
	Err    string                  `json:"err"`
}

type colonizationCoordinates struct {
	X consts.PlanetPositionX `json:"x"`
	Y consts.PlanetPositionY `json:"y"`
	Z consts.PlanetPositionZ `json:"z"`
}

func (s *Service) createColonizationNotificationEvent(ctx context.Context, colonizationNotification colonizationNotification, storage TxStorages) error {
	nID, err := s.registryProvider.GetNotificationIDByType(consts.NotificationTypeColonize)
	if err != nil {
		return fmt.Errorf("s.registryProvider.GetNotificationIDByType(): %w", err)
	}

	msg, err := json.Marshal(colonizationNotification)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	err = storage.SaveNotificationEvents(ctx, []models.NotificationEvent{{
		UserID:         colonizationNotification.UserID,
		NotificationID: nID,
		Data:           msg,
	},
	})
	if err != nil {
		return fmt.Errorf("storage.SaveNotification(): %w", err)
	}

	return nil
}
