package txmanager

import (
	"github.com/jackc/pgx/v5"

	buildingstorage "github.com/galaxy-empire-team/event-manager/internal/storage/building"
)

// I don't want to write boilerplate stuff, embed all storages ^_^
type StorageSet struct {
	*buildingstorage.BuildingStorage
}

func newStorageSet(tx pgx.Tx) StorageSet {
	return StorageSet{
		BuildingStorage: buildingstorage.New(tx),
	}
}
