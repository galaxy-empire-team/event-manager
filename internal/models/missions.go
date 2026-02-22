package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type MissionEvent struct {
	ID          uint64
	MissionID   consts.MissionID
	UserID      uuid.UUID
	PlanetFrom  uuid.UUID
	PlanetTo    Coordinates
	Fleet       []FleetUnit
	IsReturning bool
	StartedAt   time.Time
	FinishedAt  time.Time
}
