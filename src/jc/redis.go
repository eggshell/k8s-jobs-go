package jc

import (
    "os"
    "github.com/go-redis/redis"
)

func GetRedisClient() RedisClient {
    client := redis.NewClient(&redis.Options{
        Addr:       "redis-master:6379",
        Password:   os.Getenv("REDIS_PASSWORD"),
        DB:         0,
    })

    return *client
}

func CheckWorkQueue(client RedisClient) ([]string, error) {
    workItems, err := client.SMembers("stream-list").Result()
    if err != nil {
        return nil, err
    }

    return workItems, nil
}

func RenameReadKey(client RedisClient) error {
    val, err := client.Rename("stream-list", "work").Result()
    _ = val
    if err != nil {
        return err
    }

    return nil
}
