package models

import (
	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// Insert userLogin to avoid additional requests to the database;
// TODO: Remove after notification pipeline impl
type Planet struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	UserLogin   string
	Coordinates Coordinates
}

type Coordinates struct {
	X consts.PlanetPositionX
	Y consts.PlanetPositionY
	Z consts.PlanetPositionZ
}

type Resources struct {
	Metal   uint64
	Crystal uint64
	Gas     uint64
}

type Fleet struct {
	Unit []FleetUnit
}

type FleetUnit struct {
	ID    consts.FleetUnitID
	Count uint64
}
