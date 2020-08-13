package v1

import (
	"context"
	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/coreos/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type ServiceMonitorsGetter interface {
	ServiceMonitors(namespace string) ServiceMonitorInterface
}

type ServiceMonitorInterface interface {
	Create(ctx context.Context, serviceMonitor *v1.ServiceMonitor, opts metav1.CreateOptions) (*v1.ServiceMonitor, error)
	Update(ctx context.Context, serviceMonitor *v1.ServiceMonitor, opts metav1.UpdateOptions) (*v1.ServiceMonitor, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	//DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ServiceMonitor, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ServiceMonitorList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	//Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ServiceMonitor, err error)
}

type serviceMonitors struct {
	client *versioned.Clientset
	ns     string
}

func newServiceMonitors(c *versioned.Clientset, namespace string) *serviceMonitors {
	return &serviceMonitors{
		client: c,
		ns:     namespace,
	}
}

func (s *serviceMonitors) Create(ctx context.Context, serviceMonitor *v1.ServiceMonitor, opts metav1.CreateOptions) (*v1.ServiceMonitor, error) {
	return s.client.MonitoringV1().
		ServiceMonitors(s.ns).
		Create(ctx, serviceMonitor, opts)
}

func (s *serviceMonitors) Update(ctx context.Context, serviceMonitor *v1.ServiceMonitor, opts metav1.UpdateOptions) (*v1.ServiceMonitor, error) {
	return s.client.MonitoringV1().
		ServiceMonitors(s.ns).
		Update(ctx, serviceMonitor, opts)
}

func (s *serviceMonitors) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return s.client.MonitoringV1().
		ServiceMonitors(s.ns).
		Delete(ctx, name, opts)
}

func (s *serviceMonitors) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ServiceMonitor, error) {
	return s.client.MonitoringV1().
		ServiceMonitors(s.ns).
		Get(ctx, name, opts)
}

func (s *serviceMonitors) List(ctx context.Context, opts metav1.ListOptions) (*v1.ServiceMonitorList, error) {
	return s.client.MonitoringV1().
		ServiceMonitors(s.ns).
		List(ctx, opts)
}

func (s *serviceMonitors) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return s.client.MonitoringV1().
		ServiceMonitors(s.ns).
		Watch(ctx, opts)
}
