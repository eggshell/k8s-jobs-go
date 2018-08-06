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
            workItems := jc.CheckWorkQueue(redisClient)
            if workItems != nil {
                jc.CreateJob(workItems)
            } else {
                fmt.Println("No work to do. Waiting for next interval.")
            }
        }
    }
}
