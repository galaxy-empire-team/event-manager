package research

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (s *Service) Process(ctx context.Context, researchEventsCount uint16) error {
	err := s.txManager.ExecResearchTx(ctx, func(ctx context.Context, researchStorages TxStorages) error {
		researchEvents, err := researchStorages.GetResearchEvents(ctx, researchEventsCount)
		if err != nil {
			return fmt.Errorf("researchStorages.GetResearchEvents(): %w", err)
		}

		s.logger.Info("Fetched research events", zap.Int("count", len(researchEvents)))

		if len(researchEvents) == 0 {
			return nil
		}

		for _, researchEvent := range researchEvents {
			nextLvlResearchID, err := s.registry.GetResearchNextLvlID(consts.ResearchID(researchEvent.ResearchID))
			if err != nil {
				return fmt.Errorf("s.registry.GetResearchNextLvlID(): %w", err)
			}

			updatedResearch := models.ResearchUpgrade{
				UserID:            researchEvent.UserID,
				CurrentResearchID: researchEvent.ResearchID,
				UpdatedResearchID: nextLvlResearchID,
			}

			err = researchStorages.SetResearchID(ctx, updatedResearch)
			if err != nil {
				return fmt.Errorf("researchStorages.SetResearchID(): %w", err)
			}
		}

		err = researchStorages.DeleteResearchEvents(ctx, researchEvents)
		if err != nil {
			return fmt.Errorf("researchStorages.DeleteResearchEvents(): %w", err)
		}

		s.logger.Info("Completed upgrading research", zap.Int("count", len(researchEvents)))

		return nil
	})
	if err != nil {
		return fmt.Errorf("txManager.ExecResearchTx(): %w", err)
	}

	return nil
}
