package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikel-jason/kube-power-forward/pkg/cmd/powerforward"
	"github.com/mikel-jason/kube-power-forward/pkg/proxy"
)

func main() {

	cfg := &powerforward.Config{
		Forwards: []powerforward.Forward{
			{
				Namespace:   "default",
				ServiceName: "echoserver",
				PodPort:     5678,
				LocalPort:   8080,
			},
		},
	}

	if err := powerforward.Start(cfg); err != nil {
		log.Fatalln(err)
	}

	proxy.Start(
		&proxy.Config{
			SocksListenAddr: "0.0.0.0:8888",
			HttpListenAddr:  "0.0.0.0:80",
			Mappings: map[string]int{
				"everything.k8s.proxy": 8080,
			},
		},
	)

	log.Println("up and running")

	signalForShutdownChan := make(chan os.Signal, 1)
	signal.Notify(signalForShutdownChan, os.Interrupt, syscall.SIGTERM)

	<-signalForShutdownChan
	log.Println("received shutdown signal")
	proxy.Stop()

	log.Println("stopping forwarders")
	powerforward.Stop()
}
