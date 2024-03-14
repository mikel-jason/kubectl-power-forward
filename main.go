package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikel-jason/kube-power-forward/pkg/cmd/powerforward"
)

func main() {

	cfg := &powerforward.Config{
		Forwards: []powerforward.Forward{
			{
				Namespace:   "forecastle",
				ServiceName: "forecastle",
				PodPort:     3000,
				LocalPort:   8080,
			},
			{
				Namespace:   "kube-system",
				ServiceName: "hubble-ui",
				PodPort:     8081,
				LocalPort:   8081,
			},
		},
	}

	if err := powerforward.Start(cfg); err != nil {
		log.Fatalln(err)
	}

	signalForShutdownChan := make(chan os.Signal, 1)
	signal.Notify(signalForShutdownChan, os.Interrupt, syscall.SIGTERM)

	<-signalForShutdownChan
	log.Println("received shutdown signal, stopping forwarders")
	powerforward.Stop()
}
