package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"video_processing_pipeline/internal/Transcoder"
	"video_processing_pipeline/internal/queue"
	"github.com/hibiken/asynq"
)

type TranscodeHandler struct {
	 tc transcoder.TranscoderManager
	nextTask queue.AsyncPublishManager
}

func NewTranscodeHandler(trnscdr *transcoder.TranscoderManager, tsk *queue.AsyncPublishManager) *TranscodeHandler {
	return &TranscodeHandler{
		tc: *trnscdr,
		nextTask: *tsk,
	}
}

func(th *TranscodeHandler)HandleTranscoding(ctx context.Context, t *asynq.Task)error{
	fmt.Println("Entering into transcioding handler")
	var jbq queue.JobQueue
	err1:=json.Unmarshal(t.Payload(), &jbq)
	if err1!=nil{
		return err1
	}
    err:=th.tc.Transcode(ctx,&jbq)
	if err!=nil{
		return err
	}
	err=th.nextTask.Publish(ctx,queue.TypeThumbnail,&jbq)
	if err!=nil{
		return err
	}
	return nil
}