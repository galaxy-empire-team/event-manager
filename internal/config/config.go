package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server         Server `envconfig:"SERVER"`
	PgConn         PgConn `envconfig:"PG"`
	App            App    `envconfig:"APP"`
	BuildingWorker Worker `envconfig:"WORKER_BUILDINGS"`
	MissionWorker  Worker `envconfig:"WORKER_MISSIONS"`
}

func New() (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("envconfig.Process(): %w", err)
	}

	return cfg, nil
}
