package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestErrorList(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	hpa :=	autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "fail",
		},
		Status: autoscalingv2beta2.HorizontalPodAutoscalerStatus{
			Conditions: []autoscalingv2beta2.HorizontalPodAutoscalerCondition{
				{
					Type:   autoscalingv2beta2.ScalingActive,
					Status: corev1.ConditionFalse,
				},
			},
		},
	}

	_, err := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(hpa.ObjectMeta.Namespace).Create(&hpa)
	assert.Nil(t, err)

	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fail",
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeReady,
					Status: corev1.ConditionFalse,
				},
			},
		},
	}

	_, err = clientset.CoreV1().Nodes().Create(&node)
	assert.Nil(t, err)

	issues, err := IssueList(clientset)
	assert.Nil(t, err)

	expected := []Issue{
		{
			Namespace:   "test",
			Name:        "fail",
			Issue:       "NodeScalingIssue",
			Description: "The HPA controller is unable to scale if necessary",
			Command:     "kubectl -n test describe hpa fail",
		},
		{
			Name:        "fail",
			Issue:       "NodeNotReady",
			Description: "The kubelet is not healthy or ready to accept pods.",
			Command:     "kubectl describe node fail",
		},
	}

	assert.ElementsMatch(t, expected, issues)
}
