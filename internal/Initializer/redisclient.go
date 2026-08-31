package initializer

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(){
   fmt.Println("entering into redis initilaizition")
   redis_host:=os.Getenv("REDIS_URL")
   if redis_host==""{
	 redis_host="localhost"
   }
   psswd := os.Getenv("REDIS_PWD")
   rdb := redis.NewClient(&redis.Options{
	    Addr:     redis_host + ":6379",
        Password: psswd,
        DB:       0,
        Protocol: 2,
   })
    _, err := rdb.Ping(context.Background()).Result()
    if err != nil {
       fmt.Println("Failed to connect to Redis")
        return
    }
	fmt.Println("connected to redis successfully")
}