package research

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *ResearchStorage) SetResearchID(ctx context.Context, research models.ResearchUpgrade) error {
	const updateResearchQuery = `
		WITH d AS (
			DELETE FROM session_beta.user_researches
			WHERE user_id = $1 AND research_id = $2
		)
		INSERT INTO session_beta.user_researches (user_id, research_id, updated_at)
		VALUES ($1, $3, NOW())
		ON CONFLICT (user_id, research_id) DO NOTHING;
	`

	_, err := r.DB.Exec(ctx, updateResearchQuery,
		research.UserID,
		research.CurrentResearchID,
		research.UpdatedResearchID,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
