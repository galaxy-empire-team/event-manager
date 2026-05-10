package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

// getResearchByType is a wrapper on getResearchesByTypes to handle only one research type.
func (s *Service) getResearchByType(ctx context.Context, userID uuid.UUID, researchType consts.ResearchType, storage TxStorages) (registry.ResearchStats, error) {
	researches, err := s.getResearchesByTypes(ctx, userID, []consts.ResearchType{researchType}, storage)
	if err != nil {
		return registry.ResearchStats{}, fmt.Errorf("s.getResearchesByTypes(): %w", err)
	}

	researchStats, ok := researches[researchType]
	if !ok {
		return registry.ResearchStats{}, fmt.Errorf("research stats for type %s not found", researchType)
	}

	return researchStats, nil
}

// getResearchByType returns info about user technologies. If tech is not found in the database it returns zero-lvl stats.
func (s *Service) getResearchesByTypes(ctx context.Context, userID uuid.UUID, researchTypes []consts.ResearchType, storage TxStorages) (map[consts.ResearchType]registry.ResearchStats, error) {
	researchIDs, err := storage.GetUserResearchesByTypes(ctx, userID, researchTypes)
	if err != nil {
		return nil, fmt.Errorf("storage.GetUserResearchesByTypes(): %w", err)
	}

	res := make(map[consts.ResearchType]registry.ResearchStats)
	for _, researchType := range researchTypes {
		researchID, ok := researchIDs[researchType]
		if !ok {
			researchID, err = s.registry.GetResearchZeroLvlIDByType(researchType)
			if err != nil {
				return nil, fmt.Errorf("registry.GetResearchZeroLvlIDByType(%s): %w", researchType, err)
			}
		}

		researchStats, err := s.registry.GetResearchStatsByID(researchID)
		if err != nil {
			return nil, fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
		}

		res[researchType] = researchStats
	}

	return res, nil
}
