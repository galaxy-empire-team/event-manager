package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

type researchStorage interface {
	GetUserResearchesByTypes(ctx context.Context, userID uuid.UUID, researchTypes []consts.ResearchType) (map[consts.ResearchType]consts.ResearchID, error)
}

type registryProvider interface {
	GetBuildingStatsByID(buildingID consts.BuildingID) (registry.BuildingStats, error)
	GetBuildingZeroLvlIDByType(buildingType consts.BuildingType) (consts.BuildingID, error)
	GetResearchZeroLvlIDByType(researchType consts.ResearchType) (consts.ResearchID, error)
	GetResearchStatsByID(researchID consts.ResearchID) (registry.ResearchStats, error)
}

type Repository struct {
	researchStorage researchStorage
	registry        registryProvider
}

func New(researchStorage researchStorage, registry registryProvider) *Repository {
	return &Repository{
		researchStorage: researchStorage,
		registry:        registry,
	}
}
