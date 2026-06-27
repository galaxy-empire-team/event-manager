package notifications

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type RewardType string

const (
	RewardTypeResource RewardType = "resource"
	RewardTypeFleet    RewardType = "fleet"
	RewardTypeBoost    RewardType = "boost"
	RewardTypeMatter   RewardType = "matter"
)

type MistV1 struct {
	RewardType RewardType `json:"reward_type"`
	Reward     Reward     `json:"reward"`
}

type Reward struct {
	Resource MistResourceReward `json:"resource,omitzero"`
	Fleet    MistFleetReward    `json:"fleet,omitzero"`
	Boost    MistBoostReward    `json:"boost,omitzero"`
	Matter   MistMatterReward   `json:"matter,omitzero"`
}

type MistResourceReward struct {
	Metal   uint64 `json:"metal,omitzero"`
	Crystal uint64 `json:"crystal,omitzero"`
	Gas     uint64 `json:"gas,omitzero"`
}

type MistFleetReward struct {
	ID    consts.FleetUnitID `json:"id,omitzero"`
	Count uint64             `json:"count,omitzero"`
}

type MistBoostReward struct {
	Count uint64         `json:"count,omitzero"`
	ID    consts.BoostID `json:"id,omitzero"`
}

type MistMatterReward struct {
	Count uint64 `json:"count,omitzero"`
}
