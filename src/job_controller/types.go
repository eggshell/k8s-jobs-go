package job_controller

import (
    "k8s.io/client-go/kubernetes"
)

type Client struct {
    clientset kubernetes.Interface
}

type Stream struct {
    displayName string
    alive       int
}
