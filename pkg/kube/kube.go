package kube

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client interface {
	Service(namespace string, name string) Service
	Forwarder(service Service, options ForwarderOptions) Forwarder
}

type client struct {
	restConfig   *rest.Config
	k8sClientSet *kubernetes.Clientset
}

func NewClient() (Client, error) {
	var err error
	c := &client{}

	c.restConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		log.Println("error creating config")
		log.Fatalln(err)
	}

	c.k8sClientSet, err = kubernetes.NewForConfig(c.restConfig)
	if err != nil {
		log.Println("error creating k8s clientset")
		log.Fatalln(err)
	}

	return c, nil
}
