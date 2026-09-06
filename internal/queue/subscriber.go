package queue

import (
	"github.com/redis/go-redis/v9"
)

type RedisSubscriber struct{
	Client *redis.Client
	Topic string
}

func NewSubscriber(rdb *redis.Client, TopicConsume string)*RedisSubscriber{
	return &RedisSubscriber{
		Client: rdb,
		Topic:TopicConsume,
	}
}