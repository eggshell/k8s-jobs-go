package main

import (
    "fmt"
    "time"
    jc "./src/job_controller"
)

func main() {
    var lock bool    = false
    redisClient     := jc.GetRedisClient()
    kubeClient, err := jc.KubeClientInCluster()
    if err != nil {
        panic(err.Error())
    }
    workInterval := time.NewTicker(1 * time.Second).C
    lockInterval := time.NewTicker(333 * time.Millisecond).C

    for {
        select {
        case <- workInterval:
            fmt.Println("Checking for work to do")
            fmt.Println("lock status: ", lock)

            workItems, err := jc.CheckWorkQueue(redisClient)
            if err != nil {
                fmt.Println(err.Error())
                break
            }

            if workItems != nil && len(workItems) != 0 && lock == false {
                lock = true
                job, err := jc.CreateJob(kubeClient, redisClient, workItems)
                if err != nil {
                    fmt.Println(err.Error())
                }

                fmt.Println("Created job %q.", job)
            } else {
                fmt.Println("No work to do. Waiting for next interval.")
            }
        case <- lockInterval:
            fmt.Println("Checking to see if it's time to unlock the queue.")
            fmt.Println("lock status: ", lock)

            workExists, err := redisClient.Exists("work").Result()
            if err != nil {
                fmt.Println(err.Error())
            }

            jobs, err := jc.ListJobs(kubeClient, "default")
            if err != nil {
                fmt.Println(err.Error())
            }

            if len(jobs.Items) != 0 {
                fmt.Println("Active: ", jobs.Items[0].Status.Active)
                fmt.Println("Succeeded: ", jobs.Items[0].Status.Succeeded)
                fmt.Println("Failed: ", jobs.Items[0].Status.Failed)
            }

            if workExists == 0 && len(jobs.Items) != 0 && jc.IsJobFinished(jobs.Items[0]) == true {
                fmt.Println("Deleting finished job")
                err := jc.DeleteJob(kubeClient, jobs.Items[0])
                if err != nil {
                    fmt.Println(err.Error())
                }

                fmt.Println("Unlocking the queue")
                lock = false
            }
        }
    }
}
