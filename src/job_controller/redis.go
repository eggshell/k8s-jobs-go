package job_controller

import (
    "fmt"
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

func CheckWorkQueue() string {
    client := GetRedisClient()
    val, err := client.Get("stream-list").Result()
    if err != nil {
        fmt.Println(err.Error())
        return ""
    } else if val == "" {
        fmt.Println("Empty stream list")
        return ""
    }

    return val
}

func RenameWorkKey() int {
    client := GetRedisClient()
    err := client.Rename("stream-list", "stream-list-old")
    if err != nil {
        fmt.Println(err)
        return 1
    }

    return 0
}
