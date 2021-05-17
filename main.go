package main

import (
	"k8s.io/client-go/kubernetes"
)

func main() {

}

// https://rancher.com/using-kubernetes-api-go-kubecon-2017-session-recap
func getClient(pathToCfg string) (*kubernetes.Clientset, error) {
	var config *rest.Confit
	var err error
	return nil, nil
}

// https://dev.to/narasimha1997/create-kubernetes-jobs-in-golang-using-k8s-client-go-api-59ej