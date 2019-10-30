package checks

import (
	"k8s.io/client-go/kubernetes"
)

// ErrorList from all checks.
func ErrorList(clientset kubernetes.Interface) ([]Error, error) {
	var list []Error

	funcs := []func(clientset kubernetes.Interface) ([]Error, error){
		AutoscalerStatus,
		NodeStatus,
	}

	for _, f := range funcs {
		errors, err := f(clientset)
		if err != nil {
			return list, err
		}

		list = append(list, errors...)
	}

	return list, nil
}