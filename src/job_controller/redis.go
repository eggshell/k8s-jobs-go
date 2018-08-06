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

func CheckWorkQueue(client redis.Client) ([]string, error) {
    //workKey := os.Getenv("stream-list")

    workItems, err := client.SMembers("stream-list").Result()
    if err != nil {
        return nil, err
    }

    return workItems, nil
}

func RenameWorkKey(client redis.Client) int {
    //workKey := os.Getenv("stream-list")
    err := client.Rename("stream-list", "stream-list-old")
    if err != nil {
        fmt.Println(err)
        return 1
    }

    return 0
}
