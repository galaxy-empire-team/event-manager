package mission

import (
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const (
	cargoMetalTargetPercent   = 45
	cargoCrystalTargetPercent = 35
	cargoGasTargetPercent     = 20
)

type filledCargo struct {
	leftOnPlanet models.Resources
	gained       models.Resources
}

type preparedResult struct {
	leftOnPlanet uint64
	gained       uint64
}

func (s *Service) fillFleetCargo(planetResources models.Resources, cargo uint64) filledCargo {
	const (
		metalIdx   = 0
		crystalIdx = 1
		gasIdx     = 2
	)

	preparedRes := []preparedResult{
		initialFill(planetResources.Metal, cargo, cargoMetalTargetPercent),
		initialFill(planetResources.Crystal, cargo, cargoCrystalTargetPercent),
		initialFill(planetResources.Gas, cargo, cargoGasTargetPercent),
	}

	availableCargo := cargo - preparedRes[metalIdx].gained - preparedRes[crystalIdx].gained - preparedRes[gasIdx].gained

	fillRemainResources := func(idx int) {
		if preparedRes[idx].leftOnPlanet > 0 && availableCargo > 0 {
			fill := min(preparedRes[idx].leftOnPlanet, availableCargo)
			preparedRes[idx].gained += fill
			preparedRes[idx].leftOnPlanet -= fill
			availableCargo -= fill
		}
	}

	fillRemainResources(0)
	fillRemainResources(1)
	fillRemainResources(2)

	return filledCargo{
		gained: models.Resources{
			Metal:   preparedRes[metalIdx].gained,
			Crystal: preparedRes[crystalIdx].gained,
			Gas:     preparedRes[gasIdx].gained,
		},
		leftOnPlanet: models.Resources{
			Metal:     preparedRes[metalIdx].leftOnPlanet,
			Crystal:   preparedRes[crystalIdx].leftOnPlanet,
			Gas:       preparedRes[gasIdx].leftOnPlanet,
			UpdatedAt: planetResources.UpdatedAt,
		},
	}
}

func initialFill(resource uint64, cargo uint64, targetPercent uint64) preparedResult {
	target := cargo * targetPercent / 100
	gained := min(resource, target)

	return preparedResult{
		leftOnPlanet: resource - gained,
		gained:       gained,
	}
}
