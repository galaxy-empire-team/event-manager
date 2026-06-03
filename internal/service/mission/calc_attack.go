package mission

import (
	"fmt"

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
	attackerFleet      []models.FleetUnit
	attackerResearches userResearchBonuses
	defenderFleet      []models.FleetUnit
	defenderResearches userResearchBonuses
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

func (s *Service) calcAttackResult(input attackSetup) (attackResult, error) {
	attackerPower, err := s.calcFleetPower(input.attackerFleet, input.attackerResearches)
	if err != nil {
		return attackResult{}, fmt.Errorf("calcFleetPower(attacker): %w", err)
	}

	defenderPower, err := s.calcFleetPower(input.defenderFleet, input.defenderResearches)
	if err != nil {
		return attackResult{}, fmt.Errorf("calcFleetPower(defender): %w", err)
	}

	// Determines damage based on attack vs defense ratio.
	// Returns proportionally scaled damage
	// If defense greatly exceeds attack, damage approaches zero.
	damageToAttacker := defenderPower.attackPower * defenderPower.attackPower / (defenderPower.attackPower + attackerPower.defensePower)
	damageToDefender := attackerPower.attackPower * attackerPower.attackPower / (attackerPower.attackPower + defenderPower.defensePower)

	// Apply damage to fleets
	updatedAttackerFleet, err := s.calcFleetLoss(input.attackerFleet, damageToAttacker, input.attackerResearches)
	if err != nil {
		return attackResult{}, fmt.Errorf("calcFleetLoss(attacker): %w", err)
	}

	updatedDefenderFleet, err := s.calcFleetLoss(input.defenderFleet, damageToDefender, input.defenderResearches)
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
