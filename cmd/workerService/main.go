package main

import (
	"log"
	"os"
	initializer "video_processing_pipeline/internal/Initializer"
	thumbnail "video_processing_pipeline/internal/Thumbnail"
	transcoder "video_processing_pipeline/internal/Transcoder"
	"video_processing_pipeline/internal/handler"
	"video_processing_pipeline/internal/notify"
	"video_processing_pipeline/internal/queue"
	"video_processing_pipeline/internal/repository"

	"github.com/hibiken/asynq"
)

func main() {
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

	mux:=asynq.NewServeMux()
	handler1:=handler.NewTranscodeHandler(newffmpeg,publisher)
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
    mailDetails := notify.NewNotifyManager(smtpCfg)
	handler3:=handler.NewNotificationHandler(*mailDetails,jobs_repo)
	mux.HandleFunc(queue.TypeNotify,handler3.NotificationHandler)
	log.Println("starting worker service to do processing of thumnaila nd transcode")
	log.Fatal(srv.Run(mux))
}