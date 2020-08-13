package v1

import (
	"context"
	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/coreos/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type PrometheusGetter interface {
	Prometheuses(namespace string) PrometheusInterface
}

// PrometheusInterface has methods to work with Prometheus resources.
type PrometheusInterface interface {
	Create(ctx context.Context, prometheus *v1.Prometheus, opts metav1.CreateOptions) (*v1.Prometheus, error)
	Update(ctx context.Context, prometheus *v1.Prometheus, opts metav1.UpdateOptions) (*v1.Prometheus, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Prometheus, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.PrometheusList, error)
	Watch(ctx context.Context) (watch.Interface, error)
}

type prometheuses struct {
	client *versioned.Clientset
	ns     string
}

func newPrometheuses(c *versioned.Clientset, namespace string) *prometheuses {
	return &prometheuses{
		client: c,
		ns:     namespace,
	}
}

func (p *prometheuses) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Prometheus, error) {
	return p.client.MonitoringV1().
		Prometheuses(p.ns).
		Get(ctx, name, opts)
}

func (p *prometheuses) Update(ctx context.Context, prometheus *v1.Prometheus, opts metav1.UpdateOptions) (*v1.Prometheus, error) {
	return p.client.MonitoringV1().
		Prometheuses(p.ns).
		Update(ctx, prometheus, opts)
}

func (p *prometheuses) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return p.client.MonitoringV1().
		Prometheuses(p.ns).
		Delete(ctx, name, opts)
}

func (p *prometheuses) Create(ctx context.Context, prometheus *v1.Prometheus, opts metav1.CreateOptions) (*v1.Prometheus, error) {
	return p.client.MonitoringV1().
		Prometheuses(p.ns).
		Create(ctx, prometheus, opts)
}

func (p *prometheuses) List(ctx context.Context, opts metav1.ListOptions) (*v1.PrometheusList, error) {
	return p.client.MonitoringV1().
		Prometheuses(p.ns).
		List(ctx, opts)
}

func (p *prometheuses) Watch(ctx context.Context) (watch.Interface, error) {
	var opts metav1.ListOptions
	return p.client.MonitoringV1().
		Prometheuses(p.ns).
		Watch(ctx, opts)
}
