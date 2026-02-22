package planet

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) GetPlanetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error) {
	const getResourecsForUpdateQuery = `
		SELECT
		 	metal,
		 	crystal,
		 	gas 
		FROM session_beta.planet_resources
		WHERE planet_id = $1
		FOR UPDATE;
	`

	var resources models.Resources
	err := s.DB.QueryRow(ctx, getResourecsForUpdateQuery, planetID).Scan(
		&resources.Metal,
		&resources.Crystal,
		&resources.Gas,
	)
	if err != nil {
		return models.Resources{}, err
	}

	return resources, nil
}
