package config

type BridgeAPIClient struct {
	Endpoint string `envconfig:"ENDPOINT" required:"true"`
}
