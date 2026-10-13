package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	initializer "video_processing_pipeline/internal/Initializer"
	"video_processing_pipeline/internal/handler"
	"video_processing_pipeline/internal/model"
	"video_processing_pipeline/internal/queue"
	"video_processing_pipeline/internal/repository"
	"video_processing_pipeline/internal/service"
	"video_processing_pipeline/internal/uploader/chunkersse"

	"github.com/hibiken/asynq"
)

func main() {
	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	db := initializer.DbInitializer{}
	initializer.ConnectDB(&db)
	db.DB.AutoMigrate(&model.Chunk{}, &model.Chunk_Session{},&model.Jobs_Database{})
	repo := &repository.ChunkRepo{DB: db.DB}
	jobs_repo:=&repository.JobRepo{DB: db.DB}
	jobs_db := &service.InterfaceInjectRepoJob{Repo: jobs_repo}
	redis_host:=os.Getenv("REDIS_ADDR")
	redis_pswd:=os.Getenv("REDIS_PWD")
	redisoption := asynq.RedisClientOpt{
		Addr: redis_host,
		Password: redis_pswd,
    }
	client:=asynq.NewClient(redisoption)
	defer client.Close()
	// dir := os.Getenv("UPLOAD_DIR")
	// if dir == "" {
	// 	dir = "/data/uploads"
	// }
	bucketname:="videoProcessing"
    minioclient,err:=initializer.ConnectMiniIOAndCreateBucket(ctx,bucketname)
	if err!=nil{
		log.Fatal("Statup failed due to miniiIo",err)
		return
	}
    publisher:=queue.NewAsyncPublisher(client)
	// manager:=chunkersse.NewDBManager(repo, dir)
	manager:=chunkersse.NewMiniIOManager(*repo,minioclient,bucketname)
	handler1:=handler.NewHandlerStruct(manager,*jobs_db,*repo,*publisher)
	// manager:=chunkersse.NewChunkedUploadManager("uploads")
	// handler1:=handler.NewHandlerStruct(manager)
	http.HandleFunc("/upload/init",handler1.HandleStartUpload)
	http.HandleFunc("/upload/chunk",handler1.HandleUploadChunks)
	http.HandleFunc("/upload/complete",handler1.HandleCompleteUpload)
	http.HandleFunc("/upload/status",handler1.HandleGetStatusUpload)
	http.HandleFunc("/",handler.UploadHandler)
	http.HandleFunc("/upload",handler.UploadStreamingHandler)
	// http.HandleFunc("/handleUpload",uploader.UploadChunkedHandler)
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}