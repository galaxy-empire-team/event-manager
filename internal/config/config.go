package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PgConn PgConn `envconfig:"PG"`
	App    App    `envconfig:"APP"`
}

type App struct {
	Logger `envconfig:"LOGGER"`
}

type Logger struct {
	Level  string `envconfig:"LEVEL"  default:"info"`
	Format string `envconfig:"FORMAT" default:"json"`
}

type PgConn struct {
	Host     string `envconfig:"HOST"     required:"true"`
	Port     uint16 `envconfig:"PORT"     required:"true"`
	Username string `envconfig:"USERNAME" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	DBName   string `envconfig:"DB_NAME"  required:"true"`
	SSLMode  string `envconfig:"SSL_MODE" default:"disable"`
}

func New() (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("envconfig.Process(): %w", err)
	}

	return cfg, nil
}
