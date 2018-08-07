package job_controller

import (
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
    workItems, err := client.SMembers("stream-list").Result()
    if err != nil {
        return nil, err
    }

    return workItems, nil
}

func RenameReadKey(client redis.Client) error {
    val, err := client.Rename("stream-list", "work").Result()
    _ = val
    if err != nil {
        return err
    }

    return nil
}
