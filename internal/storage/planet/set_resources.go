package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *PlanetStorage) SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error {
	const updateResourcesQuery = `
		INSERT INTO session_beta.planet_resources (
			planet_id,
			metal,
			crystal,
			gas,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
		ON CONFLICT (planet_id) DO UPDATE SET
			metal = excluded.metal,
			crystal = excluded.crystal,
			gas = excluded.gas,
			updated_at = excluded.updated_at;
	`

	_, err := r.DB.Exec(
		ctx,
		updateResourcesQuery,
		planetID,
		updatedResources.Metal,
		updatedResources.Crystal,
		updatedResources.Gas,
		updatedResources.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
