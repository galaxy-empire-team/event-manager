package mission

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const (
	// Multiplier depends on both power and capacity
	fleedPowerToGainMaxChanceBoostAmount    = 200_000    // equals 30 anihilators
	fleedCapacityToGainMaxChanceBoostAmount = 10_000_000 // equals 1000 transporters
	maxBoostFleetPowerMultiplier            = 0.5
	maxBoostFleetCapacityMultiplier         = 0.5
)

type boostMistReward struct {
	count uint64
	tier  consts.BoostTier
}

func (s *Service) calcBoostMistReward(fleet []models.FleetUnit) (models.Boost, error) {
	totalFleetPower, err := s.calcBaseFleetPower(fleet)
	if err != nil {
		return models.Boost{}, fmt.Errorf("s.calcBaseFleetPower(): %w", err)
	}

	fleetMultiplier := min(float64(totalFleetPower)/float64(fleedPowerToGainMaxChanceBoostAmount), maxBoostFleetPowerMultiplier)

	totalFleetCapacity, err := s.calcFleetCapacity(fleet)
	if err != nil {
		return models.Boost{}, fmt.Errorf("s.calcFleetCapacity(): %w", err)
	}

	capacityMultiplier := min(float64(totalFleetCapacity)/float64(fleedCapacityToGainMaxChanceBoostAmount), maxBoostFleetCapacityMultiplier)

	totalMultiplier := fleetMultiplier + capacityMultiplier

	if totalMultiplier < 0.33 {
		return models.Boost{
			Count: getBoostCount(totalMultiplier),
			ID:    consts.BoostID1,
		}, nil
	}

	if totalMultiplier < 0.66 {
		return models.Boost{
			Count: getBoostCount(totalMultiplier - 0.33),
			ID:    consts.BoostID2,
		}, nil
	}

	return models.Boost{
		Count: getBoostCount(totalMultiplier - 0.66),
		ID:    consts.BoostID3,
	}, nil
}

// getBoostCount calculates the count of target boost tier
func getBoostCount(boostCountMultiplier float64) uint64 {
	// each tier has 0.33 percent. There are only 3 tiers and 3 count rewards
	const maxChanceMultiplier = 0.33 / 3
	if boostCountMultiplier < maxChanceMultiplier {
		return 1
	}

	if boostCountMultiplier < maxChanceMultiplier*2 {
		return 2
	}

	return 3
}
