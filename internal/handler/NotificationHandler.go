package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"video_processing_pipeline/internal/notify"
	"video_processing_pipeline/internal/queue"
	"video_processing_pipeline/internal/repository"

	"github.com/hibiken/asynq"
)

type NotificationHandler struct {
	Ntfy notify.NotifyManager
	Repo repository.JobRepo
}

func NewNotificationHandler( Ntfyi notify.NotifyManager, Repos repository.JobRepo)*NotificationHandler{
	return &NotificationHandler{
	  Ntfy: Ntfyi,
	  Repo: Repos,
	}
}

func(n *NotificationHandler)NotificationHandler(ctx context.Context, task *asynq.Task)error{
  fmt.Println("Enter into simple notification handler")
  var jbq queue.JobQueue
	err:=json.Unmarshal(task.Payload(),&jbq)
	if err!=nil{
		return err
	}
	jobs,err:=n.Repo.FindByVideoID(ctx,jbq.VideoId)
	if err!=nil{
		return err
	}
	if jobs.NotifiedAt!=nil{
		return nil
	}
    err=n.Ntfy.Notify(ctx,jobs)
	if err!=nil{
		return err
	}
	now:=time.Now().UTC()
	jobs.NotifiedAt=&now
    return n.Repo.Update(ctx,jobs)
}