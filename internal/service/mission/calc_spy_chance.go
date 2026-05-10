package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

const (
	defaultSpyResourcesChance  = 1.5
	defaultSpyBuildingsChance  = 1.0
	defaultSpyFleetChance      = 0.5
	defaultSpyResearchesChance = 0.3
)

type spyChancesResult struct {
	spyResources  bool
	spyBuildings  bool
	spyFleet      bool
	spyResearches bool
}

func (s *Service) calcSpyChance(ctx context.Context, userID uuid.UUID, spyShipsCount uint64, storage TxStorages) (spyChancesResult, error) {
	spyResearchStats, err := s.getResearchByType(ctx, userID, consts.ResearchTypeSpyTechnology, storage)
	if err != nil {
		return spyChancesResult{}, fmt.Errorf("s.getResearchByType(): %w", err)
	}

	res := spyChancesResult{
		spyResources:  s.getSpyChance(float32(spyShipsCount) * float32(spyResearchStats.Bonuses.SpyChanceImprove*defaultSpyResourcesChance)),
		spyBuildings:  s.getSpyChance(float32(spyShipsCount) * float32(spyResearchStats.Bonuses.SpyChanceImprove*defaultSpyBuildingsChance)),
		spyFleet:      s.getSpyChance(float32(spyShipsCount) * float32(spyResearchStats.Bonuses.SpyChanceImprove*defaultSpyFleetChance)),
		spyResearches: s.getSpyChance(float32(spyShipsCount) * float32(spyResearchStats.Bonuses.SpyChanceImprove*defaultSpyResearchesChance)),
	}

	return res, nil
}

func (s *Service) getSpyChance(persent float32) bool {
	const maxChance = 100

	return s.randGenerator.Intn(maxChance) < int(persent)
}
