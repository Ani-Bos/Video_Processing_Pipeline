package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	thumbnail "video_processing_pipeline/internal/Thumbnail"
	transcoder "video_processing_pipeline/internal/Transcoder"
	mediastore "video_processing_pipeline/internal/mediaStore"
	"video_processing_pipeline/internal/queue"

	"github.com/hibiken/asynq"
)

type ThumbnailHandler struct {
	next queue.AsyncPublishManager
	thmbnl thumbnail.ThubnailManager
	stre *mediastore.MediaStore
	trnscdr       transcoder.FFMPEG
}
func NewTHubnailHandler(tsk queue.AsyncPublishManager, thmb thumbnail.ThubnailManager, store *mediastore.MediaStore, transcdr transcoder.FFMPEG)*ThumbnailHandler{
	return &ThumbnailHandler{
		next: tsk,
		thmbnl: thmb,
		stre: store,
		trnscdr: transcdr,
	}
}
func(t *ThumbnailHandler)HandleThumbnail(ctx context.Context, task *asynq.Task)error{
	fmt.Println("Enter into thubnail Handler")
	var jbq queue.JobQueue
	err:=json.Unmarshal(task.Payload(),&jbq)
	if err!=nil{
		return err
	}
	prefix:=jbq.OutPutDir
	localdir:=filepath.Join(t.trnscdr.WORKING_DIR,jbq.VideoId)
	fmt.Println("prefix under thumbnail is :",prefix)
	fmt.Println("localdir under thumbnail is :",localdir)
	err=t.stre.DownloadHLSObject(ctx,localdir,prefix)
	if err!=nil{
		return err
	}
	jbq.OutPutDir=localdir
    err1:=t.thmbnl.GenerateThumbnail(ctx,jbq)
	if err1!=nil{
		return err1
	}
	imagepath:=filepath.Join(localdir,"ani.jpg")
	fmt.Println("image path final thubnai;l image in handler is ",imagepath)
	err=t.stre.PutFileObject(ctx,prefix+"/ani.jpg",imagepath)
	if err!=nil{
		return err
	}
	_=os.RemoveAll(localdir)
	fmt.Println("thumbnail image uploaded into bucket")
	jbq.OutPutDir=prefix
	err2:=t.next.Publish(ctx,queue.TypeNotify,&jbq)
	if err2!=nil{
		return err2
	}
	return nil
}