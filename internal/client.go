package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// Connect creates a CoreV1 client for external access to the Kubernetes API. For this demo,
// we only support reading from ~/.kube/config. The KUBECONFIG env and list of configs are not
// supported.
func Connect() (corev1.CoreV1Interface, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1(), nil
}

// WatchNodeLabels watches nodes for label changes. LabelEvents will be sent on the
// returned channel upon startup and on subsequent label updates.
//
// TODO: return <-chan LabelEvent (name only) and print
// TODO: create LabelEvents from existing + delta node objects
// TODO: channel closing
func WatchNodeLabels(api corev1.CoreV1Interface) {
	watcher, err := api.Nodes().Watch(metaV1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	nodes := make(map[string]v1.Node)

	ch := watcher.ResultChan()
	for event := range ch {
		if n, ok := event.Object.(*v1.Node); ok {
			eventNode := *n

			// FIXME: messy
			existingNode, ok := nodes[eventNode.Name]
			if ok {
				fmt.Printf("Existing (%v)...", sameLabels(existingNode, eventNode))
			} else {
				nodes[eventNode.Name] = eventNode
			}

			printNode(eventNode)

		} else {
			log.Println("Watch received unexpected object.")
		}
	}
}

// FIXME: This goes somewhere, but not here (remember that the server will use the same writer interface)
// First factor printing out of watcher. We just want a read channel
func printNode(node v1.Node) {
	fmt.Println(node.Name)
	for k, v := range node.GetLabels() {
		fmt.Printf("  %s = %s\n", k, v)
	}
	fmt.Println()
}
