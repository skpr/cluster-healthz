package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestPodCrashLoopBackOff(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fail",
				Namespace: "kube-system",
				Labels: map[string]string{
					LabelClusterHealthz: LevelCritical,
				},
			},
			Status: corev1.PodStatus{
				ContainerStatuses: []corev1.ContainerStatus{
					{
						Name: "container1",
						State: corev1.ContainerState{
							Waiting: &corev1.ContainerStateWaiting{
								Reason: ReasonCrashLoopBackOff,

							},
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pass",
				Namespace: "kube-system",
				Labels: map[string]string{
					LabelClusterHealthz: LevelCritical,
				},
			},
		},
	}

	for _, pod := range pods {
		_, err := clientset.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(&pod)
		assert.Nil(t, err)
	}

	issues, err := PodCrashLoopBackOff(clientset)
	assert.Nil(t, err)

	expected := []Issue{
		{
			Name:        "fail~container1",
			Issue:       ReasonCrashLoopBackOff,
			Description: "The Pod has restarted too many times.",
			Command:     "kubectl -n kube-system describe pod fail",
		},
	}

	assert.Equal(t, expected, issues)
}
