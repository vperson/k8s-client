package v1

import (
	"context"
	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/coreos/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type PrometheusRulesGetter interface {
	PrometheusRules(namespace string) PrometheusRuleInterface
}

type PrometheusRuleInterface interface {
	Create(ctx context.Context, prometheusRule *v1.PrometheusRule, opts metav1.CreateOptions) (*v1.PrometheusRule, error)
	Update(ctx context.Context, prometheusRule *v1.PrometheusRule, opts metav1.UpdateOptions) (*v1.PrometheusRule, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	//DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.PrometheusRule, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.PrometheusRuleList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

// prometheusRules implements PrometheusRuleInterface
type prometheusRules struct {
	client *versioned.Clientset
	ns     string
}

// newPrometheusRules returns a PrometheusRules
func newPrometheusRules(c *versioned.Clientset, namespace string) *prometheusRules {
	return &prometheusRules{
		client: c,
		ns:     namespace,
	}
}

func (p *prometheusRules) Create(ctx context.Context, prometheusRule *v1.PrometheusRule, opts metav1.CreateOptions) (*v1.PrometheusRule, error) {
	return p.client.MonitoringV1().
		PrometheusRules(p.ns).
		Create(ctx, prometheusRule, opts)
}

func (p *prometheusRules) Update(ctx context.Context, prometheusRule *v1.PrometheusRule, opts metav1.UpdateOptions) (*v1.PrometheusRule, error) {
	return p.client.MonitoringV1().
		PrometheusRules(p.ns).
		Update(ctx, prometheusRule, opts)
}

func (p *prometheusRules) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return p.client.MonitoringV1().
		PrometheusRules(p.ns).
		Delete(ctx, name, opts)
}

func (p *prometheusRules) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.PrometheusRule, error) {
	return p.client.MonitoringV1().
		PrometheusRules(p.ns).
		Get(ctx, name, opts)
}

func (p *prometheusRules) List(ctx context.Context, opts metav1.ListOptions) (*v1.PrometheusRuleList, error) {
	return p.client.MonitoringV1().
		PrometheusRules(p.ns).
		List(ctx, opts)
}

func (p *prometheusRules) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return p.client.MonitoringV1().
		PrometheusRules(p.ns).
		Watch(ctx, opts)
}
