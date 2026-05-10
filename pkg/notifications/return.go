package notifications

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type ReturnV1 struct {
	MissionType consts.MissionType `json:"mission_type"`
	Status      string             `json:"status"`
}
