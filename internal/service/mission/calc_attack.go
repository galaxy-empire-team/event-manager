package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const (
	// Minimum percentage of fleet that can be lost in an attack
	fleetLossMinMultiplier = 0.3
	// Multiplier to determine if attacker wins based on damage
	winTresholdMultiplier = 1.2
)

type attackSetup struct {
	attackerID    uuid.UUID
	defenderID    uuid.UUID
	attackerFleet []models.FleetUnit
	defenderFleet []models.FleetUnit
}

type attackResult struct {
	attackerWins      bool
	attackerFleetLeft []models.FleetUnit
	defenderFleetLeft []models.FleetUnit
}

type userResearchBonuses struct {
	attackBonus  float64
	defenseBonus float64
}

type totalPower struct {
	attackPower  uint64
	defensePower uint64
}

func (s *Service) calcAttackResult(ctx context.Context, input attackSetup, storage TxStorages) (attackResult, error) {
	attackerResearchBonuses, err := s.getResearchBonuses(ctx, input.attackerID, storage)
	if err != nil {
		return attackResult{}, fmt.Errorf("getResearchBonuses(attacker): %w", err)
	}

	attackerPower, err := s.calcFleetPower(input.attackerFleet, attackerResearchBonuses)
	if err != nil {
		return attackResult{}, fmt.Errorf("calcFleetPower(attacker): %w", err)
	}

	defenderResearchBonuses, err := s.getResearchBonuses(ctx, input.defenderID, storage)
	if err != nil {
		return attackResult{}, fmt.Errorf("getResearchBonuses(defender): %w", err)
	}

	defenderPower, err := s.calcFleetPower(input.defenderFleet, defenderResearchBonuses)
	if err != nil {
		return attackResult{}, fmt.Errorf("calcFleetPower(defender): %w", err)
	}

	// Determines damage based on attack vs defense ratio.
	// Returns proportionally scaled damage
	// If defense greatly exceeds attack, damage approaches zero.
	damageToAttacker := defenderPower.attackPower * defenderPower.attackPower / (defenderPower.attackPower + attackerPower.defensePower)
	damageToDefender := attackerPower.attackPower * attackerPower.attackPower / (attackerPower.attackPower + defenderPower.defensePower)

	// Apply damage to fleets
	updatedAttackerFleet, err := s.calcFleetLoss(input.attackerFleet, damageToAttacker, attackerResearchBonuses)
	if err != nil {
		return attackResult{}, fmt.Errorf("calcFleetLoss(attacker): %w", err)
	}

	updatedDefenderFleet, err := s.calcFleetLoss(input.defenderFleet, damageToDefender, defenderResearchBonuses)
	if err != nil {
		return attackResult{}, fmt.Errorf("calcFleetLoss(defender): %w", err)
	}

	var attackerWins bool
	if float64(damageToDefender) > float64(damageToAttacker)*winTresholdMultiplier {
		attackerWins = true
	}

	return attackResult{
		attackerWins:      attackerWins,
		attackerFleetLeft: updatedAttackerFleet,
		defenderFleetLeft: updatedDefenderFleet,
	}, nil
}

func (s *Service) calcFleetLoss(fleet []models.FleetUnit, damage uint64, researches userResearchBonuses) ([]models.FleetUnit, error) {
	type precomputedFleetLoss struct {
		unitID       consts.FleetUnitID
		shipWeight   uint64
		shipDefense  uint64
		damageToShip uint64
	}

	var totalFleetWeight uint64
	precomputedResults := make([]precomputedFleetLoss, 0, len(fleet))
	for _, unit := range fleet {
		unitStats, err := s.registry.GetFleetUnitStatsByID(unit.ID)
		if err != nil {
			return nil, fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
		}

		shipWeight := unitStats.Attack + unitStats.Defense
		totalFleetWeight += shipWeight
		precomputedResults = append(precomputedResults, precomputedFleetLoss{
			unitID:       unit.ID,
			shipWeight:   shipWeight,
			shipDefense:  uint64(float64(unitStats.Defense*unit.Count) * researches.defenseBonus),
			damageToShip: 0,
		})
	}

	// Calculate damage to each ship based on its weight in the fleet
	for i := range precomputedResults {
		precomputedResults[i].damageToShip = damage * precomputedResults[i].shipWeight / totalFleetWeight
	}

	// Calculate remaining ships after damage
	updatedFleet := make([]models.FleetUnit, 0, len(fleet))
	for i, fleetUnit := range fleet {
		shipLossMultiplier := 1 - float32(precomputedResults[i].damageToShip)/float32(precomputedResults[i].shipDefense)
		shipLossMultiplier = max(shipLossMultiplier, fleetLossMinMultiplier)
		updatedFleet = append(updatedFleet, models.FleetUnit{
			ID:    fleetUnit.ID,
			Count: uint64(float32(fleetUnit.Count) * shipLossMultiplier),
		})
	}

	return updatedFleet, nil
}

func (s *Service) calcFleetPower(fleet []models.FleetUnit, researches userResearchBonuses) (totalPower, error) {
	var fleetAttackPower, fleetDefensePower uint64
	for _, unit := range fleet {
		unitStats, err := s.registry.GetFleetUnitStatsByID(unit.ID)
		if err != nil {
			return totalPower{}, fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
		}

		fleetAttackPower += unitStats.Attack * unit.Count
		fleetDefensePower += unitStats.Defense * unit.Count
	}

	fleetAttackPower = uint64(float64(fleetAttackPower) * researches.attackBonus)
	fleetDefensePower = uint64(float64(fleetDefensePower) * researches.defenseBonus)

	return totalPower{
		attackPower:  fleetAttackPower,
		defensePower: fleetDefensePower,
	}, nil
}

func (s *Service) getResearchBonuses(ctx context.Context, userID uuid.UUID, storage TxStorages) (userResearchBonuses, error) {
	researchIDs, err := storage.GetUserResearchesByTypes(ctx, userID, []consts.ResearchType{consts.ResearchTypeWeaponTech, consts.ResearchTypeArmorTech})
	if err != nil {
		return userResearchBonuses{}, fmt.Errorf("storage.GetUserResearchesByTypes(): %w", err)
	}

	weaponTechID, ok := researchIDs[consts.ResearchTypeWeaponTech]
	if !ok {
		return userResearchBonuses{}, fmt.Errorf("weapon tech research not found for user %s", userID.String())
	}

	weaponTechStats, err := s.registry.GetResearchStatsByID(weaponTechID)
	if err != nil {
		return userResearchBonuses{}, fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
	}

	armorTechID, ok := researchIDs[consts.ResearchTypeArmorTech]
	if !ok {
		return userResearchBonuses{}, fmt.Errorf("armor tech research not found for user %s", userID.String())
	}

	armorTechStats, err := s.registry.GetResearchStatsByID(armorTechID)
	if err != nil {
		return userResearchBonuses{}, fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
	}

	return userResearchBonuses{
		attackBonus:  float64(weaponTechStats.Bonuses.AttackPower),
		defenseBonus: float64(armorTechStats.Bonuses.ArmorPower),
	}, nil
}
