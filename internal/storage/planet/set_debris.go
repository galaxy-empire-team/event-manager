package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) SetDebris(ctx context.Context, planetID uuid.UUID, debris models.Resources) error {
	const setDebrisQuery = `
		UPDATE session_beta.planet_debris 
		SET metal = $2, crystal = $3 
		WHERE planet_id = $1;
	`

	_, err := s.DB.Exec(ctx, setDebrisQuery, planetID, debris.Metal, debris.Crystal)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
