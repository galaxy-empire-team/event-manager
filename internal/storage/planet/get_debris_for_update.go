package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) GetDebrisForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error) {
	const getDebrisQuery = `
		SELECT metal, crystal 
		FROM session_beta.planet_debris 
		WHERE planet_id = $1
		FOR UPDATE;
	`

	var debris models.Resources
	err := s.DB.QueryRow(ctx, getDebrisQuery, planetID).Scan(&debris.Metal, &debris.Crystal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Resources{}, nil
		}

		return models.Resources{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return debris, nil
}
