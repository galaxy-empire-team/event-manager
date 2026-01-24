package models

import (
	"time"

	"github.com/google/uuid"
)

type BuildType string

const (
	BuildingTypeMetalMine   BuildType = "metal_mine"
	BuildingTypeCrystalMine BuildType = "crystal_mine"
	BuildingTypeGasMine     BuildType = "gas_mine"
)

type BuildEvent struct {
	ID         uint64
	PlanetID   uuid.UUID
	BuildType  BuildType
	StartdAt   time.Time
	FinishedAt time.Time
}
