package checks

import "k8s.io/client-go/kubernetes"

type Check interface {
	Errors(kubernetes.Interface) ([]Error, error)
}

type Error struct {
	Namespace   string
	Name        string
	Issue       string
	Description string
	Command     string
}
