package mission

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *Service) handleColonization(ctx context.Context, colonizationEvent models.MissionEvent, storage TxStorages) error {
	notificationMsg := colonizationNotification{
		Planet: colonizationCoordinates{
			X: colonizationEvent.PlanetTo.X,
			Y: colonizationEvent.PlanetTo.Y,
			Z: colonizationEvent.PlanetTo.Z,
		},
	}

	err := s.bridgeAPIClient.ColonizePlanet(ctx, colonizationEvent.UserID, colonizationEvent)
	if err != nil {
		if !errors.Is(err, models.ErrPlanetCoordinatesAlreadyTaken) {
			return fmt.Errorf("s.bridgeAPIClient.ColonizePlanet(): %w", err)
		}

		notificationMsg.Err = "Planet coordinates already taken"
	}

	err = s.createColonizationNotificationEvent(ctx, colonizationEvent.UserID, notificationMsg, storage)
	if err != nil {
		return fmt.Errorf("storage.CreateColonizationNotification(): %w", err)
	}

	return nil
}

type colonizationNotification struct {
	Planet colonizationCoordinates `json:"planet"`
	Err    string                  `json:"error,omitempty"`
}

type colonizationCoordinates struct {
	X consts.PlanetPositionX `json:"x"`
	Y consts.PlanetPositionY `json:"y"`
	Z consts.PlanetPositionZ `json:"z"`
}

func (s *Service) createColonizationNotificationEvent(ctx context.Context, userID uuid.UUID, colonizationNotification colonizationNotification, storage TxStorages) error {
	nID, err := s.registry.GetNotificationIDByType(consts.NotificationTypeColonize)
	if err != nil {
		return fmt.Errorf("s.registry.GetNotificationIDByType(): %w", err)
	}

	msg, err := json.Marshal(colonizationNotification)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	err = storage.SaveNotificationEvents(ctx, []models.NotificationEvent{{
		UserID:         userID,
		NotificationID: nID,
		Data:           msg,
	},
	})
	if err != nil {
		return fmt.Errorf("storage.SaveNotification(): %w", err)
	}

	return nil
}
