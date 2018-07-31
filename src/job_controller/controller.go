package job_controller

import (
    "fmt"
    v1 "k8s.io/api/core/v1"
    batchv1 "k8s.io/api/batch/v1"
)

// IsJobFinished returns whether the given Job has finished or not
func IsJobFinished(job batchv1.Job) bool {
    return job.Status.Succeeded > 0
}

// IsPodFinished returns whether the given Pod has finished or not
func IsPodFinished(pod v1.Pod) bool {
    return pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed
}

// Starts a new set of jobs
func StartNewJobSet() {
    streams := UpdateStreamsList()
    for k := range streams {
        fmt.Println(streams[k])
        //CreateJob()
    }
}
