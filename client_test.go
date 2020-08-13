package k8s_client

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNewForConfig(t *testing.T) {
	c, err := KubeRestConfigGetter()
	if err != nil {
		t.Fatal(err)
	}

	client, _ := NewForConfig(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitoringClient := client.k8sCluster
	proms, err := monitoringClient.Deployment("monitoring").List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range proms.Items {
		fmt.Printf("deployment : %s \n", i.Name)
	}
}
