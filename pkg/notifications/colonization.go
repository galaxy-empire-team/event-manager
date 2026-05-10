package notifications

type ColonizationV1 struct {
	Planet Coordinates `json:"planet"`
	Err    string      `json:"error,omitempty"`
}
