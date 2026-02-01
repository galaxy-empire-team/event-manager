package building

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *BuildingStorage) DeleteBuildEvents(ctx context.Context, events []models.BuildEvent) error {
	if len(events) == 0 {
		return nil
	}

	const deleteEventQuery = `
		DELETE FROM session_beta.building_events WHERE id = ANY($1);
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
