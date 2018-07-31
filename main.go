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
            user_ids := jc.UpdateStreamsList()
            for k := range user_ids {
                fmt.Println(user_ids[k])
            }
        case <- bestInterval:
            fmt.Println("Finding best stream")
      }
    }
}
