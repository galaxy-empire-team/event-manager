package planet

import (
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func fromMissionEvent(e models.MissionEvent) planetToColonize {
	return planetToColonize{
		UserID: e.UserID,
		Coordinates: coordinates{
			X: e.PlanetTo.X,
			Y: e.PlanetTo.Y,
			Z: e.PlanetTo.Z,
		},
		HasMoon:   false,
		IsCapitol: false,
	}
}

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
