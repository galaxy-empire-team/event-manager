package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// GetAllBuildings retrieves all buildings information from the target planet.
func (s *PlanetStorage) GetAllBuildings(ctx context.Context, planetID uuid.UUID) ([]consts.BuildingID, error) {
	const getAllBuildingsQuery = `
		SELECT 
			building_id
		FROM session_beta.planet_buildings
		WHERE planet_id = $1;
	`

	rows, err := s.DB.Query(ctx, getAllBuildingsQuery, planetID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query(): %w", err)
	}

	var buildingIDs []consts.BuildingID

	for rows.Next() {
		var buildingID consts.BuildingID

		err = rows.Scan(&buildingID)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		buildingIDs = append(buildingIDs, buildingID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return buildingIDs, nil
}
