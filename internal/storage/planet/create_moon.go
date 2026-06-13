package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *PlanetStorage) CreateMoon(ctx context.Context, planetID uuid.UUID) error {
	const createMoonQuery = `
		UPDATE session_beta.planets
		SET has_moon = true, 
			updated_at = NOW()
		WHERE id = $1 AND has_moon = false;
	`

	_, err := s.DB.Exec(ctx, createMoonQuery, planetID)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
