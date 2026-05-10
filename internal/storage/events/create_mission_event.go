package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *EventsStorage) CreateMissionEvent(ctx context.Context, missionEvent models.MissionEvent) error {
	const createEventQuery = `
		INSERT INTO session_beta.event_missions (
			mission_id,
			user_id,
			planet_from,
			planet_to_x, 
			planet_to_y, 
			planet_to_z, 
			fleet,
			cargo,
			is_returning,
			started_at,
			finished_at
		) VALUES (
			$1,    -- mission_id
			$2,    -- user_id
			$3,    -- planet_from
			$4,    -- planet_to_x
			$5,    -- planet_to_y
			$6,    -- planet_to_z
			$7,    -- fleet
			$8,    -- cargo
			$9,    -- is_returning
			$10,   -- started_at
			$11    -- finished_at
		)  
	`

	fleetJson, err := json.Marshal(toFleetUnits(missionEvent.Fleet))
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	cargoJson, err := json.Marshal(toResources(missionEvent.Cargo))
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	_, err = s.DB.Exec(ctx, createEventQuery,
		missionEvent.MissionID,
		missionEvent.UserID,
		missionEvent.PlanetFrom,
		missionEvent.PlanetTo.X,
		missionEvent.PlanetTo.Y,
		missionEvent.PlanetTo.Z,
		fleetJson,
		cargoJson,
		missionEvent.IsReturning,
		missionEvent.StartedAt,
		missionEvent.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
