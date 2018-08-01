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
            // some conditional.
            // if work to do, fire up jobs
            // check if pods are done doing work and rename key
            // else nothing
        }
    }
}
