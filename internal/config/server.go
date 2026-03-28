package config

type Server struct {
	Endpoint string `envconfig:"ENDPOINT" default:"localhost:8000"`
}
