package thumbnail

import (
	"context"
	"video_processing_pipeline/internal/queue"
)

type ThubnailManager interface {
	GenerateThumbnail(ctx context.Context, job *queue.JobQueue)error
}