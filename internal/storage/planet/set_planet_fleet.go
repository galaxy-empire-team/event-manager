package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) SetPlanetFleet(ctx context.Context, planetID uuid.UUID, fleet []models.FleetUnit) error {
	const setPlanetFleetQuery = `
		UPDATE session_beta.planet_fleet AS pf
		SET 
			count = x.count,
			updated_at = NOW()
		FROM jsonb_to_recordset($1) AS x (planet_id uuid, fleet_id int, count int)
		WHERE pf.planet_id = x.planet_id AND pf.fleet_id = x.fleet_id;
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
