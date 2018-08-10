package main

import (
	"./src/jc"
	"fmt"
	"time"
)

func attemptJobCreation(k *jc.KubeClient, lock bool, r jc.RedisClient) (bool, error) {
	workItems, err := jc.CheckWorkQueue(r)
	if err != nil {
		return false, err
	}

	if workItems != nil && len(workItems) != 0 && lock == false {
		lock = true
		job, err := jc.CreateJob(k, r, workItems)
		if err != nil {
			return false, err
		}
		fmt.Println("Created job: ", job)
	} else {
		fmt.Println("No work to do. Waiting for next interval.")
	}
	return lock, nil
}

func attemptQueueUnlock(k *jc.KubeClient, lock bool, r jc.RedisClient) (bool, error) {
	workExists, err := r.Exists("work").Result()
	if err != nil {
		return lock, err
	}

	jobs, err := jc.ListJobs(k, "default")
	if err != nil {
		return lock, err
	}

	if len(jobs.Items) != 0 {
		fmt.Println("Active: ", jobs.Items[0].Status.Active)
		fmt.Println("Succeeded: ", jobs.Items[0].Status.Succeeded)
		fmt.Println("Failed: ", jobs.Items[0].Status.Failed)
	}

	if workExists == 0 && len(jobs.Items) != 0 && jc.IsJobFinished(jobs.Items[0]) == true {
		fmt.Println("Deleting finished job")
		err := jc.DeleteJob(k, jobs.Items[0])
		if err != nil {
			return lock, err
		}

		fmt.Println("Unlocking the queue")
		lock = false
	}
	return lock, nil
}

func main() {
	lock := false
	r := jc.GetRedisClient()
	k, err := jc.KubeClientInCluster()
	if err != nil {
		panic(err.Error())
	}
	workInterval := time.NewTicker(1 * time.Second).C
	lockInterval := time.NewTicker(333 * time.Millisecond).C

	for {
		select {
		case <-workInterval:
			fmt.Println("Checking for work to do")
			fmt.Println("lock status: ", lock)
			lock, err = attemptJobCreation(k, lock, r)
			if err != nil {
				fmt.Println(err.Error())
			}

		case <-lockInterval:
			fmt.Println("Checking to see if it's time to unlock the queue.")
			fmt.Println("lock status: ", lock)
			lock, err = attemptQueueUnlock(k, lock, r)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
