package job_controller

import (
    "fmt"
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

func KubeClientInCluster() (*Client, error) {
    // gets in-cluster config using serviceaccount token
    config, err := rest.InClusterConfig()
    if err != nil {
        fmt.Println(err.Error())
    }
    // creates the clientset from config
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Println(err.Error())
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

func CreateJob(workItems string) {
    c, err := KubeClientInCluster()
    jobsClient := c.clientset.BatchV1().Jobs("default")
    job := ConstructJob()

    fmt.Println("Creating job...")
    result, err := jobsClient.Create(job)
    if err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println("Created job %q.", result)
    }
}
