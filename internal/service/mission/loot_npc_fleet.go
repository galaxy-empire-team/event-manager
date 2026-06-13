package mission

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const (
	defaultLootFleetChance = 25
)

func (s *Service) lootNPCFleet(playerFleet []models.FleetUnit, npcPosition consts.PlanetPositionZ, lootingMultiplier float32) ([]models.FleetUnit, error) {
	if defaultLootFleetChance < s.randGenerator.Intn(maxChance) {
		return playerFleet, nil
	}

	npcStats, err := s.registry.GetNPCStatsByPosition(npcPosition)
	if err != nil {
		return nil, fmt.Errorf("s.registry.GetNPCStatsByPosition(): %w", err)
	}

	playerFleetMap := lo.SliceToMap(playerFleet, func(fleetUnit models.FleetUnit) (consts.FleetUnitID, uint64) {
		return fleetUnit.ID, fleetUnit.Count
	})

	for _, fleetUnit := range npcStats.LootFleet {
		playerFleetMap[fleetUnit.ID] += uint64(lootingMultiplier * float32(fleetUnit.Count))
	}

	return lo.MapToSlice(playerFleetMap, func(id consts.FleetUnitID, count uint64) models.FleetUnit {
		return models.FleetUnit{
			ID:    id,
			Count: count,
		}
	}), nil
}
