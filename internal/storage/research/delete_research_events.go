package research

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *ResearchStorage) DeleteResearchEvents(ctx context.Context, events []models.ResearchEvent) error {
	if len(events) == 0 {
		return nil
	}

	const deleteEventQuery = `
		DELETE FROM session_beta.event_researches WHERE id = ANY($1);
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
