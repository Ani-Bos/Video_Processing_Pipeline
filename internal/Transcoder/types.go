package transcoder

import (
	"context"
	"video_processing_pipeline/internal/queue"
)

type TranscoderManager interface {
	Transcode(ctx context.Context,job *queue.JobQueue)(error)
}
