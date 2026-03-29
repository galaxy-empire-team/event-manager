package bridgeapi

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	planetpb "github.com/galaxy-empire-team/bridge-api/api/gen/go/planet/v1"
	"github.com/galaxy-empire-team/event-manager/internal/config"
)

type Client struct {
	conn *grpc.ClientConn

	planetServiceClient planetpb.PlanetServiceClient
}

func New(cfg config.BridgeAPIClient) (*Client, error) {
	conn, err := grpc.NewClient(cfg.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient(): %w", err)
	}

	return &Client{
		conn:                conn,
		planetServiceClient: planetpb.NewPlanetServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
