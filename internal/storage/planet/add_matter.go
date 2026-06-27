package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *PlanetStorage) AddMatter(ctx context.Context, userID uuid.UUID, matter uint64) error {
	const addMatterQuery = `
		INSERT INTO session_beta.user_resources (user_id, matter, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			matter = user_resources.matter + excluded.matter,
			updated_at = excluded.updated_at;
	`

	_, err := s.DB.Exec(
		ctx,
		addMatterQuery,
		userID,
		matter,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
