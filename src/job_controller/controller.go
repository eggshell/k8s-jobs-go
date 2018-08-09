package job_controller

import (
    "os"
    "github.com/go-redis/redis"
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
        return nil, err
    }
    // creates the clientset from config
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }

    return &Client{
        clientset: clientset,
    }, nil
}

func DeleteJob(c *Client, job batchv1.Job) error {
    var policy metav1.DeletionPropagation = "Background"

    if err := c.clientset.BatchV1().Jobs(job.Namespace).Delete(job.Name,
    &metav1.DeleteOptions{PropagationPolicy: &policy}); err != nil {
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

func IsJobFinished(j batchv1.Job) bool {
    for _, c := range j.Status.Conditions {
        if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == v1.ConditionTrue {
            return true
        }
    }
    return false
}

func ConstructJob(workItems []string) *batchv1.Job {
    count := int32(len(workItems))

    jobSpec := &batchv1.Job{
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
                            Image: "registry.ng.bluemix.net/eggshell/rotisserie-sp:f7095a1",
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
    return jobSpec
}

func CreateJob(kubeClient *Client, redisClient redis.Client, workItems []string) (*batchv1.Job, error) {
    err := RenameReadKey(redisClient)
    if err != nil {
        return nil,err
    }

    jobsClient  := kubeClient.clientset.BatchV1().Jobs("default")
    jobSpec     := ConstructJob(workItems)
    job, err := jobsClient.Create(jobSpec)
    if err != nil {
        return nil,err
    }

    return job, nil
}
