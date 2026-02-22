package planet

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
	Fleet []fleetUnit
}

type fleetUnit struct {
	ID    consts.FleetUnitID
	Count uint64
}

type planetFleetUnit struct {
	PlanetID uuid.UUID          `json:"planet_id"`
	FleetID  consts.FleetUnitID `json:"fleet_id"`
	Count    uint64             `json:"count"`
}
