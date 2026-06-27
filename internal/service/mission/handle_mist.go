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

const (
	// Mist zone reward chances
	resourceGainChance = 60
	fleetGainChance    = 30
	boostGainChance    = 5
	matterGainChance   = 5
	maxMistChance      = 100
)

type RewardType string

const (
	RewardTypeResource RewardType = "resource"
	RewardTypeFleet    RewardType = "fleet"
	RewardTypeBoost    RewardType = "boost"
	RewardTypeMatter   RewardType = "matter"
)

func (s *Service) handleMist(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) error {
	attackerPlanet, err := storage.GetPlanetInfoByID(ctx, missionEvent.PlanetFrom)
	if err != nil {
		return fmt.Errorf("storage.GetPlanetInfoByID(): %w", err)
	}

	var (
		fleetReward   models.FleetUnit
		rewardTypeStr notifications.RewardType
		resources     models.Resources
	)
	switch s.getMistReward() {
	case RewardTypeResource:
		rewardTypeStr = notifications.RewardTypeResource
		resources, err = s.calcResourceMistReward(missionEvent.Fleet)
		if err != nil {
			return fmt.Errorf("s.calcResourceMistReward(): %w", err)
		}

	case RewardTypeFleet:
		rewardTypeStr = notifications.RewardTypeFleet
		fleetReward, err = s.calcFleetMistReward(missionEvent.Fleet)
		if err != nil {
			return fmt.Errorf("s.calcFleetMistReward(): %w", err)
		}

	case RewardTypeBoost:
		rewardTypeStr = notifications.RewardTypeBoost
		boostReward, err := s.calcBoostMistReward(missionEvent.Fleet)
		if err != nil {
			return fmt.Errorf("s.calcBoostMistReward(): %w", err)
		}

		resources.Boost = boostReward

	case RewardTypeMatter:
		rewardTypeStr = notifications.RewardTypeMatter
		matterReward, err := s.calcMatterMistReward(missionEvent.Fleet)
		if err != nil {
			return fmt.Errorf("s.calcMatterMistReward(): %w", err)
		}

		resources.Matter = matterReward
	}

	err = storage.CreateMissionEvent(ctx, models.MissionEvent{
		MissionID:   missionEvent.MissionID,
		UserID:      missionEvent.UserID,
		PlanetFrom:  uuid.Nil,
		PlanetTo:    attackerPlanet.Coordinates,
		Fleet:       s.addFleetMistReward(missionEvent.Fleet, fleetReward),
		Cargo:       resources,
		IsReturning: true,
		StartedAt:   missionEvent.FinishedAt,
		FinishedAt:  missionEvent.FinishedAt.Add(missionEvent.FinishedAt.Sub(missionEvent.StartedAt)),
	})
	if err != nil {
		return fmt.Errorf("storage.CreateMissionEvent(): %w", err)
	}

	// --- create mist notification ---
	notificationMsg := notifications.MistV1{
		RewardType: rewardTypeStr,
		Reward: notifications.Reward{
			Resource: notifications.MistResourceReward{
				Metal:   resources.Metal,
				Crystal: resources.Crystal,
				Gas:     resources.Gas,
			},
			Fleet: notifications.MistFleetReward{
				ID:    fleetReward.ID,
				Count: fleetReward.Count,
			},
			Boost: notifications.MistBoostReward{
				Count: resources.Boost.Count,
				ID:    resources.Boost.ID,
			},
			Matter: notifications.MistMatterReward{
				Count: resources.Matter,
			},
		},
	}

	err = s.createMistNotificationEvent(ctx, missionEvent.UserID, notificationMsg, storage)
	if err != nil {
		return fmt.Errorf("s.createAttackNotificationEvent(): %w", err)
	}

	return nil
}

func (s *Service) createMistNotificationEvent(ctx context.Context, userID uuid.UUID, mistNotification notifications.MistV1, storage TxStorages) error {
	nID, err := s.registry.GetNotificationIDByType(consts.NotificationTypeMist)
	if err != nil {
		return fmt.Errorf("s.registry.GetNotificationIDByType(): %w", err)
	}

	msg, err := json.Marshal(mistNotification)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	const mistNotificationVersion = 1
	notificationEvents := []models.NotificationEvent{
		{
			UserID:         userID,
			Version:        mistNotificationVersion,
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

func (s *Service) addFleetMistReward(fleet []models.FleetUnit, reward models.FleetUnit) []models.FleetUnit {
	var found bool
	for _, unit := range fleet {
		if unit.ID == reward.ID {
			unit.Count += reward.Count
			found = true
			break
		}
	}

	if !found {
		fleet = append(fleet, reward)
	}

	return fleet
}

func (s *Service) getMistReward() RewardType {
	chance := s.randGenerator.Intn(maxMistChance)

	if chance < resourceGainChance {
		return RewardTypeResource
	}

	if chance < resourceGainChance+fleetGainChance {
		return RewardTypeFleet
	}

	if chance < resourceGainChance+fleetGainChance+boostGainChance {
		return RewardTypeBoost
	}

	return RewardTypeMatter
}
