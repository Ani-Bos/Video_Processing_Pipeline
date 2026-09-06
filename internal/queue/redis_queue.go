package queue

import "github.com/redis/go-redis/v9"

type RedisQueue struct {
	Client *redis.Client
}

func(r *RedisQueue)publish()(error){
return redis.Nil
}

func(r *RedisQueue)Subscribe()(error){
 return redis.Nil
}