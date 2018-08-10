package jc

import (
    "k8s.io/client-go/kubernetes"
    "github.com/go-redis/redis"
)

type RedisClient = redis.Client

type KubeClient struct {
    clientset kubernetes.Interface
}
