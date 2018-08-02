package job_controller

import (
    "fmt"
    "os"
    "github.com/go-redis/redis"
)

func GetRedisClient() redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:       "redis-master:6379",
        Password:   os.Getenv("REDIS_PASSWORD"),
        DB:         0,
    })

    return *client
}

func CheckWorkQueue() []string {
    client    := GetRedisClient()
    workKey   := os.Getenv("REDIS_WORK_KEY")

    workItems, err := client.ZRange(workKey, 0, 20).Result()
    if err != nil {
        return nil
    }

    return workItems
}

func RenameWorkKey() int {
    client := GetRedisClient()
    workKey := os.Getenv("REDIS_WORK_KEY")
    err := client.Rename(workKey, workKey + "-old")
    if err != nil {
        fmt.Println(err)
        return 1
    }

    return 0
}
