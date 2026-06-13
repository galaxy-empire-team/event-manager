package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) GetPlanetIDByCoordinates(ctx context.Context, coordinates models.Coordinates) (uuid.UUID, error) {
	const getPlanetIDQuery = `
		SELECT 
			id
		FROM session_beta.planets
		WHERE x = $1 AND y = $2 AND z = $3;
	`

	var planetID uuid.UUID
	err := s.DB.QueryRow(ctx, getPlanetIDQuery, coordinates.X, coordinates.Y, coordinates.Z).Scan(&planetID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return planetID, nil
}
