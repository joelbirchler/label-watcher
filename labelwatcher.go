package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	api := connect()

	watcher, err := api.Nodes().Watch(metaV1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	nodes := make(map[string]*v1.Node)

	ch := watcher.ResultChan()
	for event := range ch {
		n, ok := event.Object.(*v1.Node)
		if !ok {
			log.Fatal("watch received unexpected object")
		}

		existing, ok := nodes[n.Name]
		if ok {
			fmt.Printf("Existing (%v)...", sameLabels(existing, n))
		} else {
			nodes[n.Name] = n
		}

		printNode(*n)
	}
}

func printNode(node v1.Node) {
	fmt.Println(node.Name)
	for k, v := range node.GetLabels() {
		fmt.Printf("  %s = %s\n", k, v)
	}
	fmt.Println()
}

func sameLabels(node1, node2 *v1.Node) bool {
	return reflect.DeepEqual(node1.GetLabels(), node2.GetLabels())
}

func connect() corev1.CoreV1Interface {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
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
