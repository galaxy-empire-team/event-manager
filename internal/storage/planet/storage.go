package planet

import (
	"github.com/galaxy-empire-team/event-manager/internal/db"
)

// Embed txManager requires different naming -> can't use 'storage' storage name :().
type PlanetStorage struct {
	DB db.DB
}

func New(db db.DB) *PlanetStorage {
	return &PlanetStorage{
		DB: db,
	}
}
