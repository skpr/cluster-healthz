package checks

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/client-go/kubernetes"
)

// NodeStatus reviews the status of all HorizontalPodAutoscaler objects.
func AutoscalerStatus(clientset kubernetes.Interface) ([]Error, error) {
	var list []Error

	hpas, err := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return list, err
	}

	for _, hpa := range hpas.Items {
		for _, condition := range hpa.Status.Conditions {
			if condition.Type == autoscalingv2beta2.ScalingActive && condition.Status == corev1.ConditionFalse {
				list = append(list, Error{
					Namespace: hpa.ObjectMeta.Namespace,
					Name:        hpa.ObjectMeta.Name,
					Issue:       "NodeScalingIssue",
					Command: fmt.Sprintf("kubectl -n %s describe hpa %s", hpa.ObjectMeta.Namespace, hpa.ObjectMeta.Name),
				})
			}
		}
	}

	return list, nil
}