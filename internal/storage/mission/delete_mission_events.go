package mission

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *MissionStorage) DeleteMissionEvents(ctx context.Context, events []models.MissionEvent) error {
	if len(events) == 0 {
		return nil
	}

	const deleteMissionEventQuery = `
		DELETE FROM session_beta.mission_events WHERE id = ANY($1);
	`

	idsToDelete := make([]uint64, 0, len(events))
	for _, e := range events {
		idsToDelete = append(idsToDelete, e.ID)
	}

	_, err := r.DB.Exec(ctx, deleteMissionEventQuery, idsToDelete)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
