package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	initializer "video_processing_pipeline/internal/Initializer"
	thumbnail "video_processing_pipeline/internal/Thumbnail"
	transcoder "video_processing_pipeline/internal/Transcoder"
	"video_processing_pipeline/internal/handler"
	mediastore "video_processing_pipeline/internal/mediaStore"
	"video_processing_pipeline/internal/notify"
	"video_processing_pipeline/internal/queue"
	"video_processing_pipeline/internal/repository"

	"github.com/hibiken/asynq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	db := initializer.DbInitializer{}
	initializer.ConnectDB(&db)
	jobs_repo:=repository.JobRepo{DB:db.DB}
	redis_host:=os.Getenv("REDIS_ADDR")
	redis_pswd:=os.Getenv("REDIS_PWD")
	redisoption := asynq.RedisClientOpt{
		Addr: redis_host,
		Password: redis_pswd,
    }
	srv := asynq.NewServer(redisoption, asynq.Config{
        Concurrency: 2,
        Queues: map[string]int{
            "critical": 6,
            "default":  3,
        },
    })
	client:=asynq.NewClient(redisoption)
	defer client.Close()

	publisher:=queue.NewAsyncPublisher(client)
	//we need to install it in container app image 
	newffmpeg:=transcoder.NewFFMPEG(
		os.Getenv("FFMPEG_BIN"),
		os.Getenv("FFMPEG_WORKDIR"),
	)
	newffmpeg1:=thumbnail.NewFFMPEG(
		os.Getenv("FFMPEG_BIN"),
		os.Getenv("FFMPEG_WORKDIR"),
	)
	newffmpeg2:=transcoder.NewFFMPEG(
		os.Getenv("FFMPEG_BIN"),
		os.Getenv("FFMPEG_WORKDIR"),
	)
	bucketname:="videoprocessing"
    minioclient,err:=initializer.ConnectMiniIOAndCreateBucket(ctx,bucketname)
	defer cancel()
	if err!=nil{
		log.Fatal("Statup failed due to miniiIo",err)
		return
	}

	mux:=asynq.NewServeMux()
	MediaStrMgr:=mediastore.NewMediaStore(nil,minioclient,bucketname)
	handler1:=handler.NewTranscodeHandler(*newffmpeg2,newffmpeg,publisher,MediaStrMgr)
	// handler2:=handler1.NewTranscodeHandler()
	mux.HandleFunc(queue.TypeTranscode,handler1.HandleTranscoding)
	handler2:=handler.NewTHubnailHandler(publisher,newffmpeg1)
    mux.HandleFunc(queue.TypeThumbnail,handler2.HandleThumbnail)
	smtpCfg := notify.ConfigSMTP{
	Host:     os.Getenv("SMTP_HOST"),    
	Port:     587,
	Username: os.Getenv("SMTP_USER"),     
	Password: os.Getenv("SMTP_PASS"),   
	From:     os.Getenv("SMTP_FROM"),    
	To:       []string{os.Getenv("SMTP_TO")},
}
    access_key:=os.Getenv("S3_ACCESS_KEY")
	secret_key:=os.Getenv("S3_SECRET_KEY")
	notifyhst:=os.Getenv("S3_PUBLIC_ENDPOINT")
	fmt.Println("notifyhost is",notifyhst)
    notifyclient,err1:=minio.New(notifyhst,&minio.Options{
		Creds: credentials.NewStaticV4(access_key,secret_key,""),
		Secure: false,
		Region: "us-east-1",
	})
	if err1!=nil{
		log.Fatal("Statup failed due to miniiIo notification client",err1)
		return
	}
	Notifystrmgr:=mediastore.NewMediaStore(nil,notifyclient,bucketname)
    mailDetails := notify.NewNotifyManager(smtpCfg,*Notifystrmgr)
	handler3:=handler.NewNotificationHandler(*mailDetails,jobs_repo)
	mux.HandleFunc(queue.TypeNotify,handler3.NotificationHandler)
	log.Println("starting worker service to do processing of thumnaila nd transcode")
	log.Fatal(srv.Run(mux))
}