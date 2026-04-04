package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *PlanetStorage) GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error) {
	const getResourcesQuery = `
		SELECT 
			metal,
			crystal,
			gas,
			updated_at
		FROM session_beta.planet_resources
		WHERE planet_id = $1
		FOR UPDATE;
	`

	var resources models.Resources

	err := r.DB.QueryRow(ctx, getResourcesQuery, planetID).Scan(
		&resources.Metal,
		&resources.Crystal,
		&resources.Gas,
		&resources.UpdatedAt,
	)
	if err != nil {
		return models.Resources{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return resources, nil
}
