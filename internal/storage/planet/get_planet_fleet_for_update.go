package planet

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) GetPlanetFleetForUpdate(ctx context.Context, planetID uuid.UUID) ([]models.FleetUnit, error) {
	const getPlanetFleetQuery = `
		SELECT 
			fleet_id, 
			count
		FROM session_beta.planet_fleet
		WHERE planet_id = $1
		FOR UPDATE;
	`

	var fleet []models.FleetUnit
	rows, err := s.DB.Query(ctx, getPlanetFleetQuery, planetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fleetUnit models.FleetUnit

		err = rows.Scan(&fleetUnit.ID, &fleetUnit.Count)
		if err != nil {
			return nil, err
		}

		fleet = append(fleet, fleetUnit)
	}

	return fleet, nil
}
