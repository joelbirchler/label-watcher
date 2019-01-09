package internal

import (
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

	// FIXME: tidy this up
	if prev != nil && current != nil {
		for key, prevVal := range prev.Labels {
			curVal, ok := current.Labels[key]
			if ok {
				if curVal != prevVal {
					event.Modified[key] = curVal
				}
			} else {
				event.Removed[key] = prevVal
			}
		}

		for key, curVal := range current.Labels {
			_, ok := prev.Labels[key]
			if !ok {
				event.Added[key] = curVal
			}
		}
	}

	if prev == nil && current != nil {
		event.Added = current.Labels
	}

	return event
}
