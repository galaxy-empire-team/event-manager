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

func (s *Service) handleRecycle(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) error {
	planetID, err := storage.GetPlanetIDByCoordinates(ctx, missionEvent.PlanetTo)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetIDByCoordinates(): %w", err)
	}

	debris, err := storage.GetDebrisForUpdate(ctx, planetID)
	if err != nil {
		return fmt.Errorf("storage.GetDebrisForUpdate(): %w", err)
	}

	fleetCapactity, err := s.calcFleetCapacity(missionEvent.Fleet)
	if err != nil {
		return fmt.Errorf("s.calcFleetCapacity(): %w", err)
	}

	fillResult := s.fillFleetCargo(debris, fleetCapactity)

	err = storage.SetDebris(ctx, planetID, fillResult.leftOnPlanet)
	if err != nil {
		return fmt.Errorf("storage.SetDebris(): %w", err)
	}

	planetCoordinates, err := storage.GetPlanetCoordinatesByID(ctx, missionEvent.PlanetFrom)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetCoordinatesByID(): %w", err)
	}

	err = storage.CreateMissionEvent(ctx, models.MissionEvent{
		MissionID:   missionEvent.MissionID,
		UserID:      missionEvent.UserID,
		PlanetFrom:  uuid.Nil,
		PlanetTo:    planetCoordinates,
		Fleet:       missionEvent.Fleet,
		Cargo:       fillResult.gained,
		IsReturning: true,
		StartedAt:   missionEvent.FinishedAt,
		FinishedAt:  missionEvent.FinishedAt.Add(missionEvent.FinishedAt.Sub(missionEvent.StartedAt)),
	})
	if err != nil {
		return fmt.Errorf("storage.CreateMissionEvent(): %w", err)
	}

	recycleNotification := notifications.RecycleV1{
		Coordinates: notifications.Coordinates{
			X: planetCoordinates.X,
			Y: planetCoordinates.Y,
			Z: planetCoordinates.Z,
		},
		Resources: notifications.Resources{
			Metal:   fillResult.gained.Metal,
			Crystal: fillResult.gained.Crystal,
			Gas:     fillResult.gained.Gas,
		},
	}

	err = s.createRecycleNotificationEvent(ctx, missionEvent.UserID, recycleNotification, storage)
	if err != nil {
		return fmt.Errorf("s.createRecycleNotificationEvent(): %w", err)
	}

	return nil
}

func (s *Service) createRecycleNotificationEvent(ctx context.Context, userID uuid.UUID, recycleResult notifications.RecycleV1, storage TxStorages) error {
	nID, err := s.registry.GetNotificationIDByType(consts.NotificationTypeRecycle)
	if err != nil {
		return fmt.Errorf("s.registry.GetNotificationIDByType(): %w", err)
	}

	msg, err := json.Marshal(recycleResult)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	const recycleNotificationVersion = 1
	notificationEvents := []models.NotificationEvent{
		{
			UserID:         userID,
			Version:        recycleNotificationVersion,
			NotificationID: nID,
			Data:           msg,
		},
	}

	err = storage.SaveNotificationEvents(ctx, notificationEvents)
	if err != nil {
		return fmt.Errorf("storage.SaveNotificationEvents(): %w", err)
	}

	return nil
}
