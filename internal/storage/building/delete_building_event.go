package building

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *BuildingStorage) DeleteBuildingEvent(ctx context.Context, event models.BuildEvent) error {
	const deleteEventQuery = `
		DELETE FROM session_beta.building_events WHERE id = $1;
	`

	_, err := r.DB.Exec(ctx, deleteEventQuery, event.ID)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
