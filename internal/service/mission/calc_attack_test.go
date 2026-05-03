package mission

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
	"github.com/galaxy-empire-team/event-manager/internal/models"
	"github.com/galaxy-empire-team/event-manager/internal/service/mission/mocks"
)

func TestService_calcAttackResult(t *testing.T) {
	ctx := t.Context()

	attackerID := uuid.New()
	defenderID := uuid.New()

	// Ship stats as specified
	ship1ID := consts.FleetUnitID(1)
	ship2ID := consts.FleetUnitID(2)
	ship3ID := consts.FleetUnitID(3)
	ship4ID := consts.FleetUnitID(4)

	weaponResID := consts.ResearchID(10)
	armorResID := consts.ResearchID(20)

	researchMap := map[consts.ResearchType]consts.ResearchID{
		consts.ResearchTypeWeaponTech: weaponResID,
		consts.ResearchTypeArmorTech:  armorResID,
	}

	tests := []struct {
		name        string
		input       attackSetup
		expectSetup func(reg *mocks.RegistryProvider, storage *mocks.TxStorages)
		want        attackResult
	}{
		{
			name: "attacker wins with stronger fleet",
			input: attackSetup{
				attackerID:    attackerID,
				defenderID:    defenderID,
				attackerFleet: []models.FleetUnit{{ID: ship1ID, Count: 100}},
				defenderFleet: []models.FleetUnit{{ID: ship1ID, Count: 10}},
			},
			expectSetup: func(reg *mocks.RegistryProvider, storage *mocks.TxStorages) {
				shipStats := registry.FleetUnitStats{ID: ship1ID, Attack: 1, Defense: 1}
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(shipStats, nil).Maybe()
				reg.EXPECT().GetResearchStatsByID(weaponResID).Return(
					registry.ResearchStats{Bonuses: registry.ResearchBonuses{AttackPower: 0.5}}, nil,
				).Maybe()
				reg.EXPECT().GetResearchStatsByID(armorResID).Return(
					registry.ResearchStats{Bonuses: registry.ResearchBonuses{ArmorPower: 0.5}}, nil,
				).Maybe()

				storage.EXPECT().
					GetUserResearchesByTypes(ctx, attackerID, []consts.ResearchType{consts.ResearchTypeWeaponTech, consts.ResearchTypeArmorTech}).
					Return(researchMap, nil).Once()
				storage.EXPECT().
					GetUserResearchesByTypes(ctx, defenderID, []consts.ResearchType{consts.ResearchTypeWeaponTech, consts.ResearchTypeArmorTech}).
					Return(researchMap, nil).Once()
			},
			want: attackResult{
				attackerWins:      true,
				attackerFleetLeft: []models.FleetUnit{{ID: ship1ID, Count: 100}},
				defenderFleetLeft: []models.FleetUnit{{ID: ship1ID, Count: 3}},
			},
		},
		{
			name: "attacker wins with stronger mixed fleet",
			input: attackSetup{
				attackerID: attackerID,
				defenderID: defenderID,
				attackerFleet: []models.FleetUnit{
					{ID: ship1ID, Count: 100},
					{ID: ship2ID, Count: 50},
					{ID: ship3ID, Count: 30},
					{ID: ship4ID, Count: 20},
				},
				defenderFleet: []models.FleetUnit{
					{ID: ship1ID, Count: 30},
					{ID: ship2ID, Count: 20},
					{ID: ship3ID, Count: 10},
					{ID: ship4ID, Count: 6},
				},
			},
			expectSetup: func(reg *mocks.RegistryProvider, storage *mocks.TxStorages) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 1, Defense: 1}, nil,
				).Maybe()
				reg.EXPECT().GetFleetUnitStatsByID(ship2ID).Return(
					registry.FleetUnitStats{ID: ship2ID, Attack: 6, Defense: 5}, nil,
				).Maybe()
				reg.EXPECT().GetFleetUnitStatsByID(ship3ID).Return(
					registry.FleetUnitStats{ID: ship3ID, Attack: 15, Defense: 12}, nil,
				).Maybe()
				reg.EXPECT().GetFleetUnitStatsByID(ship4ID).Return(
					registry.FleetUnitStats{ID: ship4ID, Attack: 100, Defense: 80}, nil,
				).Maybe()
				reg.EXPECT().GetResearchStatsByID(weaponResID).Return(
					registry.ResearchStats{Bonuses: registry.ResearchBonuses{AttackPower: 1.0}}, nil,
				).Maybe()
				reg.EXPECT().GetResearchStatsByID(armorResID).Return(
					registry.ResearchStats{Bonuses: registry.ResearchBonuses{ArmorPower: 1.0}}, nil,
				).Maybe()

				storage.EXPECT().
					GetUserResearchesByTypes(ctx, attackerID, []consts.ResearchType{consts.ResearchTypeWeaponTech, consts.ResearchTypeArmorTech}).
					Return(researchMap, nil).Once()
				storage.EXPECT().
					GetUserResearchesByTypes(ctx, defenderID, []consts.ResearchType{consts.ResearchTypeWeaponTech, consts.ResearchTypeArmorTech}).
					Return(researchMap, nil).Once()
			},
			want: attackResult{
				attackerWins: true,
				attackerFleetLeft: []models.FleetUnit{
					{ID: ship1ID, Count: 98},
					{ID: ship2ID, Count: 47},
					{ID: ship3ID, Count: 27},
					{ID: ship4ID, Count: 17},
				},
				defenderFleetLeft: []models.FleetUnit{
					{ID: ship1ID, Count: 9},
					{ID: ship2ID, Count: 6},
					{ID: ship3ID, Count: 3},
					{ID: ship4ID, Count: 1},
				},
			},
		},
		{
			name: "defender wins with the same fleet",
			input: attackSetup{
				attackerID: attackerID,
				defenderID: defenderID,
				attackerFleet: []models.FleetUnit{
					{ID: ship1ID, Count: 100},
					{ID: ship2ID, Count: 50},
					{ID: ship3ID, Count: 30},
					{ID: ship4ID, Count: 20},
				},
				defenderFleet: []models.FleetUnit{
					{ID: ship1ID, Count: 100},
					{ID: ship2ID, Count: 50},
					{ID: ship3ID, Count: 30},
					{ID: ship4ID, Count: 20},
				},
			},
			expectSetup: func(reg *mocks.RegistryProvider, storage *mocks.TxStorages) {
				reg.EXPECT().GetFleetUnitStatsByID(ship1ID).Return(
					registry.FleetUnitStats{ID: ship1ID, Attack: 1, Defense: 1}, nil,
				).Maybe()
				reg.EXPECT().GetFleetUnitStatsByID(ship2ID).Return(
					registry.FleetUnitStats{ID: ship2ID, Attack: 6, Defense: 5}, nil,
				).Maybe()
				reg.EXPECT().GetFleetUnitStatsByID(ship3ID).Return(
					registry.FleetUnitStats{ID: ship3ID, Attack: 15, Defense: 12}, nil,
				).Maybe()
				reg.EXPECT().GetFleetUnitStatsByID(ship4ID).Return(
					registry.FleetUnitStats{ID: ship4ID, Attack: 100, Defense: 80}, nil,
				).Maybe()
				reg.EXPECT().GetResearchStatsByID(weaponResID).Return(
					registry.ResearchStats{Bonuses: registry.ResearchBonuses{AttackPower: 1.0}}, nil,
				).Maybe()
				reg.EXPECT().GetResearchStatsByID(armorResID).Return(
					registry.ResearchStats{Bonuses: registry.ResearchBonuses{ArmorPower: 1.0}}, nil,
				).Maybe()

				storage.EXPECT().
					GetUserResearchesByTypes(ctx, attackerID, []consts.ResearchType{consts.ResearchTypeWeaponTech, consts.ResearchTypeArmorTech}).
					Return(researchMap, nil).Once()
				storage.EXPECT().
					GetUserResearchesByTypes(ctx, defenderID, []consts.ResearchType{consts.ResearchTypeWeaponTech, consts.ResearchTypeArmorTech}).
					Return(researchMap, nil).Once()
			},
			want: attackResult{
				attackerWins: false,
				attackerFleetLeft: []models.FleetUnit{
					{ID: ship1ID, Count: 86},
					{ID: ship2ID, Count: 34},
					{ID: ship3ID, Count: 13},
					{ID: ship4ID, Count: 6},
				},
				defenderFleetLeft: []models.FleetUnit{
					{ID: ship1ID, Count: 86},
					{ID: ship2ID, Count: 34},
					{ID: ship3ID, Count: 13},
					{ID: ship4ID, Count: 6},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := mocks.NewRegistryProvider(t)
			storage := mocks.NewTxStorages(t)
			tt.expectSetup(reg, storage)

			svc := &Service{
				registry: reg,
				logger:   zap.NewNop(),
			}

			got, err := svc.calcAttackResult(ctx, tt.input, storage)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
