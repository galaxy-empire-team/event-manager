package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type BuildEvent struct {
	ID         uint64
	PlanetID   uuid.UUID
	BuildingID consts.BuildingID
	StartedAt  time.Time
	FinishedAt time.Time
}

type ResearchEvent struct {
	ID         uint64
	UserID     uuid.UUID
	ResearchID consts.ResearchID
	StartedAt  time.Time
	FinishedAt time.Time
}

type MissionEvent struct {
	ID          uint64
	MissionID   consts.MissionID
	UserID      uuid.UUID
	PlanetFrom  uuid.UUID
	PlanetTo    Coordinates
	Fleet       []FleetUnit
	Resources   Resources
	IsReturning bool
	StartedAt   time.Time
	FinishedAt  time.Time
}

type FleetConstructionEvent struct {
	ID         uint64
	PlanetID   uuid.UUID
	FleetID    consts.FleetUnitID
	Count      uint64
	StartedAt  time.Time
	FinishedAt time.Time
}

type NotificationEvent struct {
	UserID         uuid.UUID
	NotificationID consts.NotificationID
	Data           []byte
}
