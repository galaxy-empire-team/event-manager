package mission

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/models"
	"github.com/galaxy-empire-team/event-manager/internal/service/mission/mocks"
)

func TestService_isMoonCreated(t *testing.T) {
	tests := []struct {
		name      string
		debris    models.Resources
		setupRand func(randGen *mocks.RandGenerator)
		want      bool
	}{
		{
			name: "below the threshold",
			debris: models.Resources{
				Metal:   0,
				Crystal: 1_000_000,
			},
			setupRand: func(randGen *mocks.RandGenerator) {
				randGen.EXPECT().Intn(100).Return(0).Once()
			},
			want: true,
		},
		{
			name: "over the threshold",
			debris: models.Resources{
				Metal:   5_000_000,
				Crystal: 1_000_000,
			},
			setupRand: func(randGen *mocks.RandGenerator) {
				randGen.EXPECT().Intn(100).Return(21).Once()
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			randGen := mocks.NewRandGenerator(t)
			tt.setupRand(randGen)

			svc := &Service{
				randGenerator: randGen,
				logger:        zap.NewNop(),
			}

			got := svc.isMoonCreated(tt.debris)

			assert.Equal(t, tt.want, got)
		})
	}
}
