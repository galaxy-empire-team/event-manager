package notification

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/event-manager/internal/models"
)

func (r *NotificationStorage) SaveNotificationEvents(ctx context.Context, notificationEvents []models.NotificationEvent) error {
	const createNotificationQuery = `
		INSERT INTO session_beta.user_notifications (user_id, notification_id, data)
		VALUES ($1, $2, $3);
	`

	var batch pgx.Batch
	for _, event := range notificationEvents {
		batch.Queue(createNotificationQuery, event.UserID, event.NotificationID, event.Data)
	}

	br := r.DB.SendBatch(ctx, &batch)
	defer br.Close()

	for range notificationEvents {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("br.Exec(): %w", err)
		}
	}

	return nil
}
