package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	Client *redis.Client
}

func NewPublisher(rdb *redis.Client)*RedisPublisher{
	return &RedisPublisher{
     Client: rdb,
	}
}

func(r *RedisPublisher)publish(ctx context.Context, topic string, msg string)(error){   
	fmt.Println("Publishing new events")
    ctx,cancel:=context.WithTimeout(ctx,5*time.Second)
	defer cancel()
	err:=r.Client.Publish(ctx,topic,msg).Err()
	if err!=nil{
		return err
	}
	return nil
}