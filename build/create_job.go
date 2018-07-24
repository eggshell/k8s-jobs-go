package main

import (
    "fmt"
    "time"
    /* oidc needed for auth to IBM Cloud but is not referenced specifically in
    this script */
    _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
    v1 "k8s.io/api/core/v1"
    batchv1 "k8s.io/api/batch/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

type Client struct {
    clientset kubernetes.Interface
}

func NewClientInCluster() (*Client, error) {
    // gets in-cluster config using serviceaccount token
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err.Error())
    }
    // creates the clientset from config
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    return &Client{
        clientset: clientset,
    }, nil
}

func ConstructJob() *batchv1.Job {
    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            GenerateName: "whalesay-job-",
            Namespace: "default",
        },
        Spec: batchv1.JobSpec{
            Template: v1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    GenerateName: "whalesay-job-",
                },
                Spec: v1.PodSpec{
                    Containers: []v1.Container{
                        {
                            Name:  "whalesay",
                            Image: "docker/whalesay",
                        },
                    },
                    RestartPolicy: v1.RestartPolicyOnFailure,
                },
            },
        },
    }

    return job
}

func CreateJobs(t time.Time) {
    // get k8s client
    c, err := NewClientInCluster()
    // create jobs client
    jobsClient := c.clientset.BatchV1().Jobs("default")
    // construct kubernetes job
    job := ConstructJob()

    // send job to kubernetes api
    fmt.Println("Creating job... ")
    result1, err1 := jobsClient.Create(job)
    if err != nil {
        fmt.Println(err1)
        panic(err1)
    }
    fmt.Printf("Created job %q.\n", result1)
}

func doEvery(d time.Duration, f func(time.Time)) {
    for x := range time.Tick(d) {
        f(x)
    }
}

func main() {
    doEvery(15000*time.Millisecond, CreateJobs)
}
