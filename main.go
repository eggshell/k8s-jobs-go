package main

import (
    "fmt"
    "time"
    jc "./src/job_controller"
)

func main() {
    redisClient := jc.GetRedisClient()
    interval := time.NewTicker(3 * time.Second).C

    for {
        select {
        case <- interval:
            fmt.Println("Checking for work to do")

            workExists, err := redisClient.Exists("work").Result()
            if err != nil {
                fmt.Println(err.Error())
                break
            }

            workItems, err := jc.CheckWorkQueue(redisClient)
            if err != nil {
                fmt.Println(err.Error())
                break
            }

            if workItems != nil && len(workItems) != 0 && workExists == 0 {
                fmt.Println("workItems: ", len(workItems))
                jc.CreateJob(redisClient, workItems)
            } else {
                fmt.Println("No work to do. Waiting for next interval.")
            }
        }
    }
}
