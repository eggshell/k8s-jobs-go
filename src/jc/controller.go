package jc

import (
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	// blank import needed for auth with certain k8s service providers
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
)

// KubeClientInCluster gets an in-cluster kubeconfig using serviceaccount token
// Returns a kubernetes client object defined in types.go
func KubeClientInCluster() (*KubeClient, error) {
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

	return &KubeClient{
		clientset: clientset,
	}, nil
}

// DeleteJob deletes a given job in a kubernetes cluster
func DeleteJob(c *KubeClient, job batchv1.Job) error {
	var policy metav1.DeletionPropagation = "Background"

	err := c.clientset.BatchV1().Jobs(job.Namespace).Delete(job.Name,
		&metav1.DeleteOptions{PropagationPolicy: &policy})
	if err != nil {
		return err
	}

	return nil
}

// ListJobs returns a JobsList object for a given kubernetes namespace
func ListJobs(c *KubeClient, namespace string) (*batchv1.JobList, error) {
	jobs, err := c.clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// IsJobFinished determines whether or not a given kubernetes job is running
// Bases this boolean value on the status of the job and its dependent pods
func IsJobFinished(j batchv1.Job) bool {
	for _, c := range j.Status.Conditions {
		if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

// ConstructJob creates a batchv1.Job object from a JSON spec and returns it
func ConstructJob(workItems []string) *batchv1.Job {
	count := int32(len(workItems))

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "job-wq-",
			Namespace:    "default",
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
									Name:  "REDIS_PASSWORD",
									Value: os.Getenv("REDIS_PASSWORD"),
								},
								{
									Name:  "TOKEN",
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

// CreateJob gets a kubernetes job spec from construct job and sends that spec
// to the kubernetes api to be created.
func CreateJob(k *KubeClient, r RedisClient, workItems []string) (*batchv1.Job, error) {
	err := RenameReadKey(r)
	if err != nil {
		return nil, err
	}

	jobsClient := k.clientset.BatchV1().Jobs("default")
	jobSpec := ConstructJob(workItems)
	job, err := jobsClient.Create(jobSpec)
	if err != nil {
		return nil, err
	}

	return job, nil
}
