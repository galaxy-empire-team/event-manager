package mission

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func TestService_fillFleetCargo(t *testing.T) {
	tests := []struct {
		name            string
		planetResources models.Resources
		fleetCargo      uint64
		gained          models.Resources
		left            models.Resources
	}{
		// Основной сценарий
		{
			name: "standard case with surplus resources",
			planetResources: models.Resources{
				Metal:   1000,
				Crystal: 500,
				Gas:     250,
			},
			fleetCargo: 1000,
			gained: models.Resources{
				Metal:   450,
				Crystal: 350,
				Gas:     200,
			},
			left: models.Resources{
				Metal:   550,
				Crystal: 150,
				Gas:     50,
			},
		},
		{
			name: "empty resources on planet",
			planetResources: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     0,
			},
			fleetCargo: 1000,
			gained: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     0,
			},
			left: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     0,
			},
		},
		{
			name: "zero cargo capacity",
			planetResources: models.Resources{
				Metal:   1000,
				Crystal: 500,
				Gas:     250,
			},
			fleetCargo: 0,
			gained: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     0,
			},
			left: models.Resources{
				Metal:   1000,
				Crystal: 500,
				Gas:     250,
			},
		},
		{
			name: "all resources fit in cargo",
			planetResources: models.Resources{
				Metal:   100,
				Crystal: 100,
				Gas:     100,
			},
			fleetCargo: 1000,
			gained: models.Resources{
				Metal:   100,
				Crystal: 100,
				Gas:     100,
			},
			left: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     0,
			},
		},
		{
			name: "only metal available",
			planetResources: models.Resources{
				Metal:   2000,
				Crystal: 0,
				Gas:     0,
			},
			fleetCargo: 1000,
			gained: models.Resources{
				Metal:   1000,
				Crystal: 0,
				Gas:     0,
			},
			left: models.Resources{
				Metal:   1000,
				Crystal: 0,
				Gas:     0,
			},
		},
		{
			name: "only crystal available",
			planetResources: models.Resources{
				Metal:   0,
				Crystal: 2000,
				Gas:     0,
			},
			fleetCargo: 1000,
			gained: models.Resources{
				Metal:   0,
				Crystal: 1000,
				Gas:     0,
			},
			left: models.Resources{
				Metal:   0,
				Crystal: 1000,
				Gas:     0,
			},
		},
		{
			name: "only gas available",
			planetResources: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     2000,
			},
			fleetCargo: 1000,
			gained: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     1000,
			},
			left: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     1000,
			},
		},
		{
			name: "fill remaining cargo with excess metal",
			planetResources: models.Resources{
				Metal:   1500,
				Crystal: 300,
				Gas:     100,
			},
			fleetCargo: 1000,
			gained: models.Resources{
				Metal:   600,
				Crystal: 300,
				Gas:     100,
			},
			left: models.Resources{
				Metal:   900,
				Crystal: 0,
				Gas:     0,
			},
		},
		{
			name: "fill remaining cargo with multiple resources",
			planetResources: models.Resources{
				Metal:   400,
				Crystal: 400,
				Gas:     2000,
			},
			fleetCargo: 1000,
			gained: models.Resources{
				Metal:   400,
				Crystal: 400,
				Gas:     200,
			},
			left: models.Resources{
				Metal:   0,
				Crystal: 0,
				Gas:     1800,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				logger: zap.NewNop(),
			}

			fillResult := svc.fillFleetCargo(tt.planetResources, tt.fleetCargo)

			assert.Equal(t, tt.gained, fillResult.gained, "gained resources mismatch")
			assert.Equal(t, tt.left, fillResult.leftOnPlanet, "left resources mismatch")

			totalGained := fillResult.gained.Metal + fillResult.gained.Crystal + fillResult.gained.Gas
			totalLeft := fillResult.leftOnPlanet.Metal + fillResult.leftOnPlanet.Crystal + fillResult.leftOnPlanet.Gas
			totalOriginal := tt.planetResources.Metal + tt.planetResources.Crystal + tt.planetResources.Gas
			assert.Equal(t, totalOriginal, totalGained+totalLeft, "total resources not conserved")

			assert.LessOrEqual(t, totalGained, tt.fleetCargo, "gained resources exceed cargo capacity")
		})
	}
}
