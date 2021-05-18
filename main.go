package main

import (
	"bufio"
	"context"
	"fmt"
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
	fmt.Print("Task (view or create): ")
	task := readInput(reader)
	if task == "view" {
		clientset := connectToK8s()
		fmt.Print("Namespace (empty for all): ")
		namespace := readInput(reader)
		getPods(clientset, &namespace)
	} else if task == "create" {
		clientset := connectToK8s()
		fmt.Print("Namespace: ")
		namespace := readInput(reader)
		fmt.Print("Deployment name: ")
		deploymentName := readInput(reader)
		fmt.Print("Container image: ")
		image := readInput(reader)
		launchK8sDeployment(clientset, &namespace, &deploymentName, &image)
	} else if task == "exit" {
		os.Exit(0)
	}
}

func readInput(reader *bufio.Reader) string {
	result, err := reader.ReadString('\n')
	if err != nil {
		fmt.Errorf("Cannot read input: %v", err.Error())
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

func launchK8sDeployment(clientset *kubernetes.Clientset, namespace *string, deploymentName *string,
	image *string) {
	if *namespace == "" {
		*namespace = "default"
	}

	deployments := clientset.CoreV1().Pods(*namespace)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: *deploymentName,
			Namespace: *namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: *deploymentName,
					Image: *image,
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	_, err := deployments.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create K8s deployment: %v", err.Error())
	}

	log.Println("Created K8s deployment.")
}

func getPods(clientset *kubernetes.Clientset, namespace *string) {
	if *namespace == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Cannot get list of namespaces: %v", err.Error())
		}
		for _, ns := range namespaces.Items {
			name := ns.Name
			getPodsOfNamespace(clientset, name)
		}
	} else {
		getPodsOfNamespace(clientset, *namespace)
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
