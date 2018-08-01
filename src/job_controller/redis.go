package job_controller

import (
    "fmt"
    "encoding/json"
    "os"
    "github.com/go-redis/redis"
)

func GetRedisClient() redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:       os.Getenv("redis_addr") + ":" + os.Getenv("redis_port"),
        Password:   os.Getenv("redis_pass"),
        DB:         0,
    })

    return *client
}

// TODO
// get list of streams stored in redis
//func CheckWorkQueue() {
//    client := GetRedisClient()
//}
