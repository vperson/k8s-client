package v1

import (
	"context"
	v2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type HorizontalPodAutoScalersGetter interface {
	Deployment(namespace string) HorizontalPodAutoScalersInterface
}

type HorizontalPodAutoScalersInterface interface {
	Create(ctx context.Context, horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler, opts metav1.CreateOptions) (*v2beta2.HorizontalPodAutoscaler, error)
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v2beta2.HorizontalPodAutoscaler, error)
	Update(ctx context.Context, horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler, opts metav1.UpdateOptions) (*v2beta2.HorizontalPodAutoscaler, error)
}

type horizontalPodAutoScaler struct {
	client *kubernetes.Clientset
	ns     string
}

func newHorizontalPodAutoScaler(c *kubernetes.Clientset, namespace string) *horizontalPodAutoScaler {
	return &horizontalPodAutoScaler{
		client: c,
		ns:     namespace,
	}
}

func (h *horizontalPodAutoScaler) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v2beta2.HorizontalPodAutoscaler, error) {
	return h.client.AutoscalingV2beta2().
		HorizontalPodAutoscalers(h.ns).
		Get(ctx, name, opts)
}

func (h *horizontalPodAutoScaler) Update(ctx context.Context, horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler, opts metav1.UpdateOptions) (*v2beta2.HorizontalPodAutoscaler, error) {
	return h.client.AutoscalingV2beta2().
		HorizontalPodAutoscalers(h.ns).
		Update(ctx, horizontalPodAutoscaler, opts)
}

func (h *horizontalPodAutoScaler) Create(ctx context.Context, horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler, opts metav1.CreateOptions) (*v2beta2.HorizontalPodAutoscaler, error) {
	return h.client.AutoscalingV2beta2().
		HorizontalPodAutoscalers(h.ns).
		Create(ctx, horizontalPodAutoscaler, opts)
}
