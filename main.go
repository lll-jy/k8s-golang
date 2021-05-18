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
	handleK8s()
}

// https://rancher.com/using-kubernetes-api-go-kubecon-2017-session-recap
/*func getClient(pathToCfg string) (*kubernetes.Clientset, error) {
	var config *rest.Confit
	var err error
	return nil, nil
}*/

// https://dev.to/narasimha1997/create-kubernetes-jobs-in-golang-using-k8s-client-go-api-59ej
func handleK8s() {
	toCreate := flag.String("tocreate", "false", "Whether to create a new job")
	jobName := flag.String("jobname", "test-job", "The name of the job")
	containerImage := flag.String("image", "ubuntu:latest", "Name of the container image")
	entryCommand := flag.String("command", "ls", "The command to run inside the container")

	flag.Parse()

	clientset := connectToK8s()
	if *toCreate == "true" {
		launchK8sJob(clientset, jobName, containerImage, entryCommand)
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Cannot get list of namespaces: %v", err.Error())
	}
	for i, ns := range namespaces.Items {
		name := ns.Name
		log.Printf("Namespace %d, %v", i, name)

		pods, err := clientset.CoreV1().Pods(name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Cannot get pods of namespace %v: %v", name, err.Error())
		}
		for j, p := range pods.Items {
			podName := p.Name
			log.Printf("Pod %d, %v", j, podName)
		}
	}
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