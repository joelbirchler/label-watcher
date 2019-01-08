package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// FIXME
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "support-stage-config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	api := clientset.CoreV1()
	nodes, err := api.Nodes().List(metaV1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", nodes)
}
