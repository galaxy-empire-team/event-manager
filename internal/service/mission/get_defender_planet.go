package mission

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

type attackDefender struct {
	fleet      []models.FleetUnit
	researches userResearchBonuses
	planet     models.Planet
}

// getDefenderPlanet get the defender info before the attack
func (s *Service) getDefenderPlanet(ctx context.Context, missionEvent models.MissionEvent, storage TxStorages) (attackDefender, error) {
	defender := attackDefender{}
	if s.isPlanetNPC(missionEvent.PlanetTo.Z) {
		npcStats, err := s.registry.GetNPCStatsByPosition(missionEvent.PlanetTo.Z)
		if err != nil {
			return attackDefender{}, fmt.Errorf("s.registry.GetNPCStatsByPosition(): %w", err)
		}

		defender.fleet = lo.Map(npcStats.Fleet, func(fleetCount registry.FleetUnitCount, _ int) models.FleetUnit {
			return models.FleetUnit{
				ID:    fleetCount.ID,
				Count: fleetCount.Count,
			}
		})

		for _, researchID := range npcStats.Researches {
			research, err := s.registry.GetResearchStatsByID(researchID)
			if err != nil {
				return attackDefender{}, fmt.Errorf("s.registry.GetResearchStatsByID(): %w", err)
			}

			if research.Type == consts.ResearchTypeWeaponTechnology {
				defender.researches.attackBonus = float64(research.Bonuses.AttackPower)
			}

			if research.Type == consts.ResearchTypeArmorTechnology {
				defender.researches.defenseBonus = float64(research.Bonuses.ArmorPower)
			}
		}

		defender.planet.UserLogin = s.getNPCLoginByCoordinates(missionEvent.PlanetTo)
		defender.planet.Resources = models.Resources{
			Metal:   npcStats.Resources.Metal,
			Crystal: npcStats.Resources.Crystal,
			Gas:     npcStats.Resources.Gas,
		}
	} else {
		planet, err := storage.GetPlanetInfoByCoordinates(ctx, missionEvent.PlanetTo)
		if err != nil {
			return attackDefender{}, fmt.Errorf("storage.GetPlanetInfoByCoordinates(): %w", err)
		}

		defenderFleet, err := storage.GetPlanetFleetForUpdate(ctx, planet.ID)
		if err != nil {
			return attackDefender{}, fmt.Errorf("storage.GetPlanetFleetForUpdate(): %w", err)
		}

		defenderResearches, err := s.getResearchBonuses(ctx, missionEvent.UserID, storage)
		if err != nil {
			return attackDefender{}, fmt.Errorf("getResearchBonuses(defender): %w", err)
		}

		defender.planet = planet
		defender.fleet = defenderFleet
		defender.researches = defenderResearches
	}

	return defender, nil
}

func (s *Service) getNPCLoginByCoordinates(planetTo models.Coordinates) string {
	switch planetTo.Z {
	case consts.NPCTierOnePositionZ:
		return consts.NPCTierOneLogin
	case consts.NPCTierTwoPositionZ:
		return consts.NPCTierTwoLogin
	case consts.NPCTierThreePositionZ:
		return consts.NPCTierThreeLogin
	default:
		return ""
	}
}
