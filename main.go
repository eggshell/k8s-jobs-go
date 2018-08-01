package main

import (
    "fmt"
    "time"
    jc "./src/job_controller"
)

func main() {
    listInterval := time.NewTicker(20 * time.Second).C
    bestInterval := time.NewTicker(30 * time.Second).C

    for {
        select {
        case <- listInterval:
            fmt.Println("Getting list of streams")
            streams := jc.GetLiveStreams()
            jc.UpdateStreamsList(streams)
        case <- bestInterval:
            fmt.Println("Finding best stream")
            //jc.StartNewJobSet()
        }
    }
}
