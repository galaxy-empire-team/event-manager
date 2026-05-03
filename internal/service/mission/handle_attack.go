package mission

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
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

	// --- save attack result ---
	err = storage.CreateMissionEvent(ctx, models.MissionEvent{
		MissionID:   missionEvent.MissionID,
		UserID:      missionEvent.UserID,
		PlanetFrom:  defenderPlanet.ID,
		PlanetTo:    attackerPlanet.Coordinates,
		Fleet:       attackResult.attackerFleetLeft,
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
	notificationMsg := attackNotification{
		AttackerWins: attackResult.attackerWins,
		Cargo: resources{
			Metal:   0,
			Crystal: 0,
			Gas:     0,
		},
		Attacker: attackInfo{
			Login: attackerPlanet.UserLogin,
			Planet: attackCoordinates{
				X: attackerPlanet.Coordinates.X,
				Y: attackerPlanet.Coordinates.Y,
				Z: attackerPlanet.Coordinates.Z,
			},
			Fleet: prepareFleetForNotification(attackerFleet, attackResult.attackerFleetLeft),
		},
		Defender: attackInfo{
			Login: defenderPlanet.UserLogin,
			Planet: attackCoordinates{
				X: defenderPlanet.Coordinates.X,
				Y: defenderPlanet.Coordinates.Y,
				Z: defenderPlanet.Coordinates.Z,
			},
			Fleet: prepareFleetForNotification(defenderFleet, attackResult.defenderFleetLeft),
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

type attackNotification struct {
	AttackerWins bool       `json:"attackerWins"`
	Cargo        resources  `json:"cargo"`
	Attacker     attackInfo `json:"attacker"`
	Defender     attackInfo `json:"defender"`
}

type attackInfo struct {
	Login  string            `json:"login"`
	Planet attackCoordinates `json:"planet"`
	Fleet  []attackFleetUnit `json:"fleet"`
}

type attackFleetUnit struct {
	ID          consts.FleetUnitID `json:"id"`
	CountBefore uint64             `json:"countBefore"`
	CountAfter  uint64             `json:"countAfter"`
}

type resources struct {
	Metal   uint64 `json:"metal"`
	Crystal uint64 `json:"crystal"`
	Gas     uint64 `json:"gas"`
}

type attackCoordinates struct {
	X consts.PlanetPositionX `json:"x"`
	Y consts.PlanetPositionY `json:"y"`
	Z consts.PlanetPositionZ `json:"z"`
}

func (s *Service) createAttackNotificationEvent(ctx context.Context, users userIDPair, attackNotification attackNotification, storage TxStorages) error {
	nID, err := s.registry.GetNotificationIDByType(consts.NotificationTypeColonize)
	if err != nil {
		return fmt.Errorf("s.registry.GetNotificationIDByType(): %w", err)
	}

	msg, err := json.Marshal(attackNotification)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	notificationEvents := []models.NotificationEvent{
		{
			UserID:         users.Attacker,
			NotificationID: nID,
			Data:           msg,
		},
		{
			UserID:         users.Defender,
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

func prepareFleetForNotification(fleetBefore []models.FleetUnit, fleetAfter []models.FleetUnit) []attackFleetUnit {
	var notificationAttackerFleet []attackFleetUnit

	for _, unit := range fleetBefore {
		result, _ := lo.Find(fleetAfter, func(x models.FleetUnit) bool {
			return x.ID == unit.ID
		})

		notificationAttackerFleet = append(notificationAttackerFleet, attackFleetUnit{
			ID:          unit.ID,
			CountBefore: unit.Count,
			CountAfter:  result.Count,
		})
	}

	return notificationAttackerFleet
}
