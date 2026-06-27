package events

import (
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func toFleetUnits(fleet []models.FleetUnit) []fleetUnit {
	units := make([]fleetUnit, 0, len(fleet))

	for _, f := range fleet {
		units = append(units, fleetUnit{
			ID:    f.ID,
			Count: f.Count,
		})
	}

	return units
}

func toResources(r models.Resources) resources {
	return resources{
		Metal:   r.Metal,
		Crystal: r.Crystal,
		Gas:     r.Gas,
		Matter:  r.Matter,
		Boost: boost{
			ID:    r.Boost.ID,
			Count: r.Boost.Count,
		},
	}
}
