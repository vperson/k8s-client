package k8s_client

import (
	"encoding/json"
	"fmt"
	k8sCluster "github.com/vperson/k8s-client/typed/cluster/v1"
	monitoringV1 "github.com/vperson/k8s-client/typed/montiroing/v1"
	"io/ioutil"
	discovery "k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"k8s.io/client-go/util/flowcontrol"
	"os"
	"path/filepath"
	"runtime"
	"sigs.k8s.io/yaml"
	"time"
)

const (
	// High enough QPS to fit all expected use cases.
	defaultQPS = 1e6
	// High enough Burst to fit all expected use cases.
	defaultBurst = 1e6
	// full resyc cache resource time
	defaultResyncPeriod = 30 * time.Second
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	monitoringV1.PrometheusMonitoringInterface
	k8sCluster.ClusterInterface
}

type ClientSet struct {
	*discovery.DiscoveryClient
	monitoring *monitoringV1.PrometheusMonitoring
	k8sCluster *k8sCluster.Cluster
}

// 获取prometheus operator相关的方法
func (c *ClientSet) MonitoringV1() monitoringV1.PrometheusMonitoringInterface {
	return c.monitoring
}

// 获取Kubernetes集群相关的方法
func (c *ClientSet) Kubernetes() k8sCluster.ClusterInterface {
	return c.k8sCluster
}

// 获取动态client,可深度再定制
func (c *ClientSet) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}

	return c.DiscoveryClient
}

func NewForConfig(c *rest.Config) (*ClientSet, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}

	var (
		cs  ClientSet
		err error
	)

	cs.k8sCluster, err = k8sCluster.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	cs.monitoring, err = monitoringV1.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(c)
	if err != nil {
		return nil, err
	}

	return &cs, nil

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

func KubeRestConfigGetter() (*rest.Config, error) {
	var (
		c          *rest.Config
		err        error
		configPath string
	)
	if dockerEnvIsExist() == true {
		c, err = rest.InClusterConfig()
	} else {
		if h := homeDir(); h != "" {
			configPath = filepath.Join(h, ".kube", "config")
		} else {
			return nil, fmt.Errorf("unknown kube config")
		}

		configFile, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		jsonConfig, err := yaml.YAMLToJSON(configFile)
		if err != nil {
			return nil, fmt.Errorf("yaml umarshal kubernetes config err: %v", err)
		}
		configV1 := clientcmdapiv1.Config{}
		err = json.Unmarshal(jsonConfig, &configV1)
		if err != nil {
			return nil, fmt.Errorf("yaml umarshal kubernetes config err: %v", err)
		}

		configObject, err := clientcmdlatest.Scheme.ConvertToVersion(&configV1, clientcmdapi.SchemeGroupVersion)
		if err != nil {
			return nil, fmt.Errorf("clientcmd latest scheme conver to version error. %v", err)
		}

		configInternal := configObject.(*clientcmdapi.Config)

		c, err = clientcmd.NewDefaultClientConfig(*configInternal, &clientcmd.ConfigOverrides{
			ClusterDefaults: clientcmdapi.Cluster{Server: ""},
		}).ClientConfig()

		if err != nil {
			return nil, err
		}

		c.QPS = defaultQPS
		c.Burst = defaultBurst
	}

	return c, err
}

func dockerEnvIsExist() bool {
	if runtime.GOOS != "linux" {
		return false
	}

	dockerEnv := "/.dockerenv"
	_, err := os.Stat(dockerEnv)
	if err == nil {
		return true
	}
	return false
}
