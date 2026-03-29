package bridgeapi

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	planetpb "github.com/galaxy-empire-team/bridge-api/api/gen/go/planet/v1"
)

func (c *Client) UpdatePlanetResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedTime time.Time) error {
	_, err := c.planetServiceClient.UpdatePlanetResources(ctx, &planetpb.UpdatePlanetResourcesRequest{
		UserId:   userID.String(),
		PlanetId: planetID.String(),
		Time:     timestamppb.New(updatedTime),
	})
	if err != nil {
		return fmt.Errorf("planetServiceClient.UpdatePlanetResources(): %w", err)
	}

	return nil
}
