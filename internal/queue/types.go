package queue

import "context"

type JobQueue struct {
	VideId   string
	RawPath  string
	FileName string
}

type PublishManager interface {
	Publish(ctx context.Context, Topic string, jb *JobQueue)error
}

type SubscribeManager interface {
	Subscribe(ctx context.Context,stream,group,subscriber string,)
}