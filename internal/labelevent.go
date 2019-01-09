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
	return &LabelEvent{NodeName: current.Name}
}
