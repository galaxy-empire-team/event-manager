package planet

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) GetPlanetInfoByCoordinates(ctx context.Context, planetFrom models.Coordinates) (models.Planet, error) {
	const getPlanetInfoQuery = `
		SELECT 
			p.id,
			p.user_id,
			u.login,
			p.x,
			p.y,
			p.z
		FROM session_beta.planets p
		JOIN session_beta.users u ON p.user_id = u.id
		WHERE p.x = $1 AND p.y = $2 AND p.z = $3;
	`

	var planetInfo models.Planet

	err := s.DB.QueryRow(ctx, getPlanetInfoQuery, planetFrom.X, planetFrom.Y, planetFrom.Z).Scan(
		&planetInfo.ID,
		&planetInfo.UserID,
		&planetInfo.UserLogin,
		&planetInfo.Coordinates.X,
		&planetInfo.Coordinates.Y,
		&planetInfo.Coordinates.Z,
	)
	if err != nil {
		return models.Planet{}, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return planetInfo, nil
}
