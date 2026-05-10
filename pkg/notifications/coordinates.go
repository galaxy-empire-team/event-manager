package notifications

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type Coordinates struct {
	X consts.PlanetPositionX `json:"x"`
	Y consts.PlanetPositionY `json:"y"`
	Z consts.PlanetPositionZ `json:"z"`
}
