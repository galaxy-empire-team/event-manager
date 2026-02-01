package building

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *BuildingStorage) UpgradeBuildingLevel(ctx context.Context, building models.PlanetBuilding) error {
	const upgradeBuildingQuery = `
		WITH next_building_id AS (
			SELECT b.id FROM session_beta.buildings b
			WHERE b.building_type = $1 AND b.level = $2
		)
		UPDATE session_beta.planet_buildings pb
		SET
			building_id = (SELECT id FROM next_building_id),
			updated_at = NOW(),
			finished_at = NULL
		WHERE
			planet_id = $3
		AND 
			building_id = $4;
	`

	_, err := r.DB.Exec(ctx, upgradeBuildingQuery,
		building.BuildType,
		building.Level,
		building.PlanetID,
		building.ID,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
