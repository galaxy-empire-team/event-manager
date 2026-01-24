package building

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *BuildingStorage) GetBuildingEvents(ctx context.Context) ([]models.BuildEvent, error) {
	const createCapitolQuery = `
		SELECT
			id,
			planet_id,
			build_type,
			started_at,
			finished_at
		FROM
			session_beta.building_events
		WHERE
			finished_at <= NOW() + INTERVAL '1 SECOND'
		FOR UPDATE SKIP LOCKED;
	`

	rows, err := r.DB.Query(ctx, createCapitolQuery)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Query(): %w", err)
	}
	defer rows.Close()

	var buildEvents []models.BuildEvent
	for rows.Next() {
		var be models.BuildEvent

		err = rows.Scan(
			&be.ID,
			&be.PlanetID,
			&be.BuildType,
			&be.StartdAt,
			&be.FinishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		buildEvents = append(buildEvents, be)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err(): %w", rows.Err())
	}

	return buildEvents, nil
}
