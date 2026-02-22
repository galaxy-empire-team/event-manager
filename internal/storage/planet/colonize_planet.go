package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *PlanetStorage) ColonizePlanet(ctx context.Context, colonizeEvent models.MissionEvent) (bool, error) {
	planetToColonize := fromMissionEvent(colonizeEvent)

	planetID, err := uuid.NewV7()
	if err != nil {
		return false, fmt.Errorf("uuid.NewV7(): %w", err)
	}

	planetToColonize.ID = planetID

	inserted, err := r.createPlanetRow(ctx, planetToColonize)
	if err != nil {
		return false, fmt.Errorf("r.createPlanetRow(): %w", err)
	}
	if !inserted {
		return false, nil
	}

	err = r.createResourcesRow(ctx, planetToColonize.ID)
	if err != nil {
		return false, fmt.Errorf("r.createResourcesRow(): %w", err)
	}

	return true, nil
}

func (r *PlanetStorage) createPlanetRow(ctx context.Context, planet planetToColonize) (bool, error) {
	const createPlanetQuery = `
		INSERT INTO session_beta.planets (
			id,
			user_id,
			x,
			y,
			z,
			has_moon,
			is_capitol,
			colonized_at
		) VALUES (
			$1,   --- planet.ID
			$2,   --- userID
			$3,	  --- planet.X
			$4,   --- planet.Y
			$5,   --- planet.Z
			$6,   --- planet.HasMoon
			$7,   --- planet.IsCapitol
			NOW() --- colonized_at
		) ON CONFLICT (x, y, z) DO NOTHING;
	`

	cmd, err := r.DB.Exec(
		ctx,
		createPlanetQuery,
		planet.ID,
		planet.UserID,
		planet.Coordinates.X,
		planet.Coordinates.Y,
		planet.Coordinates.Z,
		planet.HasMoon,
		planet.IsCapitol,
	)
	if err != nil {
		return false, fmt.Errorf("r.DB.Exec(): %w", err)
	}

	return cmd.RowsAffected() != 0, nil
}

func (r *PlanetStorage) createResourcesRow(ctx context.Context, planetID uuid.UUID) error {
	const createResourcesQuery = `
		INSERT INTO session_beta.planet_resources (
			planet_id,
			metal,
			crystal,
			gas,
			updated_at
		) VALUES (
			$1,    --- planet.ID
			0,     --- metal
			0,     --- crystal
			0,     --- gas
			NOW()  --- updated_at
		);
	`

	_, err := r.DB.Exec(
		ctx,
		createResourcesQuery,
		planetID,
	)
	if err != nil {
		return fmt.Errorf("r.DB.Exec(): %w", err)
	}

	return nil
}
