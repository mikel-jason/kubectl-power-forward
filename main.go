package main

import (
	"context"
	"log"

	"github.com/mikel-jason/kube-power-forward/pkg/kube"
)

func main() {

	client, err := kube.NewClient()
	if err != nil {
		log.Fatalf("cannot create kube client: %v", err)
	}

	service := client.Service("forecastle", "forecastle")

	forwarder := client.Forwarder(service, kube.ForwarderOptions{
		PodPort:   3000,
		LocalPort: 3000,
	})

	err = forwarder.Start(context.TODO())
	if err != nil {
		log.Fatalf("cannot port-forward: %v", err)
	}
}
