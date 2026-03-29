package models

import (
	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type ResearchUpgrade struct {
	UserID            uuid.UUID
	CurrentResearchID consts.ResearchID
	UpdatedResearchID consts.ResearchID
}
