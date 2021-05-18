package main

import (
	"bufio"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

// reference: https://dev.to/narasimha1997/create-kubernetes-jobs-in-golang-using-k8s-client-go-api-59ej

func main() {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		handleK8sCommand(stdReader)
	}
}

func handleK8sCommand(reader *bufio.Reader) {
	fmt.Print("Task (view, create, or delete): ")
	task := readInput(reader)
	if task == "view" {
		clientset := connectToK8s()
		fmt.Print("Namespace (empty for all): ")
		namespace := readInput(reader)
		getPods(clientset, namespace)
	} else if task == "create" {
		clientset := connectToK8s()
		fmt.Print("Namespace: ")
		namespace := readInput(reader)
		fmt.Print("App name: ")
		appName := readInput(reader)
		fmt.Print("Deployment name: ")
		deploymentName := readInput(reader)
		fmt.Print("Container name: ")
		containerName := readInput(reader)
		fmt.Print("Container image: ")
		image := readInput(reader)
		launchK8sDeployment(clientset, namespace, appName, deploymentName, containerName, image)
	} else if task == "delete" {
		clientset := connectToK8s()
		fmt.Print("Namespace: ")
		namespace := readInput(reader)
		fmt.Print("Deployment name: ")
		deploymentName := readInput(reader)
		deleteK8sDeployment(clientset, namespace, deploymentName)
	} else if task == "exit" {
		os.Exit(0)
	} else {
		log.Printf("Invalid task type.")
	}
}

func readInput(reader *bufio.Reader) string {
	result, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Cannot read input: %v", err.Error())
	}
	return result[:(len(result) - 1)]
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

// https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go
func launchK8sDeployment(
	clientset *kubernetes.Clientset,
	namespace string,
	appName string,
	deploymentName string,
	containerName string,
	image string) {
	if namespace == "" {
		namespace = "default"
	}
	var numOfReplicas int32 = 4

	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &numOfReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": appName,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: containerName,
							Image: image,
							Ports: []v1.ContainerPort{
								{
									Name: "http",
									Protocol: v1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Cannot create deployment: %v", err.Error())
	}
	log.Printf("Created deployment %v.", result.GetObjectMeta().GetName())
}

func deleteK8sDeployment(
	clientset *kubernetes.Clientset,
	namespace string,
	deploymentName string) {
	deletePolicy := metav1.DeletePropagationForeground
	deploymentClient := clientset.AppsV1().Deployments(namespace)

	if err := deploymentClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Fatalf("Cannot delete deployment: %v", err.Error())
	}
	log.Printf("Deleted deployment %v", deploymentName)
}

func getPods(clientset *kubernetes.Clientset, namespace string) {
	if namespace == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Cannot get list of namespaces: %v", err.Error())
		}
		for _, ns := range namespaces.Items {
			name := ns.Name
			getPodsOfNamespace(clientset, name)
		}
	} else {
		getPodsOfNamespace(clientset, namespace)
	}
}

func getPodsOfNamespace(clientset *kubernetes.Clientset, name string) {
	pods, err := clientset.CoreV1().Pods(name).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Cannot get pods of namespace %v: %v", name, err.Error())
	}

	for _, p := range pods.Items {
		log.Printf("Pod of namespace %v: %v", name, p.Name)
	}
}
