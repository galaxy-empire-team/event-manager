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

func TestService_calcBoostMistReward(t *testing.T) {
	ship1ID := consts.FleetUnitID(1)
	ship2ID := consts.FleetUnitID(2)

	tests := []struct {
		name          string
		fleet         []models.FleetUnit
		setupRegistry func(reg *mocks.RegistryProvider)
		want          models.Boost
	}{
		{
			name:  "empty fleet returns tier 1 with count 1",
			fleet: []models.FleetUnit{},
			setupRegistry: func(reg *mocks.RegistryProvider) {
			},
			want: models.Boost{
				Count: 1,
				ID:    consts.BoostID1,
			},
		},
		{
			name: "small fleet returns tier 1 with count 1",
			fleet: []models.FleetUnit{
				{ID: ship1ID, Count: 5},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 10, Defense: 10, CargoCapacity: 100}, nil,
				).Times(2)
			},
			want: models.Boost{
				Count: 1,
				ID:    consts.BoostID1,
			},
		},
		{
			name: "medium fleet returns tier 1 with count 2",
			fleet: []models.FleetUnit{
				{ID: ship1ID, Count: 10},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 1000, Defense: 1000, CargoCapacity: 50_000}, nil,
				).Times(2)
			},
			want: models.Boost{
				Count: 2,
				ID:    consts.BoostID1,
			},
		},
		{
			name: "medium fleet returns tier 1 with count 3",
			fleet: []models.FleetUnit{
				{ID: ship1ID, Count: 17},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 1000, Defense: 1000, CargoCapacity: 50_000}, nil,
				).Times(2)
			},
			want: models.Boost{
				Count: 3,
				ID:    consts.BoostID1,
			},
		},
		{
			name: "fleet returns tier 2 with count 1",
			fleet: []models.FleetUnit{
				{ID: ship1ID, Count: 40},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 500, Defense: 500, CargoCapacity: 50_000}, nil,
				).Times(2)
			},
			want: models.Boost{
				Count: 1,
				ID:    consts.BoostID2,
			},
		},
		{
			name: "fleet at max power threshold returns tier 2 with count 2",
			fleet: []models.FleetUnit{
				{ID: ship1ID, Count: 30},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 3300, Defense: 3300, CargoCapacity: 0}, nil,
				).Times(2)
			},
			want: models.Boost{
				Count: 2,
				ID:    consts.BoostID2,
			},
		},
		{
			name: "fleet returns tier 2 with count 3",
			fleet: []models.FleetUnit{
				{ID: ship1ID, Count: 35},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 1000, Defense: 1000, CargoCapacity: 71_400}, nil,
				).Times(2)
			},
			want: models.Boost{
				Count: 3,
				ID:    consts.BoostID2,
			},
		},
		{
			name: "mixed fleet returns tier 3 with count 1",
			fleet: []models.FleetUnit{
				{ID: ship1ID, Count: 30},
				{ID: ship2ID, Count: 200},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 3300, Defense: 3300, CargoCapacity: 0}, nil,
				).Times(2)
				reg.EXPECT().GetFleetUnitStatsByID(ship2ID).Return(
					registry.FleetUnitStats{ID: ship2ID, Attack: 1, Defense: 1, CargoCapacity: 10_000}, nil,
				).Times(2)
			},
			want: models.Boost{
				Count: 1,
				ID:    consts.BoostID3,
			},
		},
		{
			name: "mixed fleet returns tier 3 with count 2",
			fleet: []models.FleetUnit{
				{ID: ship1ID, Count: 30},
				{ID: ship2ID, Count: 300},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 3300, Defense: 3300, CargoCapacity: 0}, nil,
				).Times(2)
				reg.EXPECT().GetFleetUnitStatsByID(ship2ID).Return(
					registry.FleetUnitStats{ID: ship2ID, Attack: 1, Defense: 1, CargoCapacity: 10_000}, nil,
				).Times(2)
			},
			want: models.Boost{
				Count: 2,
				ID:    consts.BoostID3,
			},
		},
		{
			name: "mixed fleet at power and capacity thresholds returns tier 3 with count 3",
			fleet: []models.FleetUnit{
				{ID: ship1ID, Count: 30},
				{ID: ship2ID, Count: 1000},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 3300, Defense: 3300, CargoCapacity: 0}, nil,
				).Times(2)
				reg.EXPECT().GetFleetUnitStatsByID(ship2ID).Return(
					registry.FleetUnitStats{ID: ship2ID, Attack: 1, Defense: 1, CargoCapacity: 10_000}, nil,
				).Times(2)
			},
			want: models.Boost{
				Count: 3,
				ID:    consts.BoostID3,
			},
		},
		{
			name: "large fleet above both thresholds returns tier 3 with count 3",
			fleet: []models.FleetUnit{
				{ID: ship1ID, Count: 60},
				{ID: ship2ID, Count: 2000},
			},
			setupRegistry: func(reg *mocks.RegistryProvider) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 3300, Defense: 3300, CargoCapacity: 0}, nil,
				).Times(2)
				reg.EXPECT().GetFleetUnitStatsByID(ship2ID).Return(
					registry.FleetUnitStats{ID: ship2ID, Attack: 1, Defense: 1, CargoCapacity: 10_000}, nil,
				).Times(2)
			},
			want: models.Boost{
				Count: 3,
				ID:    consts.BoostID3,
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

			got, err := svc.calcBoostMistReward(tt.fleet)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
