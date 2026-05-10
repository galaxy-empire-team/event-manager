package notifications

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type FleetUnit struct {
	ID    consts.FleetUnitID `json:"id"`
	Count uint64             `json:"count"`
}
