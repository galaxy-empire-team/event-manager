package mission

import (
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const (
	// Multiplier depends on both power and capacity
	fleedPowerToGainMaxChanceMatterAmount    = 200_000    // equals 30 anihilators
	fleedCapacityToGainMaxChanceMatterAmount = 10_000_000 // equals 1000 transporters
	maxMatterFleetPowerMultiplier            = 0.5
	maxMatterFleetCapacityMultiplier         = 0.5
	maxMatterGainAmount                      = 30
	minMatterRewardCount                     = 1
)

func (s *Service) calcMatterMistReward(fleet []models.FleetUnit) (uint64, error) {
	totalFleetPower, err := s.calcBaseFleetPower(fleet)
	if err != nil {
		return 0, fmt.Errorf("s.calcBaseFleetPower(): %w", err)
	}

	fleetMultiplier := min(float64(totalFleetPower)/float64(fleedPowerToGainMaxChanceMatterAmount), maxMatterFleetPowerMultiplier)

	totalFleetCapacity, err := s.calcFleetCapacity(fleet)
	if err != nil {
		return 0, fmt.Errorf("s.calcFleetCapacity(): %w", err)
	}

	capacityMultiplier := min(float64(totalFleetCapacity)/float64(fleedCapacityToGainMaxChanceMatterAmount), maxMatterFleetCapacityMultiplier)

	totalMultiplier := fleetMultiplier + capacityMultiplier

	return max(uint64(float64(totalMultiplier)*float64(maxMatterGainAmount)), uint64(minMatterRewardCount)), nil
}
