package queue

import "context"

const (
	TypeTranscode = "transcode"
	TypeThumbnail = "thumbnail"
	TypeNotify    = "notify"
)
const (
	StatusUploaded   = "uploaded"    
	StatusInProgress = "in_progress" 
	StatusFailed     = "failed"
	StatusSuccessful = "successful"
)
const (
	StageNone=""
	StageTranscode="transcode"
	StageThumbnail="thumbnail"
	StageNotify="notify"
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
