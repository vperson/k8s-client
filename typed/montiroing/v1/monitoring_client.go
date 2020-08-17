package v1

import (
	"github.com/coreos/prometheus-operator/pkg/client/versioned"
	"k8s.io/client-go/rest"
)

type PrometheusMonitoringInterface interface {
	PrometheusGetter
	PrometheusRulesGetter
	//ServiceMonitorsGetter
}

type PrometheusMonitoring struct {
	client *versioned.Clientset
}

func NewForConfig(c *rest.Config) (*PrometheusMonitoring, error) {
	config := c
	client, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &PrometheusMonitoring{
		client: client,
	}, nil

}

func (c *PrometheusMonitoring) Prometheuses(namespace string) PrometheusInterface {
	return newPrometheuses(c.client, namespace)
}

func (c *PrometheusMonitoring) PrometheusRules(namespace string) PrometheusRuleInterface {
	return newPrometheusRules(c.client, namespace)
}
