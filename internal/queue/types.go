package queue

import "context"

const (
	TypeTranscode = "transcode"
	TypeThumbnail = "thumbnail"
	TypeNotify    = "notify"
)

type JobQueue struct {
	VideoId   string
	RawPath  string
	FileName string
	OutPutDir   string
}

type AsyncPublishManager interface {
	Publish(ctx context.Context, Topic string, job *JobQueue) error
}
