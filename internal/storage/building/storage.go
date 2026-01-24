package building

import (
	"github.com/galaxy-empire-team/event-manager/internal/db"
)

// Embed txManager requires different naming -> can't use 'storage' storage name :()
type BuildingStorage struct {
	DB db.DB
}

func New(db db.DB) *BuildingStorage {
	return &BuildingStorage{
		DB: db,
	}
}
