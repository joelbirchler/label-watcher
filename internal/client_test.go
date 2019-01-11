package internal

import (
	"reflect"
	"testing"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func newFakeClient() corev1.CoreV1Interface {
	fake.NewSimpleClientset().Discovery()
	return fake.NewSimpleClientset().CoreV1()
}

func mockNode(labels map[string]string) *v1.Node {
	return &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "myNode",
			Labels: labels,
		},
	}
}

func TestWatchNodeLabels(t *testing.T) {
	tables := []struct {
		node  *v1.Node
		event LabelEvent
	}{
		{
			mockNode(map[string]string{"a": "1", "b": "2"}),
			LabelEvent{
				NodeName: "myNode",
				Modified: map[string]string{},
				Added:    map[string]string{"a": "1", "b": "2"},
				Removed:  map[string]string{},
			},
		}, {
			mockNode(map[string]string{"b": "42", "c": "3"}),
			LabelEvent{
				NodeName: "myNode",
				Modified: map[string]string{"b": "42"},
				Added:    map[string]string{"c": "3"},
				Removed:  map[string]string{"a": "1"},
			},
		}, {
			mockNode(map[string]string{"b": "42"}),
			LabelEvent{
				NodeName: "myNode",
				Modified: map[string]string{},
				Added:    map[string]string{},
				Removed:  map[string]string{"c": "3"},
			},
		},
	}

	api := newFakeClient()
	ch := WatchNodeLabels(api)

	// Create node and update labels
	go func() {
		api.Nodes().Create(tables[0].node)
		for _, table := range tables[1:] {
			api.Nodes().Update(table.node)
		}
	}()

	// Check that our WatchNodeLabels read channel receives correct LabelEvents
	for _, table := range tables {
		select {
		case event := <-ch:
			if !reflect.DeepEqual(event, table.event) {
				t.Errorf("WatchNodeLabels sent incorrect LabelEvent, got: %s, want: %s", event, table.event)
			}
		case <-time.After(10 * time.Second):
			t.Errorf("Timeout waiting for node event")
		}
	}
}
