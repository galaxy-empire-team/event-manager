package research

import (
	"context"

	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

type TxStorages interface {
	GetResearchEvents(ctx context.Context, researchEventsCount uint16) ([]models.ResearchEvent, error)
	DeleteResearchEvents(ctx context.Context, events []models.ResearchEvent) error
	SetResearchID(ctx context.Context, research models.ResearchUpgrade) error
}

type txManager interface {
	ExecResearchTx(ctx context.Context, fn func(ctx context.Context, researchStorage TxStorages) error) error
}

type registryProvider interface {
	GetResearchNextLvlID(researchID consts.ResearchID) (consts.ResearchID, error)
}

type Service struct {
	txManager txManager
	registry  registryProvider
	logger    *zap.Logger
}

func New(txManager txManager, registry registryProvider, logger *zap.Logger) *Service {
	return &Service{
		txManager: txManager,
		registry:  registry,
		logger:    logger,
	}
}
