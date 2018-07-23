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

func minuteTicker() *time.Ticker {
    // Return new ticker that triggers on the minute
    return time.NewTicker(time.Second * time.Duration(60-time.Now().Second()))
}

func main() {
    // creates in-cluster config
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err.Error())
    }
    // creates the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    fmt.Println("Started ticker")
    // tick on the minute
    t := minuteTicker()
    for {
        <-t.C
        t = minuteTicker()
        fmt.Println("Calling create_job()")
        // create jobs client
        jobsClient := clientset.BatchV1().Jobs("default")
        // construct kubernetes job
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

        // send job to kubernetes api
        fmt.Println("Creating job... ")
        result1, err1 := jobsClient.Create(job)
        if err != nil {
            fmt.Println(err1)
            panic(err1)
        }
        fmt.Printf("Created job %q.\n", result1)
        fmt.Println("Listing jobs....")
        fmt.Println(jobsClient.List(metav1.ListOptions{}))
    }
}
