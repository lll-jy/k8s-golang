package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	getJobs()
}

// https://rancher.com/using-kubernetes-api-go-kubecon-2017-session-recap
/*func getClient(pathToCfg string) (*kubernetes.Clientset, error) {
	var config *rest.Confit
	var err error
	return nil, nil
}*/

// https://dev.to/narasimha1997/create-kubernetes-jobs-in-golang-using-k8s-client-go-api-59ej
func getJobs() {
	jobName := flag.String("jobname", "test-job", "The name of the job")
	containerImage := flag.String("image", "ubuntu:latest", "Name of the container image")
	entryCommand := flag.String("command", "ls", "The command to run inside the container")

	flag.Parse()

	clientset := connectToK8s()
	launchK8sJob(clientset, jobName, containerImage, entryCommand)

	//clientset.CoreV1().Endpoints()
	//clientset.CoreV1().ConfigMaps()
	//clientset.CoreV1().Events()
	//clientset.CoreV1().ComponentStatuses()
	//clientset.CoreV1().LimitRanges()
	//clientset.CoreV1().Namespaces()
	clientset.CoreV1().PersistentVolumeClaims()
	clientset.CoreV1().PersistentVolumes()
	//clientset.CoreV1().Nodes()
	clientset.CoreV1().Pods()
	clientset.CoreV1().PodTemplates()
	clientset.CoreV1().ReplicationControllers()
	clientset.CoreV1().ResourceQuotas()
	clientset.CoreV1().Secrets()
	clientset.CoreV1().ServiceAccounts()
	clientset.CoreV1().Services()

	clientset.DiscoveryV1().EndpointSlices()

	clientset.EventsV1().Events()

	clientset.NetworkingV1().NetworkPolicies()
	clientset.NetworkingV1().Ingresses()
	clientset.NetworkingV1().IngressClasses()

	clientset.NodeV1().RuntimeClasses()

	clientset.PolicyV1().PodDisruptionBudgets()

	clientset.RbacV1().ClusterRoleBindings()
	clientset.RbacV1().ClusterRoles()
	clientset.RbacV1().Roles()
	clientset.RbacV1().RoleBindings()

	clientset.SchedulingV1().PriorityClasses()

	clientset.StorageV1().StorageClasses()
	clientset.StorageV1().CSINodes()
	clientset.StorageV1().CSIDrivers()
	clientset.StorageV1().VolumeAttachments()
}

func connectToK8s() *kubernetes.Clientset {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}

	configPath := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Panicln("failed to create K8s config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("Failed to create K8s clientset")
	}

	return clientset
}

func launchK8sJob(clientset *kubernetes.Clientset, jobName *string, image *string, cmd *string) {
	jobs := clientset.BatchV1().Jobs("default")
	var backOffLimit int32 = 0

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: *jobName,
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: *jobName,
							Image: *image,
							Command: strings.Split(*cmd, " "),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Failed to create K8s job.")
	}

	log.Println("Created K8s job successfully.")
}