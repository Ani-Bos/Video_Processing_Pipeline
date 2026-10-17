package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"video_processing_pipeline/internal/Transcoder"
	mediastore "video_processing_pipeline/internal/mediaStore"
	"video_processing_pipeline/internal/queue"

	"github.com/hibiken/asynq"
)

type TranscodeHandler struct {
	 trnscdr       transcoder.FFMPEG
	 tc transcoder.TranscoderManager
	nextTask queue.AsyncPublishManager
	stre  *mediastore.MediaStore
}

func NewTranscodeHandler(tcdr transcoder.FFMPEG,trnscdr transcoder.TranscoderManager, tsk queue.AsyncPublishManager, Mediastr *mediastore.MediaStore) *TranscodeHandler {
	return &TranscodeHandler{
		trnscdr: tcdr,
		tc: trnscdr,
		nextTask: tsk,
		stre: Mediastr,
	}
}

func(th *TranscodeHandler)HandleTranscoding(ctx context.Context, t *asynq.Task)error{
	fmt.Println("Entering into transcioding handler")
	var jbq queue.JobQueue
	err1:=json.Unmarshal(t.Payload(), &jbq)
	if err1!=nil{
		return err1
	}
	localfile:= filepath.Join(th.trnscdr.WORKING_DIR,jbq.VideoId,"source.mp4")
	fmt.Println("job raw path is as: ", jbq.RawPath)
	fmt.Println("local path of jobs db for transcoding and downloading object",localfile)
	//here basiucally downloading the minio object hat is uploaded as a complete in handlecompleteupload
	//store.DownloadOnject(ctx, jbq.OutPutDir, localDir);
	err:=th.stre.DownloadObject(ctx,jbq.RawPath,localfile)
	if err!=nil{
		return err
	}
	jbq.RawPath=localfile
    err=th.tc.Transcode(ctx,&jbq)
	if err!=nil{
		return err
	}
	//pushing object into bucket to prevent concurrency issue
	_,err=th.stre.PutObjects(ctx,jbq.OutPutDir)
	if err!=nil{
		return err
	}
	err=th.nextTask.Publish(ctx,queue.TypeThumbnail,&jbq)
	if err!=nil{
		return err
	}
	return nil
}