package powerforward

import (
	"context"
	"fmt"

	"github.com/mikel-jason/kube-power-forward/pkg/kube"
)

var forwarders []kube.Forwarder

// Config is the single entry to run a powerforward command
type Config struct {
	Forwards []Forward `json:"forward,omitempty"`
}

type Forward struct {
	Namespace   string `json:"namespace"`
	ServiceName string `json:"serviceName"`
	PodPort     int    `json:"podPort"`
	LocalPort   int    `json:"localPort"`
}

func Start(cfg *Config) error {
	client, err := kube.NewClient()
	if err != nil {
		return fmt.Errorf("cannot create kube client: %v", err)
	}

	forwarders = make([]kube.Forwarder, len(cfg.Forwards))
	for forwardIndex, forward := range cfg.Forwards {
		forwardService := client.Service(forward.Namespace, forward.ServiceName)
		forwarder := client.Forwarder(forwardService, kube.ForwarderOptions{
			PodPort:   forward.PodPort,
			LocalPort: forward.LocalPort,
		})
		err := forwarder.Start(context.TODO())
		if err != nil {
			Stop()
			return fmt.Errorf("could not start forwarder for service %s/%s: %w", forward.Namespace, forward.ServiceName, err)
		}
		forwarders[forwardIndex] = forwarder
	}
	return nil
}

func Stop() {
	for _, forwarder := range forwarders {
		if forwarder == nil { // for early exit or non-started
			continue
		}
		forwarder.Stop()
	}
}
