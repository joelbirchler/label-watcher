package internal

import (
	"reflect"

	"k8s.io/api/core/v1"
)

type LabelEvent struct {
	NodeName string
	Modified map[string]string
	Added    map[string]string
	Removed  map[string]string
}

func sameLabels(node1, node2 v1.Node) bool {
	return reflect.DeepEqual(node1.GetLabels(), node2.GetLabels())
}

/*func labelDiff(node1, node2 *v1.Node) map[string]string {
}*/
