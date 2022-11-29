package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Init(kubeConfig string) (clientSet *kubernetes.Clientset, err error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	//通过config创建clientSet
	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return
}
