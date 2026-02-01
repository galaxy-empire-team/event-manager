package building

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *BuildingStorage) GetCurrentBuilding(ctx context.Context, building models.BuildEvent) (models.PlanetBuilding, error) {
	const getCurrentBuildingsQuery = `
		SELECT 
			p.planet_id,
			b.id,
			b.building_type,
			b.level
		FROM session_beta.planet_buildings p
		JOIN session_beta.buildings b ON p.building_id = b.id
		WHERE
			p.planet_id = $1
		AND
			b.building_type = $2;
	`

	var pb models.PlanetBuilding
	err := r.DB.QueryRow(ctx, getCurrentBuildingsQuery,
		building.PlanetID,
		building.BuildType,
	).Scan(
		&pb.PlanetID,
		&pb.ID,
		&pb.BuildType,
		&pb.Level,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.PlanetBuilding{}, models.ErrBuildingNotFound
		}
		return models.PlanetBuilding{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return pb, nil
}
