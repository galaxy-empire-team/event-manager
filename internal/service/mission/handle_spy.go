package mission

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *Service) handleSpy(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) error {
	targetPlanet, err := storage.GetPlanetInfoByCoordinates(ctx, missionEvent.PlanetTo)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetInfoByCoordinates(): %w", err)
	}

	if err := s.bridgeAPIClient.UpdatePlanetResources(ctx, missionEvent.UserID, targetPlanet.ID, missionEvent.FinishedAt); err != nil {
		return fmt.Errorf("bridgeAPIClient.UpdatePlanetResources(): %w", err)
	}

	planetResources, err := storage.GetResources(ctx, targetPlanet.ID)
	if err != nil {
		return fmt.Errorf("storage.GetResources(): %w", err)
	}

	buildingIDs, err := storage.GetAllBuildings(ctx, targetPlanet.ID)
	if err != nil {
		return fmt.Errorf("storage.GetAllBuildings(): %w", err)
	}

	planetFrom, err := storage.GetPlanetInfoByID(ctx, missionEvent.PlanetFrom)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetInfoByID(): %w", err)
	}

	err = storage.CreateMissionEvent(ctx, models.MissionEvent{
		MissionID:   missionEvent.MissionID,
		UserID:      missionEvent.UserID,
		PlanetFrom:  targetPlanet.ID,
		PlanetTo:    planetFrom.Coordinates,
		Fleet:       missionEvent.Fleet,
		IsReturning: true,
		StartedAt:   missionEvent.FinishedAt,
		FinishedAt:  missionEvent.FinishedAt.Add(missionEvent.FinishedAt.Sub(missionEvent.StartedAt)),
	})
	if err != nil {
		return fmt.Errorf("storage.CreateMissionEvent(): %w", err)
	}

	spyNotification := spyNotification{
		Attacker: attackerSpyNotification{
			UserLogin: targetPlanet.UserLogin,
			PlanetTo: coordinates{
				X: targetPlanet.Coordinates.X,
				Y: targetPlanet.Coordinates.Y,
				Z: targetPlanet.Coordinates.Z,
			},
			Resources: planetResources,
			Buildings: buildingIDs,
		},
		Defender: defenderSpyNotification{
			UserLogin: planetFrom.UserLogin,
			PlanetFrom: coordinates{
				X: planetFrom.Coordinates.X,
				Y: planetFrom.Coordinates.Y,
				Z: planetFrom.Coordinates.Z,
			},
		},
	}

	users := userIDPair{
		Attacker: missionEvent.UserID,
		Defender: targetPlanet.UserID,
	}

	err = s.createSpyNotificationEvent(ctx, users, spyNotification, storage)
	if err != nil {
		return fmt.Errorf("s.createSpyNotificationEvent(): %w", err)
	}

	return nil
}

type spyNotification struct {
	Attacker attackerSpyNotification `json:"attacker"`
	Defender defenderSpyNotification `json:"defender"`
}

type attackerSpyNotification struct {
	UserLogin string              `json:"userLogin"`
	PlanetTo  coordinates         `json:"planetTo"`
	Resources models.Resources    `json:"resources"`
	Buildings []consts.BuildingID `json:"buildings"`
}

type coordinates struct {
	X consts.PlanetPositionX `json:"x"`
	Y consts.PlanetPositionY `json:"y"`
	Z consts.PlanetPositionZ `json:"z"`
}

type defenderSpyNotification struct {
	UserLogin  string      `json:"userLogin"`
	PlanetFrom coordinates `json:"planetFrom"`
}

func (s *Service) createSpyNotificationEvent(ctx context.Context, users userIDPair, spyNotification spyNotification, storage TxStorages) error {
	nID, err := s.registry.GetNotificationIDByType(consts.NotificationTypeSpy)
	if err != nil {
		return fmt.Errorf("s.registry.GetNotificationIDByType(): %w", err)
	}

	attackerMsg, err := json.Marshal(spyNotification.Attacker)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	deffenderMsg, err := json.Marshal(spyNotification.Defender)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	notificationEvents := []models.NotificationEvent{
		{
			UserID:         users.Attacker,
			NotificationID: nID,
			Data:           attackerMsg,
		},
		{
			UserID:         users.Defender,
			NotificationID: nID,
			Data:           deffenderMsg,
		},
	}
	err = storage.SaveNotificationEvents(ctx, notificationEvents)
	if err != nil {
		return fmt.Errorf("storage.SaveNotificationEvents(): %w", err)
	}

	return nil
}
