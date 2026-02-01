package mission

import (
	"github.com/google/uuid"
)

type planetToColonize struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Coordinates coordinates
	HasMoon     bool
	IsCapitol   bool
}

type coordinates struct {
	X uint8
	Y uint8
	Z uint8
}
