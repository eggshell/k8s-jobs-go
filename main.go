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
            fmt.Println("Checking for work to do")
            workItems := jc.CheckWorkQueue()
            if workItems != nil {
                jc.StartNewJobSet(workItems)
            } else {
                fmt.Println("No work to do. Waiting for next interval.")
            }
        }
    }
}
