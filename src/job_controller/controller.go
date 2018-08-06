package job_controller

import (
    "fmt"
    "os"
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

// IsJobFinished returns whether the given Job has finished or not
func IsJobFinished(job batchv1.Job) bool {
    return job.Status.Succeeded > 0
}

// IsPodFinished returns whether the given Pod has finished or not
func IsPodFinished(pod v1.Pod) bool {
    return pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed
}

// TODO: figure out if this jobspec actually works
// ref: https://kubernetes.io/docs/tasks/job/coarse-parallel-processing-work-queue/
func ConstructJob(workItems []string) *batchv1.Job {
    compCount := int32(1)
    podCount := int32(len(workItems))

    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            GenerateName: "job-wq-",
            Namespace: "default",
        },
        Spec: batchv1.JobSpec{
            Completions: &compCount,
            Parallelism: &podCount,
            Template: v1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    GenerateName: "job-wq-",
                },
                Spec: v1.PodSpec{
                    Containers: []v1.Container{
                        {
                            Name:  "sp",
                            Image: "registry.ng.bluemix.net/eggshell/rotisserie-sp:d02aaf2",
                            Env: []v1.EnvVar{
                                {
                                    Name: "REDIS_PASSWORD",
                                    Value: os.Getenv("REDIS_PASSWORD"),
                                },
                                {
                                    Name: "TOKEN",
                                    Value: os.Getenv("TOKEN"),
                                },
                            },
                        },
                    },
                    RestartPolicy: v1.RestartPolicyNever,
                },
            },
        },
    }

    return job
}

func CreateJob(workItems []string) {
    c, err := KubeClientInCluster()
    jobsClient := c.clientset.BatchV1().Jobs("default")
    job := ConstructJob(workItems)

    fmt.Println("Creating job...")
    result, err := jobsClient.Create(job)
    if err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println("Created job %q.", result)
    }
}
