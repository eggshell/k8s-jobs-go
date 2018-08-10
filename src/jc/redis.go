package jc

import (
    "os"
    "github.com/go-redis/redis"
)

// GetRedisClient gets a redisClient object defined in types.go and returns it
func GetRedisClient() RedisClient {
    client := redis.NewClient(&redis.Options{
        Addr:       "redis-master:6379",
        Password:   os.Getenv("REDIS_PASSWORD"),
        DB:         0,
    })

    return *client
}

// CheckWorkQueue uses the SMEMBERS redis method to return a string array of
// all members of a given set.
func CheckWorkQueue(client RedisClient) ([]string, error) {
    workItems, err := client.SMembers("stream-list").Result()
    if err != nil {
        return nil, err
    }

    return workItems, nil
}

// RenameReadKey uses the RENAME redis method to rename a key in redis.
func RenameReadKey(client RedisClient) error {
    val, err := client.Rename("stream-list", "work").Result()
    _ = val
    if err != nil {
        return err
    }

    return nil
}
