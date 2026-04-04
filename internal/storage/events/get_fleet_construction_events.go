package events

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *EventsStorage) GetFleetConstructionEvents(ctx context.Context, fleetConstructionEventsCount uint16) ([]models.FleetConstructionEvent, error) {
	const getFleetConstructionEventsQuery = `
		SELECT
			id,
			planet_id,
			fleet_id,
			count,
			started_at,
			finished_at
		FROM
			session_beta.event_fleet_constructions
		WHERE
			finished_at <= NOW() + INTERVAL '1 SECOND'
		LIMIT $1
		FOR UPDATE SKIP LOCKED;
	`

	rows, err := r.DB.Query(ctx, getFleetConstructionEventsQuery, fleetConstructionEventsCount)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Query(): %w", err)
	}
	defer rows.Close()

	var fleetConstructionEvents []models.FleetConstructionEvent

	for rows.Next() {
		var fe models.FleetConstructionEvent

		err = rows.Scan(
			&fe.ID,
			&fe.PlanetID,
			&fe.FleetID,
			&fe.Count,
			&fe.StartedAt,
			&fe.FinishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		fleetConstructionEvents = append(fleetConstructionEvents, fe)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err(): %w", rows.Err())
	}

	return fleetConstructionEvents, nil
}
