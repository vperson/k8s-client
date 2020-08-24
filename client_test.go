package k8s_client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"testing"
	"time"
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

func TestClientSet_Kubernetes(t *testing.T) {
	c, err := KubeRestConfigGetter()
	if err != nil {
		t.Fatal(err)
	}

	client, _ := NewForConfig(c)

	client.Kubernetes()
}

func TestKubeRestConfigGetter(t *testing.T) {

	namespace := "dev1-xiaomai-server"
	podName := "dev1-app-market-latest-5d9b6f84fb-g2pl6"
	container := "app-market"

	c, err := KubeRestConfigGetter()
	if err != nil {
		t.Fatal(err)
	}

	client, err := NewForConfig(c)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clusterClient := client.Kubernetes()
	fileByte, _ := ioutil.ReadFile("D:/Users/fonzie/go/src/xiaomai-sentry/script/dump.sh")
	file := strings.ReplaceAll(string(fileByte), "\r", "")
	_, err = clusterClient.Pods(namespace).CopyToPod(ctx, podName, container, bytes.NewReader([]byte(file)), "/usr/local/bin/dump.sh")
	if err != nil {
		t.Fatal(err)
	}

	t1 := time.Now().Format("2006-01-02_15-04-05")
	commands := []string{"/bin/bash", "-c"}
	commands = append(commands, fmt.Sprintf("/bin/bash /usr/local/bin/dump.sh %s %s", podName, t1))
	t.Log(commands)
	var stdout bytes.Buffer

	_, err = clusterClient.Pods(namespace).Exec(ctx, podName, container, commands, nil, &stdout)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(stdout.String())

}
