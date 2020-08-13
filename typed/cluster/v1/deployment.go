package v1

import (
	"context"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"regexp"
	"strconv"
)

type DeploymentGetter interface {
	Deployment(namespace string) DeploymentInterface
}

type DeploymentInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Deployment, error)
	Update(ctx context.Context, deployment *v1.Deployment, opts metav1.UpdateOptions) (*v1.Deployment, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	List(ctx context.Context, opts metav1.ListOptions) (*v1.DeploymentList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	ListWatch(ctx context.Context)
	ReDeploy(ctx context.Context, name string) error
}

type deployment struct {
	client *kubernetes.Clientset
	ns     string
}

func newDeployment(c *kubernetes.Clientset, namespace string) *deployment {
	return &deployment{
		client: c,
		ns:     namespace,
	}
}

// 获取deployment
func (d *deployment) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Deployment, error) {
	return d.client.AppsV1().
		Deployments(d.ns).
		Get(ctx, name, opts)
}

func (d *deployment) Update(ctx context.Context, deployment *v1.Deployment, opts metav1.UpdateOptions) (*v1.Deployment, error) {
	return d.client.AppsV1().
		Deployments(d.ns).
		Update(ctx, deployment, opts)
}

func (d *deployment) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return d.client.AppsV1().
		Deployments(d.ns).
		Delete(ctx, name, &opts)
}

func (d *deployment) List(ctx context.Context, opts metav1.ListOptions) (*v1.DeploymentList, error) {
	return d.client.AppsV1().
		Deployments(d.ns).
		List(ctx, opts)
}

func (d *deployment) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return d.client.AppsV1().
		Deployments(d.ns).
		Watch(ctx, opts)
}

func (d *deployment) ListWatch(ctx context.Context) {
	deploymentListWatcher := cache.NewListWatchFromClient(d.client.AppsV1().RESTClient(), "deployments", d.ns, fields.Everything())
	// 创建工作队列
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	indexer, informer := cache.NewIndexerInformer(deploymentListWatcher, &v1.Deployment{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	controller := NewDeploymentController(queue, indexer, informer)

	stop := make(chan struct{})
	defer close(stop)

	go controller.Run(1, stop)
	select {
	case <-ctx.Done():
		return
	}

}

// kubernetes没有重启服务的功能,业务中如果更新了配置,服务本
//身又没有动态刷新配置的功能时就需要重启服务来获取新配置
func (d *deployment) ReDeploy(ctx context.Context, name string) error {
	loopback := "127.0.0.1"
	reDeployDomainTml := "deployment-%d.redeploy.local"
	var (
		localExist    bool
		domainExist   bool
		num           int
		exitFor       bool
		err           error
		hostnameIndex int
		deployment    *v1.Deployment
	)

	deployment, err = d.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	hostAliases := deployment.Spec.Template.Spec.HostAliases
	// 检查是否有127.0.0.1的IP
	for _, host := range hostAliases {
		if host.IP == loopback {
			localExist = true
			break
		}
		localExist = false
	}

	if localExist == true {
		// 检查是否有deployment-x.redeploy.local的域名
		for _, host := range hostAliases {
			if host.IP == loopback {
				for index, d := range host.Hostnames {
					num, domainExist = regexReDeployDomain(d)
					if domainExist == true {
						hostnameIndex = index
						exitFor = true
						break
					}
				}
			}
			if exitFor == true {
				exitFor = false
				break
			}
		}

		if domainExist == true {
			// 获取域名,并解析出数字,在原有的数字上进行+1
			num += 1
		} else {
			// 在deployment的hostAliases添加域名
			num = 1
		}

		for index, host := range hostAliases {
			if host.IP == loopback {
				reDeployDomain := fmt.Sprintf(reDeployDomainTml, num)
				if domainExist == true {
					deployment.Spec.Template.Spec.HostAliases[index].Hostnames[hostnameIndex] = reDeployDomain
				} else {
					deployment.Spec.Template.Spec.HostAliases[index].Hostnames = append(deployment.Spec.Template.Spec.HostAliases[index].Hostnames, reDeployDomain)
				}

			}
		}

	} else {
		// 添加127.0.0.1的IP,并添加域名deployment-1.redeploy.local的域名
		reDeployHost := coreV1.HostAlias{
			IP: "127.0.0.1",
			Hostnames: []string{
				"deployment-1.redeploy.local",
			},
		}
		deployment.Spec.Template.Spec.HostAliases = append(deployment.Spec.Template.Spec.HostAliases, reDeployHost)
	}

	// 对已经修改完成的deployment进行更新
	_, err = d.Update(ctx, deployment, metav1.UpdateOptions{})
	return err
}

func regexReDeployDomain(domain string) (int, bool) {
	reDeployRegexp := regexp.MustCompile(`^deployment-([\d]+).redeploy.local$`)
	match := reDeployRegexp.MatchString(domain)
	if match == true {
		params := reDeployRegexp.FindStringSubmatch(domain)
		if len(params) >= 2 {
			num, err := strconv.Atoi(params[1])
			if err != nil {
				return 0, false
			}

			return num, match
		}
	}
	return 0, match
}
