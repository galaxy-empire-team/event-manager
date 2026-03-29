package models

import (
	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type BuildingUpgrade struct {
	PlanetID          uuid.UUID
	CurrentBuildingID consts.BuildingID
	UpdatedBuildingID consts.BuildingID
}

type BuildingInfo struct {
	ID          consts.BuildingID
	Type        consts.BuildingType
	ProductionS uint64
}
