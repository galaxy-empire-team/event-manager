package building

import (
	"time"

	"github.com/google/uuid"
)

type Planet struct {
	ID          uuid.UUID
	X           uint8
	Y           uint8
	Z           uint8
	Resources   Resources
	HasMoon     bool
	ColonizedAt time.Time
}

type Resources struct {
	Metal     uint64
	Crystal   uint64
	Gas       uint64
	UpdatedAt time.Time
}

type PlanetToColonize struct {
	ID uuid.UUID
	X  uint8
	Y  uint8
	Z  uint8
}
