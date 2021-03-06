package checks

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NodeStatus reviews the status of all Node objects.
func NodeStatus(clientset kubernetes.Interface) ([]Issue, error) {
	var list []Issue

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return list, err
	}

	// Allows for new Nodes to become Ready without triggering alerts.
	warmup := metav1.NewTime(time.Now().UTC().Add(-5 * time.Minute))

	for _, node := range nodes.Items {
		if !node.ObjectMeta.CreationTimestamp.Before(&warmup) {
			// This Node is new and is still becoming Ready.
			continue
		}

		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionFalse {
				list = append(list, Issue{
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
