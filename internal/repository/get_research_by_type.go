package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

// GetResearchByType is a wrapper on GetResearchesByTypes to handle only one research type.
// Method returns stats for the target research type.
// If tech is not found in the database it returns zero-lvl stats.
func (r *Repository) GetResearchByType(ctx context.Context, userID uuid.UUID, researchType consts.ResearchType) (registry.ResearchStats, error) {
	researches, err := r.GetResearchesByTypes(ctx, userID, []consts.ResearchType{researchType})
	if err != nil {
		return registry.ResearchStats{}, fmt.Errorf("r.GetResearchesByTypes(): %w", err)
	}

	researchStats, ok := researches[researchType]
	if !ok {
		return registry.ResearchStats{}, fmt.Errorf("research stats for type %s not found", researchType)
	}

	return researchStats, nil
}

// GetResearchesByTypes returns stats for research types.
// If tech is not found in the database it returns zero-lvl stats.
func (r *Repository) GetResearchesByTypes(ctx context.Context, userID uuid.UUID, researchTypes []consts.ResearchType) (map[consts.ResearchType]registry.ResearchStats, error) {
	researchTypeToID, err := r.researchStorage.GetUserResearchesByTypes(ctx, userID, researchTypes)
	if err != nil {
		return nil, fmt.Errorf("researchStorage.GetUserResearchesByTypes(): %w", err)
	}

	res := make(map[consts.ResearchType]registry.ResearchStats)
	for _, researchType := range researchTypes {
		researchID, ok := researchTypeToID[researchType]
		if !ok {
			researchID, err = r.registry.GetResearchZeroLvlIDByType(researchType)
			if err != nil {
				return nil, fmt.Errorf("registry.GetResearchZeroLvlIDByType(%s): %w", researchType, err)
			}
		}

		researchStats, err := r.registry.GetResearchStatsByID(researchID)
		if err != nil {
			return nil, fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
		}

		res[researchType] = researchStats
	}

	return res, nil
}
