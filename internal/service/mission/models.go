package mission

import "github.com/google/uuid"

type userIDPair struct {
	Attacker uuid.UUID
	Defender uuid.UUID
}
