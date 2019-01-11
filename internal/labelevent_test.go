package internal

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
)

func TestNewLabelEvent(t *testing.T) {
	tables := []struct {
		prev, current *v1.Node
		event         LabelEvent
	}{
		{
			mockNode(map[string]string{"a": "1", "b": "2"}),
			mockNode(map[string]string{"a": "4", "c": "3"}),
			LabelEvent{
				NodeName: "myNode",
				Modified: map[string]string{"a": "4"},
				Added:    map[string]string{"c": "3"},
				Removed:  map[string]string{"b": "2"},
			},
		},
		{
			mockNode(nil),
			mockNode(map[string]string{"a": "1", "b": "2"}),
			LabelEvent{
				NodeName: "myNode",
				Modified: map[string]string{},
				Added:    map[string]string{"a": "1", "b": "2"},
				Removed:  map[string]string{},
			},
		},
		{
			mockNode(map[string]string{"a": "1", "b": "2"}),
			mockNode(map[string]string{"a": "1", "b": "2"}),
			LabelEvent{
				NodeName: "myNode",
				Modified: map[string]string{},
				Added:    map[string]string{},
				Removed:  map[string]string{},
			},
		},
		{
			mockNode(map[string]string{"a": "1", "b": "2"}),
			mockNode(map[string]string{}),
			LabelEvent{
				NodeName: "myNode",
				Modified: map[string]string{},
				Added:    map[string]string{},
				Removed:  map[string]string{"a": "1", "b": "2"},
			},
		},
		{
			mockNode(map[string]string{"a": "1"}),
			mockNode(map[string]string{"a": "2"}),
			LabelEvent{
				NodeName: "myNode",
				Modified: map[string]string{"a": "2"},
				Added:    map[string]string{},
				Removed:  map[string]string{},
			},
		},
	}

	for _, table := range tables {
		event := NewLabelEvent(table.prev, table.current)

		if !reflect.DeepEqual(*event, table.event) {
			t.Errorf("Event was incorrect, got: %s, want: %s", *event, table.event)
		}
	}
}
