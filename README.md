# k8s-client
对[client-go](https://github.com/kubernetes/client-go) 进行简单的封装,后续如果需要添加CRD的也可以直接在这里添加,并初始化.

## 支持的API
 - [x] Kubernetes原生接口
 - [x] Prometheus-operator
 - [ ] istio   

## Prometheus-operator
### 获取Prometheus资源对象
```go
func main() {
	c, err := clientcmd.BuildConfigFromKubeconfigGetter("", KubeConfigGetter)
	if err != nil {
		t.Fatal(err)
	}

	client, _ := NewForConfig(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitoringClient := client.MonitoringV1()
	proms, err := monitoringClient.Prometheuses("monitoring").List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range proms.Items {
		fmt.Printf("prometheus : %s \n", i.Name)
	}
}
```

## Kubernetes集群
### 获取deployment
```go
func main() {
	c, err := clientcmd.BuildConfigFromKubeconfigGetter("", KubeConfigGetter)
	if err != nil {
		t.Fatal(err)
	}

	client, _ := NewForConfig(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitoringClient := client.k8sCluster
	proms, err := monitoringClient.Deployment("monitoring").List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range proms.Items {
		fmt.Printf("deployment : %s \n", i.Name)
	}
}
```

### 重启deployment
kubernetes没有重启服务的功能,业务中如果更新了配置,服务本身又没有动态刷新配置的功能时就需要重启服务来获取新配置。镜像等任何配置不做修改的情况下是不会触发deployment的更新的.

使用重启注意以下事项：
* 服务是无状态的
```go
func main() {
	namespace := "dev1-xxxx-server"
	deploymentName := "dev1-app-forum-latest"
	
	c, err := clientcmd.BuildConfigFromKubeconfigGetter("", KubeConfigGetter)
	if err != nil {
		t.Fatal(err)
	}

	client, _ := NewForConfig(c)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = client.Deployment(namespace).ReDeploy(ctx, deploymentName)
	if err != nil {
		t.Fatalf("update deployment fatalf: %v", err)
	}

	t.Logf("update deployment successfully")
}
```
