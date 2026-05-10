package notifications

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type SpyV1 struct {
	IsSpy       bool                `json:"isSpy"`
	Login       string              `json:"login"`
	Coordinates Coordinates         `json:"coordinates"`
	Resources   Resources           `json:"resources,omitzero"`
	Buildings   []consts.BuildingID `json:"buildings,omitempty"`
	Fleet       []FleetUnit         `json:"fleet,omitempty"`
	Researches  []consts.ResearchID `json:"researches,omitempty"`
}
