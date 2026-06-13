package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) AddDebris(ctx context.Context, planetID uuid.UUID, debris models.Resources) error {
	const addDebrisQuery = `
		INSERT INTO session_beta.planet_debris pd (planet_id, metal, crystal, updated_at) 
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (planet_id) DO UPDATE SET
			metal = pd.metal + EXCLUDED.metal,
			crystal = pd.crystal + EXCLUDED.crystal,
			updated_at = EXCLUDED.updated_at;
	`

	_, err := s.DB.Exec(
		ctx,
		addDebrisQuery,
		planetID,
		debris.Metal,
		debris.Crystal,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
