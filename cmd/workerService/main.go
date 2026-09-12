package workerservice

import (
	"log"
	"os"
	handler1 "video_processing_pipeline/internal/Handler"
	"github.com/hibiken/asynq"
)

func main() {
	redis_host:=os.Getenv("REDIS_URL")
	redisoption := asynq.RedisClientOpt{Addr: redis_host}
	srv := asynq.NewServer(redisoption, asynq.Config{
        Concurrency: 2,
        Queues: map[string]int{
            "critical": 6,
            "default":  3,
        },
    })
	mux:=asynq.NewServeMux()
	// handler2:=handler1.NewTranscodeHandler()
	// mux.HandleFunc("/transcode/v1",handler2.HandleTranscoding)
	// mux.HandleFunc()
	log.Fatal(srv.Run(mux))
}