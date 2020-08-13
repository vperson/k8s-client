package v1

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
)

func TestDeployment_ListWatch(t *testing.T) {
	c, err := clientcmd.BuildConfigFromKubeconfigGetter("", KubeConfigGetter)
	if err != nil {
		t.Fatal(err)
	}

	client, _ := NewForConfig(c)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client.Deployment("dev-xiaomai-server").ListWatch(ctx)
}

func TestDeployment_Watch(t *testing.T) {
	c, err := clientcmd.BuildConfigFromKubeconfigGetter("", KubeConfigGetter)
	if err != nil {
		t.Fatal(err)
	}

	client, _ := NewForConfig(c)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w, _ := client.Deployment("dev-xiaomai-server").Watch(ctx, metav1.ListOptions{})
	for {
		select {
		case obj := <-w.ResultChan():

			fmt.Println(obj.Type)
			fmt.Println()
		case <-ctx.Done():
			return
		}
	}
}

func TestDeployment_ReDeploy(t *testing.T) {
	namespace := "dev1-xiaomai-server"
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
