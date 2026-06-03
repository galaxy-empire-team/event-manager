package mission

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
	"github.com/galaxy-empire-team/event-manager/internal/service/mission/mocks"
)

func TestService_calcSpyChance(t *testing.T) {
	ctx := t.Context()

	userID := uuid.New()
	spyResID := consts.ResearchID(30)

	tests := []struct {
		name        string
		spyShips    uint64
		expectSetup func(reg *mocks.RegistryProvider, storage *mocks.TxStorages, randGen *mocks.RandGenerator)
		want        spyChancesResult
	}{
		{
			name:     "chances with 1 spyship",
			spyShips: 1,
			expectSetup: func(reg *mocks.RegistryProvider, storage *mocks.TxStorages, randGen *mocks.RandGenerator) {
				researchMap := map[consts.ResearchType]consts.ResearchID{
					consts.ResearchTypeSpyTechnology: spyResID,
				}
				storage.EXPECT().
					GetUserResearchesByTypes(ctx, userID, []consts.ResearchType{consts.ResearchTypeSpyTechnology}).
					Return(researchMap, nil).Once()

				reg.EXPECT().GetResearchStatsByID(spyResID).Return(
					registry.ResearchStats{
						Bonuses: registry.ResearchBonuses{SpyChanceMuliplier: 4.0},
					}, nil,
				).Once()

				randGen.EXPECT().Intn(100).Return(5).Once()  // spyResources chance
				randGen.EXPECT().Intn(100).Return(60).Once() // spyBuildings chance
				randGen.EXPECT().Intn(100).Return(30).Once() // spyFleet chance
				randGen.EXPECT().Intn(100).Return(5).Once()  // spyResearches chance
			},
			want: spyChancesResult{
				spyResources:  true,
				spyBuildings:  false,
				spyFleet:      false,
				spyResearches: false,
			},
		},
		{
			name:     "chances with 20 spyships",
			spyShips: 20,
			expectSetup: func(reg *mocks.RegistryProvider, storage *mocks.TxStorages, randGen *mocks.RandGenerator) {
				researchMap := map[consts.ResearchType]consts.ResearchID{
					consts.ResearchTypeSpyTechnology: spyResID,
				}
				storage.EXPECT().
					GetUserResearchesByTypes(ctx, userID, []consts.ResearchType{consts.ResearchTypeSpyTechnology}).
					Return(researchMap, nil).Once()

				reg.EXPECT().GetResearchStatsByID(spyResID).Return(
					registry.ResearchStats{
						Bonuses: registry.ResearchBonuses{SpyChanceMuliplier: 4.0},
					}, nil,
				).Once()

				randGen.EXPECT().Intn(100).Return(100).Once() // spyResources chance
				randGen.EXPECT().Intn(100).Return(80).Once()  // spyBuildings chance
				randGen.EXPECT().Intn(100).Return(50).Once()  // spyFleet chance
				randGen.EXPECT().Intn(100).Return(5).Once()   // spyResearches chance
			},
			want: spyChancesResult{
				spyResources:  true,
				spyBuildings:  false,
				spyFleet:      false,
				spyResearches: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := mocks.NewRegistryProvider(t)
			storage := mocks.NewTxStorages(t)
			randGen := mocks.NewRandGenerator(t)
			tt.expectSetup(reg, storage, randGen)

			svc := &Service{
				registry:      reg,
				randGenerator: randGen,
				logger:        zap.NewNop(),
			}

			got, err := svc.calcSpyChance(ctx, userID, tt.spyShips, storage)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
