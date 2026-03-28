package research

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *ResearchStorage) GetUserResearches(ctx context.Context, userID uuid.UUID) ([]consts.ResearchID, error) {
	const getAllUserResearchesQuery = `
		SELECT 
			research_id
		FROM session_beta.user_researches
		WHERE user_id = $1;
	`

	var researches []consts.ResearchID
	rows, err := r.DB.Query(ctx, getAllUserResearchesQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var researchID consts.ResearchID
		err = rows.Scan(&researchID)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		researches = append(researches, researchID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return researches, nil
}
