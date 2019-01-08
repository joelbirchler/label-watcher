package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	api := connect()

	nodes, err := api.Nodes().List(metaV1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range nodes.Items {
		fmt.Println(node.Name)
		for k, v := range node.GetLabels() {
			fmt.Printf("  %s = %s\n", k, v)
		}
		fmt.Println()
	}
}

func connect() corev1.CoreV1Interface {
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

	return clientset.CoreV1()
}
