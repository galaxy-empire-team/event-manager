package txmanager

import (
	"github.com/jackc/pgx/v5"

	eventsstorage "github.com/galaxy-empire-team/event-manager/internal/storage/events"
	notificationstorage "github.com/galaxy-empire-team/event-manager/internal/storage/notification"
	planetstorage "github.com/galaxy-empire-team/event-manager/internal/storage/planet"
	researchstorage "github.com/galaxy-empire-team/event-manager/internal/storage/research"
)

// I don't want to write boilerplate stuff, embed all storages ^_^.
type StorageSet struct {
	*eventsstorage.EventsStorage
	*planetstorage.PlanetStorage
	*notificationstorage.NotificationStorage
	*researchstorage.ResearchStorage
}

func newStorageSet(tx pgx.Tx) StorageSet {
	return StorageSet{
		EventsStorage:       eventsstorage.New(tx),
		PlanetStorage:       planetstorage.New(tx),
		NotificationStorage: notificationstorage.New(tx),
		ResearchStorage:     researchstorage.New(tx),
	}
}
