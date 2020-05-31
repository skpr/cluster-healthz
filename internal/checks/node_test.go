package checks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNodeStatus(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	nodes := []corev1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pass",
				CreationTimestamp: metav1.NewTime(time.Now().UTC().AddDate(0, 0, -1)),
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					{
						Type:   corev1.NodeReady,
						Status: corev1.ConditionTrue,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fail",
				CreationTimestamp: metav1.NewTime(time.Now().UTC().AddDate(0, 0, -1)),
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					{
						Type:   corev1.NodeReady,
						Status: corev1.ConditionFalse,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "skip",
				// This time will always trigger a skip.
				CreationTimestamp: metav1.NewTime(time.Now().UTC()),
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					{
						Type:   corev1.NodeReady,
						Status: corev1.ConditionFalse,
					},
				},
			},
		},
	}

	for _, node := range nodes {
		_, err := clientset.CoreV1().Nodes().Create(&node)
		assert.Nil(t, err)
	}

	issues, err := NodeStatus(clientset)
	assert.Nil(t, err)

	expected := []Issue{
		{
			Name:        "fail",
			Issue:       "NodeNotReady",
			Description: "The kubelet is not healthy or ready to accept pods.",
			Command:     "kubectl describe node fail",
		},
	}

	assert.Equal(t, expected, issues)
}
