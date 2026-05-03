package bridgeapi

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	planetpb "github.com/galaxy-empire-team/bridge-api/api/gen/go/planet/v1"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (c *Client) ColonizePlanet(ctx context.Context, userID uuid.UUID, event models.MissionEvent) error {
	_, err := c.planetServiceClient.ColonizePlanet(ctx, &planetpb.ColonizePlanetRequest{
		UserID:      userID.String(),
		OperationID: event.ID,
		IsCapitol:   false,
		Coordinates: &planetpb.Coordinates{
			X: uint32(event.PlanetTo.X),
			Y: uint32(event.PlanetTo.Y),
			Z: uint32(event.PlanetTo.Z),
		},
		Resources: &planetpb.Resources{
			Metal:   event.Resources.Metal,
			Crystal: event.Resources.Crystal,
			Gas:     event.Resources.Gas,
		},
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.AlreadyExists {
			return models.ErrPlanetCoordinatesAlreadyTaken
		}

		return fmt.Errorf("planetServiceClient.ColonizePlanet(): %w", err)
	}

	return nil
}
