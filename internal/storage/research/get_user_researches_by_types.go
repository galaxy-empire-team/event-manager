package research

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *ResearchStorage) GetUserResearchesByTypes(ctx context.Context, userID uuid.UUID, researchTypes []consts.ResearchType) (map[consts.ResearchType]consts.ResearchID, error) {
	const getAllUserResearchesQuery = `
		SELECT
			ur.research_id,
			rs.research_type
		FROM session_beta.user_researches ur
		JOIN session_beta.s_researches rs ON ur.research_id = rs.id
		WHERE user_id = $1 AND rs.research_type = ANY($2);
	`

	var researches = make(map[consts.ResearchType]consts.ResearchID)

	rows, err := r.DB.Query(ctx, getAllUserResearchesQuery, userID, researchTypes)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var researchID consts.ResearchID
		var researchType consts.ResearchType

		err = rows.Scan(&researchID, &researchType)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		researches[researchType] = researchID
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return researches, nil
}
