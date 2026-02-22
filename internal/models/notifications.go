package models

import (
	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type NotificationEvent struct {
	UserID         uuid.UUID
	NotificationID consts.NotificationID
	Data           []byte
}
