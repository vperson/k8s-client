package v1

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
)

func TestPods_Exec(t *testing.T) {
	c, err := clientcmd.BuildConfigFromKubeconfigGetter("", KubeConfigGetter)
	if err != nil {
		t.Fatal(err)
	}
	client, _ := NewForConfig(c)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	namespace := "dev-xiaomai-server"
	pod := "dev-app-gateway-latest-76ddb96f4c-nbkrr"
	container := "app-gateway"

	commands := []string{"/bin/bash", "-c"}
	commands = append(commands, fmt.Sprintf("ls -l"))

	var stdout bytes.Buffer

	stderr, err := client.Pods(namespace).Exec(ctx, pod, container, commands, nil, &stdout)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(stderr)
	fmt.Println(stdout.String())

}

func TestPods_CopyToPod(t *testing.T) {
	c, err := clientcmd.BuildConfigFromKubeconfigGetter("", KubeConfigGetter)
	if err != nil {
		t.Fatal(err)
	}
	client, _ := NewForConfig(c)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	namespace := "dev-xiaomai-server"
	pod := "dev-app-gateway-latest-76ddb96f4c-nbkrr"
	container := "app-gateway"

	testFile, err := ioutil.ReadFile("D:/tmp/awesomeProject")
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := client.Pods(namespace).CopyToPod(ctx, pod, container, bytes.NewReader(testFile), "/tmp/main")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("stderr : %s", string(stderr))
}
