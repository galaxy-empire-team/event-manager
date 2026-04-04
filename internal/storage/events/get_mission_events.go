package events

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *EventsStorage) GetMissionEvents(ctx context.Context, missionEventsCount uint16) ([]models.MissionEvent, error) {
	const getMissionEventsQuery = `
		SELECT
			id,
			mission_id,
			user_id,
			planet_from,
			planet_to_x,
			planet_to_y,
			planet_to_z,
			fleet,
			is_returning,
			started_at,
			finished_at
		FROM
			session_beta.event_missions
		WHERE
			finished_at <= NOW() + INTERVAL '1 SECOND'
		LIMIT $1
		FOR UPDATE SKIP LOCKED;
	`

	rows, err := r.DB.Query(ctx, getMissionEventsQuery, missionEventsCount)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Query(): %w", err)
	}
	defer rows.Close()

	var (
		missionEvents []models.MissionEvent
		fleet         []fleetUnit
	)

	for rows.Next() {
		var me models.MissionEvent

		err = rows.Scan(
			&me.ID,
			&me.MissionID,
			&me.UserID,
			&me.PlanetFrom,
			&me.PlanetTo.X,
			&me.PlanetTo.Y,
			&me.PlanetTo.Z,
			&fleet,
			&me.IsReturning,
			&me.StartedAt,
			&me.FinishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		me.Fleet = lo.Map(fleet, func(f fleetUnit, _ int) models.FleetUnit {
			return models.FleetUnit{
				ID:    f.ID,
				Count: f.Count,
			}
		})

		missionEvents = append(missionEvents, me)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err(): %w", rows.Err())
	}

	return missionEvents, nil
}
