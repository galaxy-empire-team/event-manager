package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) GetPlanetCoordinatesByID(ctx context.Context, planetID uuid.UUID) (models.Coordinates, error) {
	const getPlanetCoordinatesQuery = `
		SELECT 
			x,
			y,
			z
		FROM session_beta.planets
		WHERE id = $1;
	`

	var coordinates models.Coordinates
	err := s.DB.QueryRow(ctx, getPlanetCoordinatesQuery, planetID).Scan(&coordinates.X, &coordinates.Y, &coordinates.Z)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return coordinates, nil
}
