package building

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *BuildingStorage) SetBuildingID(ctx context.Context, building models.BuildingUpgrade) error {
	const createBuildingQuery = `
		UPDATE session_beta.planet_buildings
		SET 
			building_id = $3,
			updated_at = NOW(),
			finished_at = NULL
		WHERE planet_id = $1 AND building_id = $2;
	`

	_, err := r.DB.Exec(ctx, createBuildingQuery,
		building.PlanetID,
		building.CurrentBuildingID,
		building.UpdatedBuildingID,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
