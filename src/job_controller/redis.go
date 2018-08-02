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
    workKey := os.Getenv("redis_work_key")

    val, err := client.Get(workKey).Result()
    if err != nil {
        fmt.Println(err.Error())
        return ""
    } else if val == "" {
        fmt.Println("Empty work queue")
        return ""
    }

    return val
}

func RenameWorkKey() int {
    client := GetRedisClient()
    err := client.Rename(workKey, workKey + "-old")
    if err != nil {
        fmt.Println(err)
        return 1
    }

    return 0
}
