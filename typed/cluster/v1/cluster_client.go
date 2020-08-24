package v1

import (
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"path/filepath"
)

type ClusterInterface interface {
	DeploymentGetter
	PodsGetter
}

type Cluster struct {
	client     *kubernetes.Clientset
	restConfig *rest.Config
}

func NewForConfig(c *rest.Config) (*Cluster, error) {
	config := c

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Cluster{
		client:     client,
		restConfig: c,
	}, nil
}

func (c *Cluster) Deployment(namespace string) DeploymentInterface {
	if c == nil {
		return nil
	}

	return newDeployment(c.client, namespace)
}

func (c *Cluster) Pods(namespace string) PodsInterface {
	return newPods(c.client, namespace, c.restConfig)
}

func (c *Cluster) ConfigMap(namespace string) ConfigMapInterface {
	return newConfigMap(c.client, namespace)
}

func (c *Cluster) HorizontalPodAutoScalers(namespace string) HorizontalPodAutoScalersInterface {
	return newHorizontalPodAutoScaler(c.client, namespace)
}

func KubeConfigGetter() (*clientcmdapi.Config, error) {
	var configPath string

	if h := homeDir(); h != "" {
		configPath = filepath.Join(h, ".kube", "config")
	} else {
		return nil, fmt.Errorf("unknown kube config")
	}

	c, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	return clientcmd.Load(c)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	return os.Getenv("USERPROFILE") // windows
}
