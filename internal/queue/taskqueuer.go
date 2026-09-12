package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
)

type AsyncPublishersrvc struct {
	Client *asynq.Client
}

func NewAsyncPublisher(client *asynq.Client) *AsyncPublishersrvc {
    return &AsyncPublishersrvc{
		Client: client,
	}
}

func (p *AsyncPublishersrvc)Publish(ctx context.Context,Topic string,job *JobQueue)(error){
	fmt.Println("Entering into publish events")
	payload,err:=json.Marshal(job)
	if err!=nil{
		return err
	}
    task:=asynq.NewTask(Topic,payload,asynq.MaxRetry(3))
	_,err=p.Client.Enqueue(task)
	return err
}