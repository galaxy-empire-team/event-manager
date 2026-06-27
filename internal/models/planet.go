package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// Insert userLogin to avoid additional requests to the database;
// TODO: Remove after notification pipeline impl.
type Planet struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	UserLogin   string
	Resources   Resources
	Coordinates Coordinates
}

type Coordinates struct {
	X consts.PlanetPositionX
	Y consts.PlanetPositionY
	Z consts.PlanetPositionZ
}

type Resources struct {
	Metal     uint64
	Crystal   uint64
	Gas       uint64
	Matter    uint64
	Boost     Boost
	UpdatedAt time.Time
}

func (r *Resources) IsEmptyBase() bool {
	return r.Metal == 0 && r.Crystal == 0 && r.Gas == 0
}

type Boost struct {
	ID    consts.BoostID
	Count uint64
}

type Fleet struct {
	Unit []FleetUnit
}

type FleetUnit struct {
	ID    consts.FleetUnitID
	Count uint64
}
