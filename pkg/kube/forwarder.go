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
	context  Context
	service  Service
	options  ForwarderOptions
	ctx      context.Context
	running  bool
	stopChan chan (struct{})
}

func (c *client) Forwarder(ctx Context, service Service, options ForwarderOptions) (Forwarder, error) {
	return &forwarder{
		client:   c,
		context:  ctx,
		service:  service,
		options:  options,
		running:  false,
		stopChan: make(chan struct{}, 1),
	}, nil
}

type ForwarderOptions struct {
	PodPort   int
	LocalPort int
}

func (f *forwarder) Start(ctx context.Context) error {

	if f.running {
		return fmt.Errorf("forwarder is already running")
	}

	ports := []string{fmt.Sprintf("%d:%d", f.options.LocalPort, f.options.PodPort)}

	startChan := make(chan error, 1)

	go func() {
		for {
			readyChan := make(chan struct{}, 1)
			errChan := make(chan error, 1)

			pod, err := f.service.GetPod(ctx)
			if err != nil {
				startChan <- err
				return
			}

			dialer, err := f.dialer(pod)
			if err != nil {
				startChan <- err
				return
			}

			discard := io.Discard
			portForward, err := portforward.New(dialer, ports, f.stopChan, readyChan, discard, discard)
			if err != nil {
				startChan <- err
				return
			}

			go func() {
				f.running = true
				err := portForward.ForwardPorts()
				errChan <- err
			}()

			err = nil
			select { // needed to first receive from readyChan
			case <-ctx.Done():
				log.Println("Stopping port-forward due to canceled context")
				f.Stop()
				startChan <- fmt.Errorf("cannot start forwarder due to canceled context")
				return
			case <-readyChan:
				log.Printf("%s:%d forwarded to local port %d\n", pod.Name, f.options.PodPort, f.options.LocalPort)
				startChan <- nil // happy case, no error
			case err = <-errChan:
				f.running = false
				if !errors.Is(err, portforward.ErrLostConnectionToPod) {
					startChan <- err
					return
				}
				panic(fmt.Errorf("got lost connection error from non-running forwarder: %w", err))
			}

			select {
			case <-ctx.Done():
				log.Println("Stopping port-forward due to canceled context")
				f.Stop()
				return
			case <-f.stopChan:
				log.Println("Received stop signal")
				return
			case err = <-errChan:
				f.running = false
				if !errors.Is(err, portforward.ErrLostConnectionToPod) {
					log.Println(fmt.Errorf("unexpected error terminated the connection to the pod: %w", err))
					return
				}
				log.Println("Lost connection to the pod, reconnecting")
			}
		}
	}()

	return <-startChan
}

func (f *forwarder) Stop() {
	if f.running {
		f.stopChan <- struct{}{}
	}

}

func (f *forwarder) dialer(pod *corev1.Pod) (httpstream.Dialer, error) {
	kubeClient, err := f.client.getKubeClient(f.context.Name)
	if err != nil {
		return nil, fmt.Errorf("cannot get kubernetes client for context '%s': %w", f.context.Name, err)
	}
	restClient, err := kubeClient.clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot get rest client for context '%s': %w", f.context.Name, err)
	}
	url := kubeClient.clientset.CoreV1().RESTClient().Post().Resource("pods").Namespace(pod.Namespace).Name(pod.Name).SubResource("portforward").URL()

	transport, upgrader, err := spdy.RoundTripperFor(restClient)
	if err != nil {
		return nil, err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", url)
	return dialer, nil
}
