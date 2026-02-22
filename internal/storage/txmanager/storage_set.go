package txmanager

import (
	"github.com/jackc/pgx/v5"

	buildingstorage "github.com/galaxy-empire-team/event-manager/internal/storage/building"
	missionstorage "github.com/galaxy-empire-team/event-manager/internal/storage/mission"
	notificationstorage "github.com/galaxy-empire-team/event-manager/internal/storage/notification"
	planetstorage "github.com/galaxy-empire-team/event-manager/internal/storage/planet"
)

// I don't want to write boilerplate stuff, embed all storages ^_^.
type StorageSet struct {
	*buildingstorage.BuildingStorage
	*missionstorage.MissionStorage
	*planetstorage.PlanetStorage
	*notificationstorage.NotificationStorage
}

func newStorageSet(tx pgx.Tx) StorageSet {
	return StorageSet{
		BuildingStorage:     buildingstorage.New(tx),
		MissionStorage:      missionstorage.New(tx),
		PlanetStorage:       planetstorage.New(tx),
		NotificationStorage: notificationstorage.New(tx),
	}
}
