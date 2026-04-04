package events

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *EventsStorage) DeleteFleetConstructionEvents(ctx context.Context, events []models.FleetConstructionEvent) error {
	if len(events) == 0 {
		return nil
	}

	const deleteEventQuery = `
		DELETE FROM session_beta.event_fleet_constructions WHERE id = ANY($1);
	`

	idsToDelete := make([]uint64, 0, len(events))
	for _, e := range events {
		idsToDelete = append(idsToDelete, e.ID)
	}

	_, err := r.DB.Exec(ctx, deleteEventQuery, idsToDelete)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
