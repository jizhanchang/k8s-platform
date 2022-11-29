package lib

import (
	"k8s-platfrom/lib/k8s"
	"k8s.io/client-go/kubernetes"
)

var K8s *kubernetes.Clientset

func Init(kubeConfig string) (err error) {

	K8s, err = k8s.Init(kubeConfig)
	return
}
