package models

import (
	"time"

	"github.com/google/uuid"
)

type MissionType string

const (
	MissionTypeColonize MissionType = "colonize"
)

type MissionEvent struct {
	ID         uint64
	UserID     uuid.UUID
	PlanetFrom uuid.UUID
	PlanetTo   Coordinates
	Type       MissionType
	StartedAt  time.Time
	FinishedAt time.Time
}

type Coordinates struct {
	X uint8
	Y uint8
	Z uint8
}
