package notify

import (
	"context"
	"video_processing_pipeline/internal/model"
)

type NotifyInterface interface {
	Notify(ctx context.Context, job *model.Jobs_Database)error
}