package mission

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
