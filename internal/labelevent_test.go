package internal

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func fakeNode(labels map[string]string) *v1.Node {
	return &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "myNode",
			Labels: labels,
		},
	}
}

func TestNewLabelEvent(t *testing.T) {
	tables := []struct {
		prev, current,
		modified, added, removed map[string]string
	}{
		{
			map[string]string{"a": "1", "b": "2"}, // prev
			map[string]string{"a": "4", "c": "3"}, // current
			map[string]string{"a": "4"},           // modified
			map[string]string{"c": "3"},           // added
			map[string]string{"b": "2"},           // removed
		},
		{
			nil, // prev
			map[string]string{"a": "1", "b": "2"}, // current
			map[string]string{},                   // modified
			map[string]string{"a": "1", "b": "2"}, // added
			map[string]string{},                   // removed
		},
		{
			map[string]string{"a": "1", "b": "2"}, // prev
			map[string]string{"a": "1", "b": "2"}, // current
			map[string]string{},                   // modified
			map[string]string{},                   // added
			map[string]string{},                   // removed
		},
		{
			map[string]string{"a": "1", "b": "2"}, // prev
			map[string]string{},                   // current
			map[string]string{},                   // modified
			map[string]string{},                   // added
			map[string]string{"a": "1", "b": "2"}, // removed
		},
		{
			map[string]string{"a": "1"}, // prev
			map[string]string{"a": "2"}, // current
			map[string]string{"a": "2"}, // modified
			map[string]string{},         // added
			map[string]string{},         // removed
		},
	}

	for _, table := range tables {
		event := NewLabelEvent(fakeNode(table.prev), fakeNode(table.current))

		if !reflect.DeepEqual(event.Modified, table.modified) {
			t.Errorf("Modified was incorrect, got: %s, want: %s", event.Modified, table.modified)
		}

		if !reflect.DeepEqual(event.Added, table.added) {
			t.Errorf("Added was incorrect, got: %s, want: %s", event.Added, table.added)
		}

		if !reflect.DeepEqual(event.Removed, table.removed) {
			t.Errorf("Removed was incorrect, got: %s, want: %s", event.Removed, table.removed)
		}
	}
}
