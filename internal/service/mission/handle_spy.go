package mission

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
	"github.com/galaxy-empire-team/event-manager/internal/models"
	"github.com/galaxy-empire-team/event-manager/pkg/notifications"
)

type spyResult struct {
	TargetPlanet models.Planet
	Resources    models.Resources
	Buildings    []consts.BuildingID
	Fleet        []models.FleetUnit
	Researches   []consts.ResearchID
	SpyChances   spyChancesResult
}

func (s *Service) handleSpy(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) error {
	if len(missionEvent.Fleet) != 1 {
		s.logger.Warn("handleSpy: invalid fleet size", zap.Any("missionEvent", missionEvent))
		return nil
	}

	var (
		spyResult spyResult
		err       error
	)
	if s.isPlanetNPC(missionEvent.PlanetTo.Z) {
		spyResult, err = s.getSpyInfoForNPC(missionEvent)
		if err != nil {
			return fmt.Errorf("s.getSpyInfoForNPC(): %w", err)
		}
	} else {
		spyResult, err = s.getSpyInfoForUser(ctx, missionEvent, storage)
		if err != nil {
			return fmt.Errorf("s.getSpyInfoForUser(): %w", err)
		}
	}

	planetFrom, err := storage.GetPlanetInfoByID(ctx, missionEvent.PlanetFrom)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetInfoByID(): %w", err)
	}

	err = storage.CreateMissionEvent(ctx, models.MissionEvent{
		MissionID:   missionEvent.MissionID,
		UserID:      missionEvent.UserID,
		PlanetFrom:  spyResult.TargetPlanet.ID,
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
		Login: spyResult.TargetPlanet.UserLogin,
		Coordinates: notifications.Coordinates{
			X: spyResult.TargetPlanet.Coordinates.X,
			Y: spyResult.TargetPlanet.Coordinates.Y,
			Z: spyResult.TargetPlanet.Coordinates.Z,
		},
		Resources: notifications.Resources{
			Metal:   spyResult.Resources.Metal,
			Crystal: spyResult.Resources.Crystal,
			Gas:     spyResult.Resources.Gas,
		},
		Buildings: spyResult.Buildings,
		Fleet: lo.Map(spyResult.Fleet, func(f models.FleetUnit, _ int) notifications.FleetUnit {
			return notifications.FleetUnit{
				ID:    f.ID,
				Count: f.Count,
			}
		}),
		Researches: spyResult.Researches,
		Result: notifications.SpyResult{
			ResourcesGot:  spyResult.SpyChances.spyResources,
			BuildingsGot:  spyResult.SpyChances.spyBuildings,
			FleetGot:      spyResult.SpyChances.spyFleet,
			ResearchesGot: spyResult.SpyChances.spyResearches,
		},
	}

	users := userIDPair{
		Attacker: missionEvent.UserID,
		Defender: spyResult.TargetPlanet.UserID,
	}

	err = s.createSpyNotificationEvent(ctx, users, spyInitiatorNotification, spyVictimNotification, storage)
	if err != nil {
		return fmt.Errorf("s.createSpyNotificationEvent(): %w", err)
	}

	return nil
}

func (s *Service) getSpyInfoForUser(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) (spyResult, error) {
	targetPlanet, err := storage.GetPlanetInfoByCoordinates(ctx, missionEvent.PlanetTo)
	if err != nil {
		return spyResult{}, fmt.Errorf("storage.GetPlanetInfoByCoordinates(): %w", err)
	}

	if err := s.bridgeAPIClient.UpdatePlanetResources(ctx, missionEvent.UserID, targetPlanet.ID, missionEvent.FinishedAt); err != nil {
		return spyResult{}, fmt.Errorf("bridgeAPIClient.UpdatePlanetResources(): %w", err)
	}

	spyChances, err := s.calcSpyChance(ctx, missionEvent.UserID, missionEvent.Fleet[0].Count, storage)
	if err != nil {
		return spyResult{}, fmt.Errorf("s.calcSpyChance(): %w", err)
	}

	var planetResources models.Resources
	if spyChances.spyResources {
		planetResources, err = storage.GetResources(ctx, targetPlanet.ID)
		if err != nil {
			return spyResult{}, fmt.Errorf("storage.GetResources(): %w", err)
		}
	}

	var buildings []consts.BuildingID
	if spyChances.spyBuildings {
		buildings, err = storage.GetBuildings(ctx, targetPlanet.ID)
		if err != nil {
			return spyResult{}, fmt.Errorf("storage.GetBuildings(): %w", err)
		}
	}

	var fleet []models.FleetUnit
	if spyChances.spyFleet {
		fleet, err = storage.GetPlanetFleetForUpdate(ctx, targetPlanet.ID)
		if err != nil {
			return spyResult{}, fmt.Errorf("storage.GetPlanetFleetForUpdate(): %w", err)
		}
	}

	var researches []consts.ResearchID
	if spyChances.spyResearches {
		researches, err = storage.GetUserResearches(ctx, targetPlanet.UserID)
		if err != nil {
			return spyResult{}, fmt.Errorf("storage.GetUserResearchTypes(): %w", err)
		}
	}

	return spyResult{
		TargetPlanet: targetPlanet,
		Resources:    planetResources,
		Buildings:    buildings,
		Fleet:        fleet,
		Researches:   researches,
		SpyChances:   spyChances,
	}, nil
}

func (s *Service) getSpyInfoForNPC(missionEvent models.MissionEvent) (spyResult, error) {
	npcStats, err := s.registry.GetNPCStatsByPosition(missionEvent.PlanetTo.Z)
	if err != nil {
		return spyResult{}, fmt.Errorf("s.registry.GetNPCStatsByPosition(): %w", err)
	}

	return spyResult{
		TargetPlanet: models.Planet{
			ID:          uuid.Nil,
			UserID:      uuid.Nil,
			UserLogin:   s.getNPCLoginByCoordinates(missionEvent.PlanetTo),
			Coordinates: missionEvent.PlanetTo,
		},
		Resources: models.Resources{
			Metal:   npcStats.Resources.Metal,
			Crystal: npcStats.Resources.Crystal,
			Gas:     npcStats.Resources.Gas,
		},
		Buildings: nil,
		Fleet: lo.Map(npcStats.Fleet, func(fleetCount registry.FleetUnitCount, _ int) models.FleetUnit {
			return models.FleetUnit{
				ID:    fleetCount.ID,
				Count: fleetCount.Count,
			}
		}),
		Researches: lo.Map(npcStats.Researches, func(researchID consts.ResearchID, _ int) consts.ResearchID {
			return researchID
		}),
		SpyChances: spyChancesResult{
			spyResources:  true,
			spyBuildings:  true,
			spyFleet:      true,
			spyResearches: true,
		},
	}, nil
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

	const spyNotificationVersion = 1
	notificationEvents := []models.NotificationEvent{
		{
			UserID:         users.Attacker,
			Version:        spyNotificationVersion,
			NotificationID: nID,
			Data:           attackerMsg,
		},
	}

	if users.Defender != uuid.Nil {
		defenderMsg, err := json.Marshal(spyResult)
		if err != nil {
			return fmt.Errorf("json.Marshal(): %w", err)
		}

		notificationEvents = append(notificationEvents, models.NotificationEvent{
			UserID:         users.Defender,
			Version:        spyNotificationVersion,
			NotificationID: nID,
			Data:           defenderMsg,
		})
	}

	err = storage.SaveNotificationEvents(ctx, notificationEvents)
	if err != nil {
		return fmt.Errorf("storage.SaveNotificationEvents(): %w", err)
	}

	return nil
}
