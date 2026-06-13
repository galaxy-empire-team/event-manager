package mission

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
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

	attackerResearchBonuses, err := s.getResearchBonuses(ctx, attackerPlanet.UserID, storage)
	if err != nil {
		return fmt.Errorf("getResearchBonuses(attacker): %w", err)
	}

	defender, err := s.getDefenderPlanet(ctx, missionEvent, storage)
	if err != nil {
		return fmt.Errorf("s.getDefenderPlanet(): %w", err)
	}

	// --- calculate attack result ---
	attackResult, err := s.calcAttackResult(attackSetup{
		attackerFleet:      missionEvent.Fleet,
		attackerResearches: attackerResearchBonuses,
		defenderFleet:      defender.fleet,
		defenderResearches: defender.researches,
	})
	if err != nil {
		return fmt.Errorf("s.calcAttackResult(): %w", err)
	}

	var gainedResources models.Resources
	if attackResult.attackerWins {
		if s.isPlanetNPC(missionEvent.PlanetTo.Z) {
			fleetCapactity, err := s.calcFleetCapacity(missionEvent.Fleet)
			if err != nil {
				return fmt.Errorf("s.calcFleetCapacity(): %w", err)
			}

			lootingTechStats, err := s.repository.GetResearchByType(ctx, missionEvent.UserID, consts.ResearchTypeLootingTechnology)
			if err != nil {
				return fmt.Errorf("storage.GetUserResearchesByTypes(): %w", err)
			}

			defender.planet.Resources.Metal = uint64(float64(defender.planet.Resources.Metal) * float64(lootingTechStats.Bonuses.LootingNPCMuliplier))
			defender.planet.Resources.Crystal = uint64(float64(defender.planet.Resources.Crystal) * float64(lootingTechStats.Bonuses.LootingNPCMuliplier))
			defender.planet.Resources.Gas = uint64(float64(defender.planet.Resources.Gas) * float64(lootingTechStats.Bonuses.LootingNPCMuliplier))
			gainedResources = s.fillFleetCargo(defender.planet.Resources, fleetCapactity).gained

			updatedAttackerFleet, err := s.lootNPCFleet(attackResult.attackerFleetLeft, missionEvent.PlanetTo.Z, lootingTechStats.Bonuses.LootingNPCMuliplier)
			if err != nil {
				return fmt.Errorf("s.lootNPCFleet(): %w", err)
			}

			attackResult.attackerFleetLeft = updatedAttackerFleet
		} else {
			gainedResources, err = s.stealResources(ctx, defender.planet, missionEvent.Fleet, storage)
			if err != nil {
				return fmt.Errorf("s.stealResources(): %w", err)
			}
		}
	}

	// --- save attack result ---
	err = storage.CreateMissionEvent(ctx, models.MissionEvent{
		MissionID:   missionEvent.MissionID,
		UserID:      missionEvent.UserID,
		PlanetFrom:  defender.planet.ID,
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

	if !s.isPlanetNPC(missionEvent.PlanetTo.Z) {
		err = storage.SetPlanetFleet(ctx, defender.planet.ID, attackResult.defenderFleetLeft)
		if err != nil {
			return fmt.Errorf("storage.SetPlanetFleet(%s): %w", defender.planet.ID.String(), err)
		}

		debris, err := s.calcDebris(missionEvent.Fleet, attackResult.attackerFleetLeft)
		if err != nil {
			return fmt.Errorf("s.calcDebris(): %w", err)
		}

		err = storage.AddDebris(ctx, attackerPlanet.ID, debris)
		if err != nil {
			return fmt.Errorf("storage.AddDebris(): %w", err)
		}

		if s.isMoonCreated(debris) {
			err = storage.CreateMoon(ctx, defender.planet.ID)
			if err != nil {
				return fmt.Errorf("storage.CreateMoon(): %w", err)
			}
		}
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
			Fleet: prepareAttackFleetNotification(missionEvent.Fleet, attackResult.attackerFleetLeft),
		},
		Defender: notifications.AttackInfo{
			Login: defender.planet.UserLogin,
			Planet: notifications.Coordinates{
				X: missionEvent.PlanetTo.X,
				Y: missionEvent.PlanetTo.Y,
				Z: missionEvent.PlanetTo.Z,
			},
			Fleet: prepareAttackFleetNotification(defender.fleet, attackResult.defenderFleetLeft),
		},
	}

	users := userIDPair{
		Attacker: attackerPlanet.UserID,
		Defender: defender.planet.UserID,
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
	}

	// If planet is not NPC
	if users.Defender != uuid.Nil {
		notificationEvents = append(notificationEvents, models.NotificationEvent{
			UserID:         users.Defender,
			Version:        attackNotificationVersion,
			NotificationID: nID,
			Data:           msg,
		})
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
