package mission

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/samber/lo"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
	"github.com/galaxy-empire-team/event-manager/pkg/notifications"
)

func (s *Service) handleAttack(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) error {
	attackerPlanet, err := storage.GetPlanetInfoByID(ctx, missionEvent.PlanetFrom)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetInfoByID(): %w", err)
	}

	attackerFleet := missionEvent.Fleet

	defenderPlanet, err := storage.GetPlanetInfoByCoordinates(ctx, missionEvent.PlanetTo)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetInfoByCoordinates(): %w", err)
	}

	defenderFleet, err := storage.GetPlanetFleetForUpdate(ctx, defenderPlanet.ID)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetFleetForUpdate(): %w", err)
	}

	// --- calculate attack result ---
	attackResult, err := s.calcAttackResult(ctx, attackSetup{
		attackerID:    attackerPlanet.UserID,
		defenderID:    defenderPlanet.UserID,
		attackerFleet: attackerFleet,
		defenderFleet: defenderFleet,
	}, storage)
	if err != nil {
		return fmt.Errorf("s.calcAttackResult(): %w", err)
	}

	var gainedResources models.Resources
	if attackResult.attackerWins {
		gainedResources, err = s.stealResources(ctx, defenderPlanet, attackerFleet, storage)
		if err != nil {
			return fmt.Errorf("s.stealResources(): %w", err)
		}
	}

	// --- save attack result ---
	err = storage.CreateMissionEvent(ctx, models.MissionEvent{
		MissionID:   missionEvent.MissionID,
		UserID:      missionEvent.UserID,
		PlanetFrom:  defenderPlanet.ID,
		PlanetTo:    attackerPlanet.Coordinates,
		Fleet:       attackResult.attackerFleetLeft,
		Cargo:       gainedResources,
		IsReturning: true,
		StartedAt:   missionEvent.FinishedAt,
		FinishedAt:  missionEvent.FinishedAt.Add(missionEvent.FinishedAt.Sub(missionEvent.StartedAt)),
	})
	if err != nil {
		return fmt.Errorf("storage.CreateMissionEvent(): %w", err)
	}

	err = storage.SetPlanetFleet(ctx, defenderPlanet.ID, attackResult.defenderFleetLeft)
	if err != nil {
		return fmt.Errorf("storage.SetPlanetFleet(%s): %w", defenderPlanet.ID.String(), err)
	}

	// --- create attack notifications for both attacker and defender ---
	notificationMsg := notifications.AttackV1{
		AttackerWins: attackResult.attackerWins,
		Cargo: notifications.Resources{
			Metal:   gainedResources.Metal,
			Crystal: gainedResources.Crystal,
			Gas:     gainedResources.Gas,
		},
		Attacker: notifications.AttackInfo{
			Login: attackerPlanet.UserLogin,
			Planet: notifications.Coordinates{
				X: attackerPlanet.Coordinates.X,
				Y: attackerPlanet.Coordinates.Y,
				Z: attackerPlanet.Coordinates.Z,
			},
			Fleet: prepareAttackFleetNotification(attackerFleet, attackResult.attackerFleetLeft),
		},
		Defender: notifications.AttackInfo{
			Login: defenderPlanet.UserLogin,
			Planet: notifications.Coordinates{
				X: defenderPlanet.Coordinates.X,
				Y: defenderPlanet.Coordinates.Y,
				Z: defenderPlanet.Coordinates.Z,
			},
			Fleet: prepareAttackFleetNotification(defenderFleet, attackResult.defenderFleetLeft),
		},
	}

	users := userIDPair{
		Attacker: attackerPlanet.UserID,
		Defender: defenderPlanet.UserID,
	}

	err = s.createAttackNotificationEvent(ctx, users, notificationMsg, storage)
	if err != nil {
		return fmt.Errorf("s.createAttackNotificationEvent(): %w", err)
	}

	return nil
}

func (s *Service) stealResources(ctx context.Context, defenderPlanet models.Planet, attackerFleet []models.FleetUnit, storage TxStorages) (models.Resources, error) {
	if err := s.bridgeAPIClient.UpdatePlanetResources(ctx, defenderPlanet.UserID, defenderPlanet.ID, time.Now().UTC()); err != nil {
		return models.Resources{}, fmt.Errorf("bridgeAPIClient.UpdatePlanetResources(): %w", err)
	}

	resources, err := storage.GetResourcesForUpdate(ctx, defenderPlanet.ID)
	if err != nil {
		return models.Resources{}, fmt.Errorf("storage.GetResourcesForUpdate(): %w", err)
	}

	fleetCapactity, err := s.calcFleetCapacity(attackerFleet)
	if err != nil {
		return models.Resources{}, fmt.Errorf("s.calcFleetCapacity(): %w", err)
	}

	fillResult := s.fillFleetCargo(resources, fleetCapactity)

	err = storage.SetResources(ctx, defenderPlanet.ID, fillResult.leftOnPlanet)
	if err != nil {
		return models.Resources{}, fmt.Errorf("storage.SetResources(): %w", err)
	}

	return fillResult.gained, nil
}

func (s *Service) createAttackNotificationEvent(ctx context.Context, users userIDPair, attackNotification notifications.AttackV1, storage TxStorages) error {
	nID, err := s.registry.GetNotificationIDByType(consts.NotificationTypeAttack)
	if err != nil {
		return fmt.Errorf("s.registry.GetNotificationIDByType(): %w", err)
	}

	msg, err := json.Marshal(attackNotification)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	const attackNotificationVersion = 1
	notificationEvents := []models.NotificationEvent{
		{
			UserID:         users.Attacker,
			Version:        attackNotificationVersion,
			NotificationID: nID,
			Data:           msg,
		},
		{
			UserID:         users.Defender,
			Version:        attackNotificationVersion,
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

func prepareAttackFleetNotification(fleetBefore []models.FleetUnit, fleetAfter []models.FleetUnit) []notifications.AttackFleetUnit {
	var notificationAttackerFleet []notifications.AttackFleetUnit

	for _, unit := range fleetBefore {
		result, _ := lo.Find(fleetAfter, func(x models.FleetUnit) bool {
			return x.ID == unit.ID
		})

		notificationAttackerFleet = append(notificationAttackerFleet, notifications.AttackFleetUnit{
			ID:          unit.ID,
			CountBefore: unit.Count,
			CountAfter:  result.Count,
		})
	}

	return notificationAttackerFleet
}

func (s *Service) calcFleetCapacity(fleet []models.FleetUnit) (uint64, error) {
	var capacity uint64

	for _, unit := range fleet {
		fleetStats, err := s.registry.GetFleetUnitStatsByID(unit.ID)
		if err != nil {
			return 0, fmt.Errorf("s.registry.GetFleetUnitStatsByID(): %w", err)
		}

		capacity += unit.Count * fleetStats.CargoCapacity
	}

	return capacity, nil
}
