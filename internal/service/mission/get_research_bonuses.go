package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// getResearchBonuses return weapon and armor tech bonuses for an attack calculation
func (s *Service) getResearchBonuses(ctx context.Context, userID uuid.UUID, storage TxStorages) (userResearchBonuses, error) {
	researchIDs, err := storage.GetUserResearchesByTypes(ctx, userID, []consts.ResearchType{consts.ResearchTypeWeaponTechnology, consts.ResearchTypeArmorTechnology})
	if err != nil {
		return userResearchBonuses{}, fmt.Errorf("storage.GetUserResearchesByTypes(): %w", err)
	}

	weaponTechID, ok := researchIDs[consts.ResearchTypeWeaponTechnology]
	if !ok {
		weaponTechID, err = s.registry.GetResearchZeroLvlIDByType(consts.ResearchTypeWeaponTechnology)
		if err != nil {
			return userResearchBonuses{}, fmt.Errorf("registry.GetResearchZeroLvlIDByType(weapon tech): %w", err)
		}
	}

	weaponTechStats, err := s.registry.GetResearchStatsByID(weaponTechID)
	if err != nil {
		return userResearchBonuses{}, fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
	}

	armorTechID, ok := researchIDs[consts.ResearchTypeArmorTechnology]
	if !ok {
		armorTechID, err = s.registry.GetResearchZeroLvlIDByType(consts.ResearchTypeArmorTechnology)
		if err != nil {
			return userResearchBonuses{}, fmt.Errorf("registry.GetResearchZeroLvlIDByType(armor tech): %w", err)
		}
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
