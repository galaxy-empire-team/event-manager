package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) AddBoost(ctx context.Context, planetID uuid.UUID, boost models.Boost) error {
	const addBoostQuery = `
		INSERT INTO session_beta.user_boosts (user_id, boost_id, count, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			count = user_boosts.count + excluded.count,
			updated_at = excluded.updated_at;
	`

	_, err := s.DB.Exec(
		ctx,
		addBoostQuery,
		planetID,
		boost.ID,
		boost.Count,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
