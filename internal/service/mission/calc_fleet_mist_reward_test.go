package mission

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
	"github.com/galaxy-empire-team/event-manager/internal/models"
	"github.com/galaxy-empire-team/event-manager/internal/service/mission/mocks"
)

func TestService_calcFleetMistReward(t *testing.T) {
	shipID := consts.FleetUnitID(1)

	tests := []struct {
		name  string
		fleet []models.FleetUnit
		setup func(reg *mocks.RegistryProvider, randGen *mocks.RandGenerator)
		want  models.FleetUnit
	}{
		{
			name:  "empty fleet returns minimum reward count",
			fleet: []models.FleetUnit{},
			setup: func(reg *mocks.RegistryProvider, randGen *mocks.RandGenerator) {
				randGen.EXPECT().Intn(maxFleetRewardID).Return(0).Once()
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100}, nil,
				).Once()
			},
			want: models.FleetUnit{
				ID:    shipID,
				Count: 1,
			},
		},
		{
			name: "small fleet returns minimum reward count",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 1},
			},
			setup: func(reg *mocks.RegistryProvider, randGen *mocks.RandGenerator) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100}, nil,
				).Times(2)

				randGen.EXPECT().Intn(maxFleetRewardID).Return(0).Once()
			},
			want: models.FleetUnit{
				ID:    shipID,
				Count: 1,
			},
		},
		{
			name: "medium fleet returns proportional reward",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 500},
			},
			setup: func(reg *mocks.RegistryProvider, randGen *mocks.RandGenerator) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100}, nil,
				).Times(2)

				randGen.EXPECT().Intn(maxFleetRewardID).Return(0).Once()
			},
			want: models.FleetUnit{
				ID:    shipID,
				Count: 12,
			},
		},
		{
			name: "large fleet reward is capped at max power gain",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 3240},
			},
			setup: func(reg *mocks.RegistryProvider, randGen *mocks.RandGenerator) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100}, nil,
				).Times(2)

				randGen.EXPECT().Intn(maxFleetRewardID).Return(0).Once()
			},
			want: models.FleetUnit{
				ID:    shipID,
				Count: 81,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := mocks.NewRegistryProvider(t)
			randGen := mocks.NewRandGenerator(t)
			tt.setup(reg, randGen)

			svc := &Service{
				registry:      reg,
				randGenerator: randGen,
				logger:        zap.NewNop(),
			}

			got, err := svc.calcFleetMistReward(tt.fleet)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
