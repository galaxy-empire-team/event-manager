package events

import (
	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type planetToColonize struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Coordinates coordinates
	HasMoon     bool
	IsCapitol   bool
}

type coordinates struct {
	X consts.PlanetPositionX
	Y consts.PlanetPositionY
	Z consts.PlanetPositionZ
}

type fleet struct {
	Fleet []fleetUnit `json:"fleet"`
}

type fleetUnit struct {
	ID    consts.FleetUnitID `json:"id"`
	Count uint64             `json:"count"`
}
