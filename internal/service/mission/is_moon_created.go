package mission

import (
	"github.com/galaxy-empire-team/event-manager/internal/models"
)

const (
	maxMoonCreationChance = 20
	onePercentDebreeCount = 250_000
)

func (s *Service) isMoonCreated(debris models.Resources) bool {
	const maxChance = 100

	totalDebrisChance := min(int(float64(debris.Metal+debris.Crystal)/onePercentDebreeCount), maxMoonCreationChance)

	return s.randGenerator.Intn(maxChance) < totalDebrisChance
}
