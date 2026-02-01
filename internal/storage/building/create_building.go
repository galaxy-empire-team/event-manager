package building

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *BuildingStorage) CreateBuilding(ctx context.Context, building models.PlanetBuilding) error {
	const createBuildingQuery = `
		INSERT INTO session_beta.planet_buildings (
			planet_id, 
			building_id, 
			created_at, 
			updated_at,
			finished_at
		) VALUES (
			$3,
			(SELECT id FROM session_beta.buildings WHERE building_type = $1 AND level = $2),
			NOW(),
			NOW(),
			NULL
		);
	`

	_, err := r.DB.Exec(ctx, createBuildingQuery,
		building.BuildType,
		building.Level,
		building.PlanetID,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
