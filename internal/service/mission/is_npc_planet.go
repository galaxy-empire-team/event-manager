package mission

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) isPlanetNPC(positionZ consts.PlanetPositionZ) bool {
	if positionZ == consts.NPCTierOnePositionZ ||
		positionZ == consts.NPCTierTwoPositionZ ||
		positionZ == consts.NPCTierThreePositionZ {
		return true
	}

	return false
}
