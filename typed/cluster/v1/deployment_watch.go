package v1

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"time"
)

type DeploymentController struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
}

func NewDeploymentController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller) *DeploymentController {
	return &DeploymentController{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
	}
}

func (c *DeploymentController) runWorker() {
	for c.processNextItem() {

	}
}

func (c *DeploymentController) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	defer c.queue.Done(key)

	err := c.syncToStdout(key.(string))
	c.handleErr(err, key)
	return true
}

func (c *DeploymentController) syncToStdout(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		klog.Errorf("fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		klog.Infof("deployment %s does not exist anymore", key)
	} else {
		klog.Infof("sync/add/update for deployment %s", obj.(*v1.Deployment).GetName())
	}

	return nil
}

func (c *DeploymentController) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("error syncing deployment %v: %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	klog.Infof("dropping deployment %q out of the queue: %v", key, err)
}

func (c *DeploymentController) Run(threadNum int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	defer c.queue.ShutDown()

	klog.Info("start deployment controller")

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("time out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadNum; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh
	klog.Info("stopping deployment controller")
}
