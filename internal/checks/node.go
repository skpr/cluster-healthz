package checks

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NodeStatus reviews the status of all Node objects.
func NodeStatus(clientset kubernetes.Interface) ([]Error, error) {
	var list []Error

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return list, err
	}

	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionFalse {
				list = append(list, Error{
					Name:        node.ObjectMeta.Name,
					Issue:       "NodeNotReady",
					Description: "The kubelet is not healthy or ready to accept pods.",
					Command:     fmt.Sprintf("kubectl describe node %s", node.ObjectMeta.Name),
				})
			}
		}
	}

	return list, nil
}
