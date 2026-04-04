package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server                  Server          `envconfig:"SERVER"`
	PgConn                  PgConn          `envconfig:"PG"`
	App                     App             `envconfig:"APP"`
	BuildingWorker          Worker          `envconfig:"WORKER_BUILD"`
	MissionWorker           Worker          `envconfig:"WORKER_MISSION"`
	ResearchWorker          Worker          `envconfig:"WORKER_RESEARCH"`
	FleetConstructionWorker Worker          `envconfig:"WORKER_FLEET_CONSTRUCTION"`
	BridgeAPIClient         BridgeAPIClient `envconfig:"BRIDGE_API_CLIENT"`
}

func New() (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("envconfig.Process(): %w", err)
	}

	return cfg, nil
}
