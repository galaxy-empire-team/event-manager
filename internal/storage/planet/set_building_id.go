package planet

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *PlanetStorage) SetBuildingID(ctx context.Context, building models.BuildingUpgrade) error {
	const createBuildingQuery = `
		WITH d AS (
			DELETE FROM session_beta.planet_buildings
			WHERE planet_id = $1 AND building_id = $2
		)
		INSERT INTO session_beta.planet_buildings (planet_id, building_id, updated_at)
		VALUES ($1, $3, NOW())
		ON CONFLICT (planet_id, building_id) DO NOTHING;
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
