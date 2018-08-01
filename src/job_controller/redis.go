package job_controller

import (
    "encoding/json"
    "os"
    "github.com/go-redis/redis"
)

func GetRedisClient() redis.Client {
    redisAddr := os.Getenv("redis_addr")
    redisPort := os.Getenv("redis_port")
    redisPass := os.Getenv("redis_pass")
    client := redis.NewClient(&redis.Options{
        Addr:       redisAddr + ":" + redisPort,
        Password:   redisPass,
        DB:         0,
    })

    return *client
}

// puts list of streams in redis
func UpdateStreamsList(streams []Stream) {
    client := GetRedisClient()

    streamsJSON, errJSON := json.Marshal(streams)
    if errJSON != nil {
        panic(errJSON)
    }

    err := client.SAdd("streams", streamsJSON, 0).Err()
    if err != nil {
        panic(err)
    }
}

// TODO
// get list of streams stored in redis
//func CheckStreamsList() {
//    client := GetRedisClient()
//}
