package mission

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *MissionStorage) ColonizePlanet(ctx context.Context, colonizeEvent models.MissionEvent) error {
	planetToColonize := fromMissionEvent(colonizeEvent)

	planetID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("uuid.NewV7(): %w", err)
	}

	planetToColonize.ID = planetID

	err = r.createPlanetRow(ctx, planetToColonize)
	if err != nil {
		return fmt.Errorf("r.createPlanetRow(): %w", err)
	}

	err = r.createResourcesRow(ctx, planetToColonize.ID)
	if err != nil {
		return fmt.Errorf("r.createResourcesRow(): %w", err)
	}

	return nil
}

func (r *MissionStorage) createPlanetRow(ctx context.Context, planet planetToColonize) error {
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
		);
	`

	_, err := r.DB.Exec(
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				if pgErr.ConstraintName == "planet_have_unique_x_y_z" {
					return models.ErrPlanetCoordinatesAlreadyTaken
				}

				return models.ErrCapitolAlreadyExists
			}
		}

		return fmt.Errorf("r.DB.Exec(): %w", err)
	}

	return nil
}

func (r *MissionStorage) createResourcesRow(ctx context.Context, planetID uuid.UUID) error {
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
