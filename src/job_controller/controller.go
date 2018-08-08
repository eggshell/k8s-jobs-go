package job_controller

import (
    "fmt"
    "os"
    "github.com/go-redis/redis"
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

func DeleteJob(c *Client, job batchv1.Job) error {
    if err := c.clientset.BatchV1().Jobs(job.Namespace).Delete(job.Name, &metav1.DeleteOptions{}); err != nil {
        return err
    }

    return nil
}

func ListJobs(c *Client, namespace string) (*batchv1.JobList, error) {
    jobs, err := c.clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
    if err != nil {
        return nil, err
    }

    return jobs, nil
}


// IsJobFinished returns whether the given Job has finished or not
func IsJobFinished(j batchv1.Job) bool {
    for _, c := range j.Status.Conditions {
        if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == v1.ConditionTrue {
            return true
        }
    }
    return false
}

// IsPodFinished returns whether the given Pod has finished or not
func IsPodFinished(pod v1.Pod) bool {
    return pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed
}

// TODO: figure out if this jobspec actually works
// ref: https://kubernetes.io/docs/tasks/job/coarse-parallel-processing-work-queue/
func ConstructJob(workItems []string) *batchv1.Job {
    count := int32(len(workItems))

    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            GenerateName: "job-wq-",
            Namespace: "default",
        },
        Spec: batchv1.JobSpec{
            Completions: &count,
            Parallelism: &count,
            Template: v1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    GenerateName: "job-wq-1",
                },
                Spec: v1.PodSpec{
                    Containers: []v1.Container{
                        {
                            Name:  "sp",
                            Image: "registry.ng.bluemix.net/eggshell/rotisserie-sp:4472956",
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

func CreateJob(client redis.Client, workItems []string) {
    err := RenameReadKey(client)
    if err != nil {
        fmt.Println(err)
    }

    c, err := KubeClientInCluster()
    jobsClient := c.clientset.BatchV1().Jobs("default")
    job := ConstructJob(workItems)

    fmt.Println("Creating job...")
    result, err := jobsClient.Create(job)
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println("Created job %q.", result)
}
