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

type BuildingUpgrade struct {
	PlanetID          uuid.UUID
	CurrentBuildingID consts.BuildingID
	UpdatedBuildingID consts.BuildingID
}
