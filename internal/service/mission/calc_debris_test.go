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

func TestService_calcDebris(t *testing.T) {
	ship1ID := consts.FleetUnitID(1)
	ship2ID := consts.FleetUnitID(2)

	tests := []struct {
		name              string
		fleetBeforeAttack []models.FleetUnit
		fleetAfterAttack  []models.FleetUnit
		setupRegistry     func(reg *mocks.RegistryProvider)
		want              models.Resources
	}{
		{
			name: "no destroyed units",
			fleetBeforeAttack: []models.FleetUnit{
				{ID: ship1ID, Count: 100},
				{ID: ship2ID, Count: 50},
			},
			fleetAfterAttack: []models.FleetUnit{
				{ID: ship1ID, Count: 100},
				{ID: ship2ID, Count: 50},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {},
			want: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     0,
			},
		},
		{
			name: "partial fleet destruction",
			fleetBeforeAttack: []models.FleetUnit{
				{ID: ship1ID, Count: 100},
				{ID: ship2ID, Count: 50},
			},
			fleetAfterAttack: []models.FleetUnit{
				{ID: ship1ID, Count: 80},
				{ID: ship2ID, Count: 40},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, MetalCost: 100, CrystalCost: 50, GasCost: 25}, nil,
				).Once()
				reg.EXPECT().GetFleetUnitStatsByID(ship2ID).Return(
					registry.FleetUnitStats{ID: ship2ID, MetalCost: 200, CrystalCost: 100, GasCost: 50}, nil,
				).Once()
			},
			want: models.Resources{
				Metal:   800,
				Crystal: 400,
				Gas:     0,
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

			got, err := svc.calcDebris(tt.fleetBeforeAttack, tt.fleetAfterAttack)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
