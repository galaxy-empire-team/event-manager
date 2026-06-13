package mission

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
	"github.com/galaxy-empire-team/event-manager/pkg/notifications"
)

func (s *Service) handleSpy(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) error {
	if len(missionEvent.Fleet) != 1 {
		s.logger.Warn("handleSpy: invalid fleet size", zap.Any("missionEvent", missionEvent))
		return nil
	}

	targetPlanet, err := storage.GetPlanetInfoByCoordinates(ctx, missionEvent.PlanetTo)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetInfoByCoordinates(): %w", err)
	}

	if err := s.bridgeAPIClient.UpdatePlanetResources(ctx, missionEvent.UserID, targetPlanet.ID, missionEvent.FinishedAt); err != nil {
		return fmt.Errorf("bridgeAPIClient.UpdatePlanetResources(): %w", err)
	}

	spyChances, err := s.calcSpyChance(ctx, missionEvent.UserID, missionEvent.Fleet[0].Count, storage)
	if err != nil {
		return fmt.Errorf("s.calcSpyChance(): %w", err)
	}

	var planetResources models.Resources
	if spyChances.spyResources {
		planetResources, err = storage.GetResources(ctx, targetPlanet.ID)
		if err != nil {
			return fmt.Errorf("storage.GetResources(): %w", err)
		}
	}

	var buildings []consts.BuildingID
	if spyChances.spyBuildings {
		buildings, err = storage.GetBuildings(ctx, targetPlanet.ID)
		if err != nil {
			return fmt.Errorf("storage.GetBuildings(): %w", err)
		}
	}

	var fleet []models.FleetUnit
	if spyChances.spyFleet {
		fleet, err = storage.GetPlanetFleetForUpdate(ctx, targetPlanet.ID)
		if err != nil {
			return fmt.Errorf("storage.GetPlanetFleetForUpdate(): %w", err)
		}
	}

	var researches []consts.ResearchID
	if spyChances.spyResearches {
		researches, err = storage.GetUserResearches(ctx, targetPlanet.UserID)
		if err != nil {
			return fmt.Errorf("storage.GetUserResearchTypes(): %w", err)
		}
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

	spyVictimNotification := notifications.SpyV1{
		IsSpy: false,
		Login: planetFrom.UserLogin,
		Coordinates: notifications.Coordinates{
			X: planetFrom.Coordinates.X,
			Y: planetFrom.Coordinates.Y,
			Z: planetFrom.Coordinates.Z,
		},
	}

	spyInitiatorNotification := notifications.SpyV1{
		IsSpy: true,
		Login: targetPlanet.UserLogin,
		Coordinates: notifications.Coordinates{
			X: targetPlanet.Coordinates.X,
			Y: targetPlanet.Coordinates.Y,
			Z: targetPlanet.Coordinates.Z,
		},
		Resources: notifications.Resources{
			Metal:   planetResources.Metal,
			Crystal: planetResources.Crystal,
			Gas:     planetResources.Gas,
		},
		Buildings: buildings,
		Fleet: lo.Map(fleet, func(f models.FleetUnit, _ int) notifications.FleetUnit {
			return notifications.FleetUnit{
				ID:    f.ID,
				Count: f.Count,
			}
		}),
		Researches: researches,
		Result: notifications.SpyResult{
			ResourcesGot:  spyChances.spyResources,
			BuildingsGot:  spyChances.spyBuildings,
			FleetGot:      spyChances.spyFleet,
			ResearchesGot: spyChances.spyResearches,
		},
	}

	users := userIDPair{
		Attacker: missionEvent.UserID,
		Defender: targetPlanet.UserID,
	}

	err = s.createSpyNotificationEvent(ctx, users, spyInitiatorNotification, spyVictimNotification, storage)
	if err != nil {
		return fmt.Errorf("s.createSpyNotificationEvent(): %w", err)
	}

	return nil
}

func (s *Service) createSpyNotificationEvent(ctx context.Context, users userIDPair, spyInfo notifications.SpyV1, spyResult notifications.SpyV1, storage TxStorages) error {
	nID, err := s.registry.GetNotificationIDByType(consts.NotificationTypeSpy)
	if err != nil {
		return fmt.Errorf("s.registry.GetNotificationIDByType(): %w", err)
	}

	attackerMsg, err := json.Marshal(spyInfo)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	defenderMsg, err := json.Marshal(spyResult)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	const spyNotificationVersion = 1
	notificationEvents := []models.NotificationEvent{
		{
			UserID:         users.Attacker,
			Version:        spyNotificationVersion,
			NotificationID: nID,
			Data:           attackerMsg,
		},
		{
			UserID:         users.Defender,
			Version:        spyNotificationVersion,
			NotificationID: nID,
			Data:           defenderMsg,
		},
	}

	err = storage.SaveNotificationEvents(ctx, notificationEvents)
	if err != nil {
		return fmt.Errorf("storage.SaveNotificationEvents(): %w", err)
	}

	return nil
}
