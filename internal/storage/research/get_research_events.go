package research

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *ResearchStorage) GetResearchEvents(ctx context.Context, researchEventsCount uint16) ([]models.ResearchEvent, error) {
	const getResearchEventsQuery = `
		SELECT
			id,
			user_id,
			research_id,
			started_at,
			finished_at
		FROM
			session_beta.event_researches
		WHERE
			finished_at <= NOW() + INTERVAL '1 SECOND'
		LIMIT $1
		FOR UPDATE SKIP LOCKED;
	`

	rows, err := r.DB.Query(ctx, getResearchEventsQuery, researchEventsCount)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Query(): %w", err)
	}
	defer rows.Close()

	var researchEvents []models.ResearchEvent
	for rows.Next() {
		var re models.ResearchEvent

		err = rows.Scan(
			&re.ID,
			&re.UserID,
			&re.ResearchID,
			&re.StartedAt,
			&re.FinishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		researchEvents = append(researchEvents, re)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err(): %w", rows.Err())
	}

	return researchEvents, nil
}
