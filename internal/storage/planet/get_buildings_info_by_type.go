package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

// GetBuildingsInfoByType retrieves mine infromation from the target planet's buildings
func (s *PlanetStorage) GetBuildingsInfoByTypes(ctx context.Context, planetID uuid.UUID, BuildingTypes []consts.BuildingType) (map[consts.BuildingType]models.BuildingInfo, error) {
	const getMineInfoQuery = `
		SELECT 
			b.id,
			b.building_type,
			b.production_s
		FROM session_beta.planet_buildings pb
		JOIN session_beta.buildings b ON pb.building_id = b.id
		WHERE pb.planet_id = $1 AND b.building_type = ANY($2);
	`

	rows, err := s.DB.Query(ctx, getMineInfoQuery, planetID, BuildingTypes)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}

	buildingsInfo := make(map[consts.BuildingType]models.BuildingInfo)
	for rows.Next() {
		var buildingInfo models.BuildingInfo
		err = rows.Scan(
			&buildingInfo.ID,
			&buildingInfo.Type,
			&buildingInfo.ProductionS,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		buildingsInfo[buildingInfo.Type] = buildingInfo
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return buildingsInfo, nil
}
