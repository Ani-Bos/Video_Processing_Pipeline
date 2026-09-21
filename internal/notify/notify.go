package notify

import (
	"context"
	"fmt"
	"video_processing_pipeline/internal/model"
)

type NotifyManager struct {
	Ntfy NotifyInterface
}

func NewNotifyManager(ntfy NotifyInterface) *NotifyManager {
	return &NotifyManager{Ntfy: ntfy}
}

func(N *NotifyManager)Notify(ctx context.Context, job *model.Jobs_Database) error {
  fmt.Println("Enter into notification service to notify users")
  return nil
}