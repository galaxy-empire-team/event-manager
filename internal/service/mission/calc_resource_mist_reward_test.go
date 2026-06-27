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

func TestService_calcResourceMistReward(t *testing.T) {
	shipID := consts.FleetUnitID(1)

	tests := []struct {
		name          string
		fleet         []models.FleetUnit
		setupRegistry func(reg *mocks.RegistryProvider)
		want          models.Resources
	}{
		{
			name:  "empty fleet returns zero resources",
			fleet: []models.FleetUnit{},
			setupRegistry: func(reg *mocks.RegistryProvider) {
			},
			want: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     0,
			},
		},
		{
			name: "small fleet returns proportional resources",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 1},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100, CargoCapacity: 100}, nil,
				).Once()
			},
			want: models.Resources{
				Metal:   8,
				Crystal: 8,
				Gas:     8,
			},
		},
		{
			name: "medium fleet returns proportional resources",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 10},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100, CargoCapacity: 12_000}, nil,
				).Once()
			},
			want: models.Resources{
				Metal:   10_000,
				Crystal: 10_000,
				Gas:     10_000,
			},
		},
		{
			name: "large fleet reward is capped at max gain amount",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 2400},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100, CargoCapacity: 10_000}, nil,
				).Once()
			},
			want: models.Resources{
				Metal:   2_000_000,
				Crystal: 2_000_000,
				Gas:     2_000_000,
			},
		},
		{
			name: "fleet just below max gain cap",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 2399},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100, CargoCapacity: 10_000}, nil,
				).Once()
			},
			want: models.Resources{
				Metal:   1_999_166,
				Crystal: 1_999_166,
				Gas:     1_999_166,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := mocks.NewRegistryProvider(t)
			tt.setupRegistry(reg)

			svc := &Service{
				registry: reg,
				logger:   zap.NewNop(),
			}

			got, err := svc.calcResourceMistReward(tt.fleet)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
