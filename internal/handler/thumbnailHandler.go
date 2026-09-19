package handler

import (
	"context"
	"encoding/json"
	"fmt"
	thumbnail "video_processing_pipeline/internal/Thumbnail"
	"video_processing_pipeline/internal/queue"

	"github.com/hibiken/asynq"
)

type ThumbnailHandler struct {
	next queue.AsyncPublishManager
	thmbnl thumbnail.ThubnailManager
}
func NewTHubnailHandler(tsk queue.AsyncPublishManager, thmb thumbnail.ThubnailManager)*ThumbnailHandler{
	return &ThumbnailHandler{
		next: tsk,
		thmbnl: thmb,
	}
}
func(t *ThumbnailHandler)HandleThumbnail(ctx context.Context, task *asynq.Task)error{
	fmt.Println("Enter into thubnail Handler")
	var jbq queue.JobQueue
	err:=json.Unmarshal(task.Payload(),&jbq)
	if err!=nil{
		return err
	}
    err1:=t.thmbnl.GenerateThumbnail(ctx,jbq)
	if err1!=nil{
		return err1
	}
	err2:=t.next.Publish(ctx,queue.TypeNotify,&jbq)
	if err2!=nil{
		return err2
	}
	return nil
}