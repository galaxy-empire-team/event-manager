package mission

import (
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const (
	// Mist zone reward amounts
	capacityFillMultiplier = 0.25
	maxMetalGainAmount     = 2_000_000
	maxCrystalGainAmount   = 2_000_000
	maxGasGainAmount       = 2_000_000
)

func (s *Service) calcResourceMistReward(fleet []models.FleetUnit) (models.Resources, error) {
	totalCargoCapacity, err := s.calcFleetCapacity(fleet)
	if err != nil {
		return models.Resources{}, fmt.Errorf("s.calcFleetCapacity(): %w", err)
	}

	totalEachResource := (float64(totalCargoCapacity) * capacityFillMultiplier) / 3

	return models.Resources{
		Metal:   min(uint64(totalEachResource), maxMetalGainAmount),
		Crystal: min(uint64(totalEachResource), maxCrystalGainAmount),
		Gas:     min(uint64(totalEachResource), maxGasGainAmount),
	}, nil
}
