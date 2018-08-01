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

// puts list of streams in redis
func UpdateStreamsList(streams []Stream) {
    client := GetRedisClient()

    streamsJSON, errJSON := json.Marshal(streams)
    if errJSON != nil {
        fmt.Println(errJSON)
    }

    err := client.SAdd("streams", streamsJSON, 0).Err()
    if err != nil {
        fmt.Println(err)
    }
}

// TODO
// get list of streams stored in redis
//func CheckStreamsList() {
//    client := GetRedisClient()
//}
