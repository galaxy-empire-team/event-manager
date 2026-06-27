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

func TestService_calcMatterMistReward(t *testing.T) {
	shipID := consts.FleetUnitID(1)

	tests := []struct {
		name          string
		fleet         []models.FleetUnit
		setupRegistry func(reg *mocks.RegistryProvider)
		want          uint64
	}{
		{
			name:  "empty fleet returns minimum matter reward",
			fleet: []models.FleetUnit{},
			setupRegistry: func(reg *mocks.RegistryProvider) {
			},
			want: 1,
		},
		{
			name: "small fleet returns minimum matter reward",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 1},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100, CargoCapacity: 100}, nil,
				).Times(2)
			},
			want: 1,
		},
		{
			name: "medium fleet returns proportional matter reward",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 375},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100, CargoCapacity: 0}, nil,
				).Times(2)
			},
			want: 11,
		},
		{
			name: "fleet at max power threshold",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 1000},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100, CargoCapacity: 0}, nil,
				).Times(2)
			},
			want: 15,
		},
		{
			name: "fleet at power and capacity thresholds",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 1000},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100, CargoCapacity: 10_000}, nil,
				).Times(2)
			},
			want: 30,
		},
		{
			name: "large fleet above both thresholds",
			fleet: []models.FleetUnit{
				{ID: shipID, Count: 2000},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(shipID).Return(
					registry.FleetUnitStats{ID: shipID, Attack: 100, Defense: 100, CargoCapacity: 10_000}, nil,
				).Times(2)
			},
			want: 30,
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

			got, err := svc.calcMatterMistReward(tt.fleet)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
