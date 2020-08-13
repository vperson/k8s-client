package v1

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapGetter interface {
	ConfigMap(namespace string)
}

type ConfigMapInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ConfigMap, error)
	Create(ctx context.Context, configMapData *v1.ConfigMap, opts metav1.CreateOptions) (*v1.ConfigMap, error)
	Update(ctx context.Context, configMapData *v1.ConfigMap, opts metav1.UpdateOptions) (*v1.ConfigMap, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type configMap struct {
	client *kubernetes.Clientset
	ns     string
}

func newConfigMap(c *kubernetes.Clientset, ns string) *configMap {
	return &configMap{
		client: c,
		ns:     ns,
	}
}

func (c *configMap) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ConfigMap, error) {
	return c.client.CoreV1().
		ConfigMaps(c.ns).
		Get(ctx, name, opts)
}

func (c *configMap) Create(ctx context.Context, configMapData *v1.ConfigMap, opts metav1.CreateOptions) (*v1.ConfigMap, error) {
	return c.client.CoreV1().
		ConfigMaps(c.ns).
		Create(ctx, configMapData, opts)
}

func (c *configMap) Update(ctx context.Context, configMapData *v1.ConfigMap, opts metav1.UpdateOptions) (*v1.ConfigMap, error) {
	return c.client.CoreV1().
		ConfigMaps(c.ns).
		Update(ctx, configMapData, opts)
}

func (c *configMap) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.CoreV1().
		ConfigMaps(c.ns).
		Delete(ctx, name, &opts)
}

func (c *configMap) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.CoreV1().
		ConfigMaps(c.ns).
		Watch(ctx, opts)
}
