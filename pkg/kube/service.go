package kube

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service interface {
	GetPod(ctx context.Context) (*corev1.Pod, error)
}

type service struct {
	client    *client
	namespace string
	name      string
}

func (c *client) Service(namespace string, name string) Service {
	return &service{
		client:    c,
		namespace: namespace,
		name:      name,
	}
}

func (s *service) loadService(ctx context.Context) (*corev1.Service, error) {
	svc, err := s.client.k8sClientSet.CoreV1().Services(s.namespace).Get(ctx, s.name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot get service %s/%s: %w", s.namespace, s.name, err)
	}
	return svc, nil
}

func (s *service) GetPod(ctx context.Context) (*corev1.Pod, error) {
	svc, err := s.loadService(ctx)
	if err != nil {
		return nil, err
	}

	podList, err := s.client.k8sClientSet.CoreV1().Pods(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(
			&metav1.LabelSelector{
				MatchLabels: svc.Spec.Selector,
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot get pods from service's label selector: %w", err)
	}

	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodRunning {
			return &pod, nil
		}
	}

	return nil, fmt.Errorf("service %s in namespace %s has no running pods", s.name, s.namespace)
}
