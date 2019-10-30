package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestAutoscalerStatus(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	hpas := []autoscalingv2beta2.HorizontalPodAutoscaler{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Name: "pass",
			},
			Status: autoscalingv2beta2.HorizontalPodAutoscalerStatus{
				Conditions: []autoscalingv2beta2.HorizontalPodAutoscalerCondition{
					{
						Type:autoscalingv2beta2.ScalingActive,
						Status: corev1.ConditionTrue,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Name: "fail",
			},
			Status: autoscalingv2beta2.HorizontalPodAutoscalerStatus{
				Conditions: []autoscalingv2beta2.HorizontalPodAutoscalerCondition{
					{
						Type:autoscalingv2beta2.ScalingActive,
						Status: corev1.ConditionFalse,
					},
				},
			},
		},
	}

	for _, hpa := range hpas {
		_, err := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(hpa.ObjectMeta.Namespace).Create(&hpa)
		assert.Nil(t, err)
	}

	errors, err := AutoscalerStatus(clientset)
	assert.Nil(t, err)

	expected := []Error{
		{
			Namespace: "test",
			Name: "fail",
			Issue: "NodeScalingIssue",
			Description: "The HPA controller is unable to scale if necessary",
			Command: "kubectl -n test describe hpa fail",
		},
	}

	assert.Equal(t, expected, errors)
}