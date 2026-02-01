package models

import (
	"time"

	"github.com/google/uuid"
)

type BuildType string

type BuildEvent struct {
	ID         uint64
	PlanetID   uuid.UUID
	BuildType  BuildType
	StartedAt  time.Time
	FinishedAt time.Time
}

type PlanetBuilding struct {
	ID        uint64
	PlanetID  uuid.UUID
	BuildType BuildType
	Level     uint8
}
