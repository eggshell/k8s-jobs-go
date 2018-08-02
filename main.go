package main

import (
    "fmt"
    "time"
    jc "./src/job_controller"
)

func main() {
    interval := time.NewTicker(3 * time.Second).C

    for {
        select {
        case <- interval:
            fmt.Println("Checking stream list in Redis")
            streams := jc.CheckWorkQueue()
            if streams != "" {
                jc.StartNewJobSet()
            } else {
                fmt.Println("No work to do. Waiting for next interval.")
            }
        }
    }
}
