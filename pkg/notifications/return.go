package notifications

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type ReturnV1 struct {
	MissionType consts.MissionID `json:"mission_id"`
	Status      string           `json:"status"`
}
