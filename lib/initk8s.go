package lib

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

var K8sClient *kubernetes.Clientset

func init() {
	config, err := clientcmd.BuildConfigFromFlags("", "yamls/config")
	if err != nil {
		log.Fatal(err)
	}
	config.Insecure = true
	K8sClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
}
