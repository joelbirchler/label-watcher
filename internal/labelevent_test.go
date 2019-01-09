package internal

import (
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func fakeNode(name string) *v1.Node {
	return &v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}
}

func TestNewLabelEvent(t *testing.T) {
	event := NewLabelEvent(fakeNode("node00"), fakeNode("node00"))
	if event.NodeName != "node00" {
		t.Errorf("NodeName was incorrect, got: %s, want: %s", event.NodeName, "node00")
	}
}
