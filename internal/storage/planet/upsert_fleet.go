package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) UpsertFleet(ctx context.Context, planetID uuid.UUID, fleet []models.FleetUnit) error {
	const setPlanetFleetQuery = `
		INSERT INTO session_beta.planet_fleet AS pf (planet_id, fleet_id, count, updated_at)
			SELECT x.planet_id, x.fleet_id, x.count, NOW()
			FROM jsonb_to_recordset($1) AS x (planet_id uuid, fleet_id int, count int)
		ON CONFLICT (planet_id, fleet_id) DO UPDATE SET 
		    count = pf.count + EXCLUDED.count, 
			updated_at = EXCLUDED.updated_at;
		`

	fleetData := make([]planetFleetUnit, 0, len(fleet))
	for _, fleetUnit := range fleet {
		fleetData = append(fleetData, planetFleetUnit{
			PlanetID: planetID,
			FleetID:  fleetUnit.ID,
			Count:    fleetUnit.Count,
		})
	}

	_, err := s.DB.Exec(ctx, setPlanetFleetQuery, fleetData)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
