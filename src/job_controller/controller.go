package job_controller

import (
    "time"
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

// Starts a new set of jobs when hit with an http request
func StartNewJobSet(t time.Time) {
    CreateJob()
}
