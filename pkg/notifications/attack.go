package notifications

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type AttackV1 struct {
	AttackerWins bool       `json:"attackerWins"`
	Cargo        Resources  `json:"cargo,omitempty"`
	Attacker     AttackInfo `json:"attacker"`
	Defender     AttackInfo `json:"defender"`
}

type AttackInfo struct {
	Login  string            `json:"login"`
	Planet Coordinates       `json:"planet"`
	Fleet  []AttackFleetUnit `json:"fleet"`
}

type AttackFleetUnit struct {
	ID          consts.FleetUnitID `json:"id"`
	CountBefore uint64             `json:"countBefore"`
	CountAfter  uint64             `json:"countAfter"`
}
