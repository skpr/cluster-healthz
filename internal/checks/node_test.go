package checks

import (
	"testing"

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
			},
			Status: corev1.NodeStatus{
				Conditions:      []corev1.NodeCondition{
					{
						Type: corev1.NodeReady,
						Status: corev1.ConditionTrue,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fail",
			},
			Status: corev1.NodeStatus{
				Conditions:      []corev1.NodeCondition{
					{
						Type: corev1.NodeReady,
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

	errors, err := NodeStatus(clientset)
	assert.Nil(t, err)

	expected := []Error{
		{
			Name: "fail",
			Issue: "NodeNotReady",
			Description: "The kubelet is not healthy or ready to accept pods.",
			Command: "kubectl describe node fail",
		},
	}

	assert.Equal(t, expected, errors)
}