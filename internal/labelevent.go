package internal

import (
	"fmt"

	"k8s.io/api/core/v1"
)

// LabelEvent describes changes to a node's labels
type LabelEvent struct {
	NodeName string
	Modified map[string]string
	Added    map[string]string
	Removed  map[string]string
}

// NewLabelEvent creates a LabelEvent from the diff of a previous node and a current one
func NewLabelEvent(prev, current *v1.Node) *LabelEvent {
	event := &LabelEvent{
		NodeName: current.Name,
		Modified: make(map[string]string),
		Added:    make(map[string]string),
		Removed:  make(map[string]string),
	}

	switch {
	case prev == nil && current != nil:
		event.Added = current.Labels
	case prev != nil && current == nil:
		event.Removed = prev.Labels
	case prev != nil && current != nil:
		event.addModified(prev, current)
		event.addAdded(prev, current)
		event.addRemoved(prev, current)
	}

	return event
}

func (event *LabelEvent) String() string {
	s := fmt.Sprintf("%s:\n", event.NodeName)

	actions := []struct {
		bucket map[string]string
		symbol string
	}{
		{event.Modified, "Δ"},
		{event.Added, "+"},
		{event.Removed, "-"},
	}

	for _, action := range actions {
		for k, v := range action.bucket {
			s += fmt.Sprintf("%s %s = %s\n", action.symbol, k, v)
		}
	}

	return s
}

// addModified looks for any changed values with the same key, and adds them to the event
func (event *LabelEvent) addModified(prev, current *v1.Node) {
	for key, prevVal := range prev.Labels {
		curVal, ok := current.Labels[key]
		if ok && curVal != prevVal {
			event.Modified[key] = curVal
		}
	}
}

// addAdded looks for any new keys in current that were not found in previous, and adds them to the event
func (event *LabelEvent) addAdded(prev, current *v1.Node) {
	for k, v := range current.Labels {
		if _, ok := prev.Labels[k]; !ok {
			event.Added[k] = v
		}
	}
}

// addRemoved looks for any keys in previous that were dropped in current, and adds them to the event
func (event *LabelEvent) addRemoved(prev, current *v1.Node) {
	for k, v := range prev.Labels {
		if _, ok := current.Labels[k]; !ok {
			event.Removed[k] = v
		}
	}
}
