package building

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *BuildingStorage) UpgradeBuilding(ctx context.Context, building models.BuildEvent) error {
	const upgradeBuildingQuery = `
		WITH next_building_id AS (
			SELECT 
				b.id AS old_id,
				b1.id AS new_id
            FROM session_beta.planet_buildings pb
            JOIN session_beta.buildings b ON pb.building_id = b.id
            JOIN session_beta.buildings b1 ON b.level + 1 = b1.level
			WHERE
			    pb.planet_id = $1
			AND
			    b.type = $2
		)
		UPDATE
			session_beta.planet_buildings pb
		SET
		    building_id = (SELECT new_id FROM next_building_id),
			updated_at = NOW(),
			finished_at = NULL
		WHERE
			pb.planet_id = $1 
		AND 
			pb.building_id = (SELECT old_id FROM next_building_id);
	`

	_, err := r.DB.Exec(ctx, upgradeBuildingQuery,
		building.PlanetID,
		building.BuildType,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
