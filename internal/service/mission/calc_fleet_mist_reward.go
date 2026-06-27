package mission

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const (
	powerFleetMultiplier    = 0.025  // reward is 5% of total fleet power
	maxFleetRewardID        = 4      // cruiser
	maxFleetPowerGainAmount = 16_200 // equals 100 cruisers
	minFleetRewardCount     = 1
)

func (s *Service) calcFleetMistReward(fleet []models.FleetUnit) (models.FleetUnit, error) {
	totalFleetPower, err := s.calcBaseFleetPower(fleet)
	if err != nil {
		return models.FleetUnit{}, fmt.Errorf("s.calcBaseFleetPower(): %w", err)
	}

	fleetID := consts.FleetUnitID(s.randGenerator.Intn(maxFleetRewardID) + 1)
	unitStats, err := s.registry.GetFleetUnitStatsByID(fleetID)
	if err != nil {
		return models.FleetUnit{}, fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
	}

	expectFleetPower := min(uint64(float64(totalFleetPower)*powerFleetMultiplier), maxFleetPowerGainAmount)

	rewardFleetPower := unitStats.Attack + unitStats.Defense
	rewardFleetCount := max(expectFleetPower/rewardFleetPower, minFleetRewardCount)

	// I don't like remainder and remove it
	if rewardFleetCount > 1000 {
		rewardFleetCount = (1000 / 10) * rewardFleetCount
	}

	return models.FleetUnit{
		ID:    fleetID,
		Count: rewardFleetCount,
	}, nil
}

func (s *Service) calcBaseFleetPower(fleet []models.FleetUnit) (uint64, error) {
	var totalFleetPower uint64
	for _, unit := range fleet {
		unitStats, err := s.registry.GetFleetUnitStatsByID(unit.ID)
		if err != nil {
			return 0, fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
		}

		totalFleetPower += unitStats.Attack * unit.Count
		totalFleetPower += unitStats.Defense * unit.Count
	}

	return totalFleetPower, nil
}
