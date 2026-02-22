package planet

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *PlanetStorage) GetPlanetInfoByID(ctx context.Context, planetID uuid.UUID) (models.Planet, error) {
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
		WHERE p.id = $1;
	`

	var planetInfo models.Planet

	err := s.DB.QueryRow(ctx, getPlanetInfoQuery, planetID).Scan(
		&planetInfo.ID,
		&planetInfo.UserID,
		&planetInfo.UserLogin,
		&planetInfo.Coordinates.X,
		&planetInfo.Coordinates.Y,
		&planetInfo.Coordinates.Z,
	)
	if err != nil {
		return models.Planet{}, err
	}

	return planetInfo, nil
}
