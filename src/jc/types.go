package jc

import (
    "k8s.io/client-go/kubernetes"
    "github.com/go-redis/redis"
)

// RedisClient is an alias to redis.Client in the go-redis package
// Existence makes it easy to use a redis client object in main.go without
// having to import go-redis there.
type RedisClient = redis.Client

// KubeClient represents the wrapper of a Kubernetes API client
type KubeClient struct {
    clientset kubernetes.Interface
}
