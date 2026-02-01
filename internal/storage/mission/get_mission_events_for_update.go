package mission

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *MissionStorage) GetMissionEventsForUpdate(ctx context.Context) ([]models.MissionEvent, error) {
	const getMissionEventsQuery = `
		SELECT
			id,
			mission_type,
			user_id,
			planet_from,
			planet_to_x,
			planet_to_y,
			planet_to_z,
			started_at,
			finished_at
		FROM
			session_beta.mission_events
		WHERE
			finished_at <= NOW() + INTERVAL '1 SECOND'
		FOR UPDATE SKIP LOCKED;
	`

	rows, err := r.DB.Query(ctx, getMissionEventsQuery)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Query(): %w", err)
	}
	defer rows.Close()

	var missionEvents []models.MissionEvent
	for rows.Next() {
		var me models.MissionEvent

		err = rows.Scan(
			&me.ID,
			&me.Type,
			&me.UserID,
			&me.PlanetFrom,
			&me.PlanetTo.X,
			&me.PlanetTo.Y,
			&me.PlanetTo.Z,
			&me.StartedAt,
			&me.FinishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		missionEvents = append(missionEvents, me)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err(): %w", rows.Err())
	}

	return missionEvents, nil
}
