package main

import (
	"bytes"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func Test(t *testing.T) {
	clientset := connectToK8s()
	deploymentName := "kubernetes-bootcamp"
	var newPods map[string]bool

	t.Run("Create", func(t *testing.T) {
		newPods = createUnitTest(t, clientset, deploymentName)
	})

	// Test log: https://stackoverflow.com/questions/44119951/how-to-check-a-log-output-in-go-test
	t.Run("View", func(t *testing.T) {
		viewUnitTest(t, clientset, newPods)
	})

	t.Run("Delete", func(t *testing.T) {
		deletedPods := deleteUnitTest(t, clientset, deploymentName)
		checkIsValidDelete(t, deletedPods, newPods)
	})
}

func checkIsValidDelete(t *testing.T, deletedPods map[string]bool, newPods map[string]bool) {
	t.Logf("Deleted pods are %v", deletedPods)
	for k, _ := range newPods {
		if _, ok := deletedPods[k]; !ok {
			t.Errorf("The pod %v created but not deleted.", k)
		}
	}
	for k, _ := range deletedPods {
		if _, ok := newPods[k]; !ok {
			t.Errorf("The pod %v deleted but was not created in test.", k)
		}
	}
}

func viewUnitTest(t *testing.T, clientset *kubernetes.Clientset, newPods map[string]bool) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	getPodsOfNamespace(clientset, "default")
	time.Sleep(10 * time.Second)

	logs := strings.Split(buf.String(), "\n")
	pods := make(map[string]bool)
	for _, l := range logs {
		if l == "" {
			break
		}
		index := strings.LastIndex(l, ":")
		pods[strings.Trim(l[(index+1):], " \n")] = true
	}

	t.Logf("New pods are %v", newPods)
	t.Logf("Pods are %v", pods)
	for k, _ := range newPods {
		if _, ok := pods[k]; !ok {
			t.Errorf("Failed to get the pod %v", k)
		}
	}
}

func createUnitTest(t *testing.T, clientset *kubernetes.Clientset, deploymentName string) map[string]bool {
	result := make(map[string]bool)
	originalDeployments := getDeploymentsOfDefaultNamespace(t, clientset)
	existingDeployments := convertDeploymentListToMapOfName(originalDeployments)
	originalPods := getPodsOfDefaultNamespace(t, clientset)
	existingPods := convertPodListToMapOfName(originalPods)

	createSampleDeployment(clientset, "demo", deploymentName)
	time.Sleep(40*time.Second)

	deployments := getDeploymentsOfDefaultNamespace(t, clientset)
	numOfNewDeployments := len(deployments) - len(originalDeployments)
	if numOfNewDeployments != 1 {
		t.Errorf("Number of deployments newly created, got: %d, want: %d.", numOfNewDeployments, 1)
	}
	pods := getPodsOfDefaultNamespace(t, clientset)
	numOfNewPods := len(pods) - len(originalPods)
	if numOfNewPods != 4 {
		t.Errorf("Number of pods newly created, got: %d, want: %d.", numOfNewPods, 4)
	}

	for _, d := range deployments {
		if _, ok := existingDeployments[d.Name]; !ok {
			t.Logf("The new deployment is: %s", d.Name)
			if d.Name != deploymentName {
				t.Errorf("Deployment newly created, got: \"%s\", want: \"%s\".", d.Name, deploymentName)
			}
		}
	}
	for _, p := range pods {
		if _, ok := existingPods[p.Name]; !ok {
			t.Logf("New pod: %s", p.Name)
			result[p.Name] = true
			if !strings.HasPrefix(p.Name, deploymentName) {
				t.Errorf("Pod newly created, got: \"%s\", want: \"^%s.*\" (regex).", p.Name, deploymentName)
			}
		}
	}
	return result
}

func createSampleDeployment(clientset *kubernetes.Clientset, appName string, deploymentName string) {
	launchK8sDeployment(
		clientset,
		"default",
		appName,
		deploymentName,
		"kubernetes-bootcamp",
		"gcr.io/google-samples/kubernetes-bootcamp:appsv1")
}

func getDeploymentsOfDefaultNamespace(t *testing.T, clientset *kubernetes.Clientset) []appsv1.Deployment {
	deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Cannot retrieve deployments of default namespace: %v", err.Error())
	}
	return deployments.Items
}

func convertDeploymentListToMapOfName(originalDeployments []appsv1.Deployment) map[string]bool {
	result := make(map[string]bool)
	for _, d := range originalDeployments {
		result[d.Name] = true
	}
	return result
}

func getPodsOfDefaultNamespace(t *testing.T, clientset *kubernetes.Clientset) []v1.Pod {
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Cannot retrieve pods of default namespace: %v", err.Error())
	}
	return pods.Items
}

func convertPodListToMapOfName(originalPods []v1.Pod) map[string]bool {
	result := make(map[string]bool)
	for _, p := range originalPods {
		result[p.Name] = true
	}
	return result
}

func deleteUnitTest(t *testing.T, clientset *kubernetes.Clientset, deploymentName string) map[string]bool {
	result := make(map[string]bool)
	originalDeployments := getDeploymentsOfDefaultNamespace(t, clientset)
	originalPods := getPodsOfDefaultNamespace(t, clientset)

	deleteSampleDeployment(clientset, deploymentName)
	time.Sleep(90*time.Second)

	deployments := getDeploymentsOfDefaultNamespace(t, clientset)
	remainingDeployments := convertDeploymentListToMapOfName(deployments)
	numOfDeletedDeployments := len(originalDeployments) - len(deployments)
	if numOfDeletedDeployments != 1 {
		t.Errorf("Number of deployments deleted, got: %d, want: %d.", numOfDeletedDeployments, 1)
	}
	pods := getPodsOfDefaultNamespace(t, clientset)
	remainingPods := convertPodListToMapOfName(pods)
	numOfDeletedPods := len(originalPods) - len(pods)
	if numOfDeletedPods != 4 {
		t.Errorf("Number of pods deleted, got: %d, want: %d.", numOfDeletedPods, 4)
	}

	for _, d := range originalDeployments {
		if _, ok := remainingDeployments[d.Name]; !ok {
			t.Logf("Deployment %v is deleted.", d.Name)
			if d.Name != deploymentName {
				t.Errorf("Deployment deleted, got: \"%s\", want: \"%s\".", d.Name, deploymentName)
			}
		}
	}
	for _, p := range originalPods {
		if _, ok := remainingPods[p.Name]; !ok {
			t.Logf("Deleted pod %v.", p.Name)
			result[p.Name] = true
			if !strings.HasPrefix(p.Name, deploymentName) {
				t.Errorf("Pod deleted, got: \"%s\", want: \"^%s.*\" (regex).", p.Name, deploymentName)
			}
		}
	}
	return result
}

func deleteSampleDeployment(clientset *kubernetes.Clientset, deploymentName string) {
	deleteK8sDeployment(
		clientset,
		"default",
		deploymentName)
}