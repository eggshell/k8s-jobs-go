package main

import (
    "flag"
    "fmt"
    "k8s.io/client-go/util/homedir"
    "path/filepath"
    /* oidc needed for auth to IBM Cloud but is not referenced specifically in
    this script */
    _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
    v1 "k8s.io/api/core/v1"
    batchv1 "k8s.io/api/batch/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    var kubeconfig *string
    if home := homedir.HomeDir(); home != "" {
        kubeconfig = flag.String("kubeconfig", filepath.Join(home, "kubeconfig.yml"), "(optional) absolute path to the kubeconfig file")
    } else {
        kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }
    flag.Parse()

    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        panic(err)
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    jobsClient := clientset.BatchV1().Jobs("default")
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
