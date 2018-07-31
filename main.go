package main

import (
    "time"
    "./src/job_controller"
)

func doEvery(d time.Duration, f func(time.Time)) {
    for x := range time.Tick(d) {
        f(x)
    }
}

func main() {
    doEvery(15000*time.Millisecond, job_controller.StartNewJobSet)
}
