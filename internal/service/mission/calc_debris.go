package mission

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const (
	debrisMultiplier = 0.2
)

func (s *Service) calcDebris(fleetBeforeAttack []models.FleetUnit, fleetAfterAttack []models.FleetUnit) (models.Resources, error) {
	fleetAfterAttackMap := lo.SliceToMap(fleetAfterAttack, func(fleetUnit models.FleetUnit) (consts.FleetUnitID, uint64) {
		return fleetUnit.ID, fleetUnit.Count
	})

	var res models.Resources
	for _, fleetUnit := range fleetBeforeAttack {
		destroyCount := fleetUnit.Count - fleetAfterAttackMap[fleetUnit.ID]

		if destroyCount == 0 {
			continue
		}

		stats, err := s.registry.GetFleetUnitStatsByID(fleetUnit.ID)
		if err != nil {
			return models.Resources{}, fmt.Errorf("s.registry.GetFleetUnitStatsByID(): %w", err)
		}

		// I don't want to spawn gas as debris
		res.Metal += uint64(float64(stats.MetalCost*destroyCount) * debrisMultiplier)
		res.Crystal += uint64(float64(stats.CrystalCost*destroyCount) * debrisMultiplier)
	}

	return res, nil
}
