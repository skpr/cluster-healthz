package checks

import (
	"k8s.io/client-go/kubernetes"
)

// IssueList from all checks.
func IssueList(clientset kubernetes.Interface) ([]Issue, error) {
	var list []Issue

	funcs := []func(clientset kubernetes.Interface) ([]Issue, error){
		AutoscalerStatus,
		PodCrashLoopBackOff,
		NodeStatus,
	}

	for _, f := range funcs {
		issues, err := f(clientset)
		if err != nil {
			return list, err
		}

		list = append(list, issues...)
	}

	return list, nil
}
