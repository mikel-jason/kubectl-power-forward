package powerforward

import (
	"context"
	"fmt"
	"log"

	"github.com/mikel-jason/kube-power-forward/pkg/kube"
	"github.com/mikel-jason/kube-power-forward/pkg/proxy"
)

var forwarders []kube.Forwarder

// Config is the single entry to run a powerforward command
type Config struct {
	SocksListenAddr       string
	HttpReverseListenAddr string
	Forwards              []Forward
}

type Forward struct {
	Namespace       string `json:"namespace"`
	ServiceName     string `json:"serviceName"`
	PodPort         int    `json:"podPort"`
	LocalPort       int    `json:"localPort"`
	Hostname        string `json:"hostname,omitempty"`
	KubeContextName string `json:"context,omitempty"`
}

func Start(cfg *Config) error {
	client, err := kube.NewClient()
	if err != nil {
		return fmt.Errorf("cannot create kube client: %v", err)
	}

	forwarders = make([]kube.Forwarder, len(cfg.Forwards))
	mappings := make(map[string]int)
	for forwardIndex, forward := range cfg.Forwards {
		ctx := kube.Context{Name: forward.KubeContextName}
		forwardService, err := client.Service(ctx, forward.Namespace, forward.ServiceName)
		if err != nil {
			return fmt.Errorf("cannot load service (context: %s, namespace: %s, service name: %s): %w",
				forward.KubeContextName, forward.Namespace, forward.ServiceName, err)
		}
		forwarder, err := client.Forwarder(ctx, forwardService, kube.ForwarderOptions{
			PodPort:   forward.PodPort,
			LocalPort: forward.LocalPort,
		})
		if err != nil {
			return fmt.Errorf("cannot load forwarder (context: %s, namespace: %s, service name: %s): %w",
				forward.KubeContextName, forward.Namespace, forward.ServiceName, err)
		}
		if forward.Hostname != "" {
			if _, exists := mappings[forward.Hostname]; exists {
				log.Printf("Service %s/%s should be proxied as %s, but domain is already mapped to another service. Skipping...\n", forward.Namespace, forward.ServiceName, forward.Hostname)
			} else {
				mappings[forward.Hostname] = forward.LocalPort
			}
		}
		err = forwarder.Start(context.TODO())
		if err != nil {
			Stop()
			return fmt.Errorf("could not start forwarder for service %s/%s: %w", forward.Namespace, forward.ServiceName, err)
		}
		forwarders[forwardIndex] = forwarder
	}

	if err = proxy.Start(
		&proxy.Config{
			SocksListenAddr: cfg.SocksListenAddr,
			HttpListenAddr:  cfg.HttpReverseListenAddr,
			Mappings:        mappings,
		},
	); err != nil {
		Stop()
		return fmt.Errorf("cannot start proxy: %w", err)
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
