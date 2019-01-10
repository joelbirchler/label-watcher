package internal

import (
	"log"
	"os"
	"path/filepath"
	"reflect"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
// TODO: channel closing (channel may need to be passed so can be deferred closed in main)
func WatchNodeLabels(api corev1.CoreV1Interface) <-chan LabelEvent {
	watcher, err := api.Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	nodes := make(map[string]*v1.Node)
	labelChan := make(chan LabelEvent)

	go func() {
		watcherChan := watcher.ResultChan()

		for event := range watcherChan {
			eventNode, ok := event.Object.(*v1.Node)
			if !ok {
				log.Println("Watch received unexpected object.")
			} else {
				prevNode := nodes[eventNode.Name]
				nodes[eventNode.Name] = eventNode
				// on change, send an event to the the label channel
				if prevNode == nil || eventNode == nil || !reflect.DeepEqual(prevNode.Labels, eventNode.Labels) {
					labelChan <- *NewLabelEvent(prevNode, eventNode)
				}
			}
		}
	}()

	return labelChan
}
