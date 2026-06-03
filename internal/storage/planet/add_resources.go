package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) AddResources(ctx context.Context, planetID uuid.UUID, resources models.Resources) error {
	const addResourcesQuery = `
		UPDATE session_beta.planet_resources
		SET 
			metal = metal + $2,
			crystal = crystal + $3,
			gas = gas + $4
		WHERE planet_id = $1
	`

	_, err := s.DB.Exec(
		ctx,
		addResourcesQuery,
		planetID,
		resources.Metal,
		resources.Crystal,
		resources.Gas,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
