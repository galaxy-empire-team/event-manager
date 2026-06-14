package notifications

type TransportV1 struct {
	PlanetFrom Coordinates `json:"planet_from"`
	PlanetTo   Coordinates `json:"planet_to"`
	Status     string      `json:"status"`
}
