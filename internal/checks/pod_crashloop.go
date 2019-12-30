package checks

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// PodCrashLoopBackOff reviews the status of all Pod marked with a label to detmine if they are in a CrashLoopBackOff.
func PodCrashLoopBackOff(clientset kubernetes.Interface) ([]Issue, error) {
	var issues []Issue

	list, err := clientset.CoreV1().Pods(corev1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.FormatLabels(map[string]string{
			LabelClusterHealthz: LevelCritical,
		}),
	})
	if err != nil {
		return issues, err
	}

	for _, pod := range list.Items {
		for _, container := range pod.Status.ContainerStatuses {
			if container.State.Waiting == nil {
				continue
			}

			if container.State.Waiting.Reason == ReasonCrashLoopBackOff {
				issues = append(issues, Issue{
					Name:        fmt.Sprintf("%s~%s", pod.ObjectMeta.Name, container.Name),
					Issue:       ReasonCrashLoopBackOff,
					// @todo, Make this more specific eg. OOM
					Description: "The Pod has restarted too many times.",
					Command:     fmt.Sprintf("kubectl -n %s describe pod %s", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name),
				})
			}
		}
	}

	return issues, nil
}
