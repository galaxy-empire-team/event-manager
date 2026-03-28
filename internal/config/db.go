package config

type PgConn struct {
	Host     string `envconfig:"HOST" required:"true"`
	Port     uint16 `envconfig:"PORT" required:"true"`
	Username string `envconfig:"USERNAME" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	DBName   string `envconfig:"DB_NAME" required:"true"`
	SSLMode  string `envconfig:"SSL_MODE" default:"disable"`
}
