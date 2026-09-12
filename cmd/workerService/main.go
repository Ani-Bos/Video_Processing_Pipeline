package workerservice

import (
	"log"
    "os"
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
	// mux.HandleFunc()
	// mux.HandleFunc()
	log.Fatal(srv.Run(mux))
}