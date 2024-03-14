package kube

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type Forwarder interface {
	Start(ctx context.Context) error
	Stop()
}

type forwarder struct {
	client   *client
	service  Service
	options  ForwarderOptions
	stopChan chan (struct{})
}

func (c *client) Forwarder(service Service, options ForwarderOptions) Forwarder {
	return &forwarder{
		client:   c,
		service:  service,
		options:  options,
		stopChan: make(chan struct{}, 1),
	}
}

type ForwarderOptions struct {
	PodPort   int
	LocalPort int
}

func (f *forwarder) Start(ctx context.Context) error {

	ports := []string{fmt.Sprintf("%d:%d", f.options.LocalPort, f.options.PodPort)}

	for {
		readyChan := make(chan struct{}, 1)
		errChan := make(chan error, 1)

		pod, err := f.service.GetPod(ctx)
		if err != nil {
			return err
		}

		dialer, err := f.dialer(pod)
		if err != nil {
			return err
		}

		discard := io.Discard
		portForward, err := portforward.New(dialer, ports, f.stopChan, readyChan, discard, discard)
		if err != nil {
			return err
		}

		go func() {
			err := portForward.ForwardPorts()
			errChan <- err
		}()

		select {
		case <-ctx.Done():
			log.Println("Stopping port-forward due to canceled context")
			f.Stop()
			return nil
		case <-readyChan:
			log.Printf("%s:%d forwarded to local port %d\n", pod.Name, f.options.PodPort, f.options.LocalPort)
		case err = <-errChan:
			if !errors.Is(err, portforward.ErrLostConnectionToPod) {
				return fmt.Errorf("unexpected error terminated the connection to the pod: %w", err)
			}
			log.Println("Lost connection to the pod, reconnecting")
		}

		select {
		case <-f.stopChan:
			log.Println("Received stop signal")
			return nil
		case err = <-errChan:
			if !errors.Is(err, portforward.ErrLostConnectionToPod) {
				return fmt.Errorf("unexpected error terminated the connection to the pod: %w", err)
			}
			log.Println("Lost connection to the pod, reconnecting")
		}
	}
}

func (f *forwarder) Stop() {
	f.stopChan <- struct{}{}
}

func (f *forwarder) dialer(pod *corev1.Pod) (httpstream.Dialer, error) {
	url := f.client.k8sClientSet.CoreV1().RESTClient().Post().Resource("pods").Namespace(pod.Namespace).Name(pod.Name).SubResource("portforward").URL()

	transport, upgrader, err := spdy.RoundTripperFor(f.client.restConfig)
	if err != nil {
		return nil, err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", url)
	return dialer, nil
}
