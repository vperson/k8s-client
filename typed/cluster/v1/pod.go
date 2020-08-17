package v1

import (
	"bytes"
	"context"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type PodsGetter interface {
	Pods(namespace string) PodsInterface
}

type PodsInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Pod, error)
	Update(ctx context.Context, pod *v1.Pod, opts metav1.UpdateOptions) (*v1.Pod, error)
	Create(ctx context.Context, pod *v1.Pod, opts metav1.CreateOptions) (*v1.Pod, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	List(ctx context.Context, opts metav1.ListOptions) (*v1.PodList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Exec(ctx context.Context, podName, containerName string, command []string, stdin io.Reader, stdout io.Writer) ([]byte, error)
	CopyToPod(ctx context.Context, podName, containerName string, sourceFile io.Reader, targetFile string) ([]byte, error)
}

type pods struct {
	client     *kubernetes.Clientset
	ns         string
	restConfig *rest.Config
}

func newPods(c *kubernetes.Clientset, namespace string, config *rest.Config) *pods {
	return &pods{
		client:     c,
		ns:         namespace,
		restConfig: config,
	}
}

func (p *pods) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Pod, error) {
	return p.client.CoreV1().
		Pods(p.ns).
		Get(ctx, name, opts)
}

func (p *pods) Update(ctx context.Context, pod *v1.Pod, opts metav1.UpdateOptions) (*v1.Pod, error) {
	return p.client.CoreV1().
		Pods(p.ns).
		Update(ctx, pod, opts)
}

func (p *pods) Create(ctx context.Context, pod *v1.Pod, opts metav1.CreateOptions) (*v1.Pod, error) {
	return p.client.CoreV1().
		Pods(p.ns).
		Create(ctx, pod, opts)
}

func (p *pods) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return p.client.CoreV1().
		Pods(p.ns).
		Delete(ctx, name, &opts)
}

func (p *pods) List(ctx context.Context, opts metav1.ListOptions) (*v1.PodList, error) {
	return p.client.CoreV1().
		Pods(p.ns).
		List(ctx, opts)
}

func (p *pods) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return p.client.CoreV1().
		Pods(p.ns).
		Watch(ctx, opts)
}

func (p *pods) Exec(ctx context.Context, podName, containerName string, command []string, stdin io.Reader, stdout io.Writer) ([]byte, error) {
	_, err := p.Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get pod %s err : %v", podName, err)
	}
	req := p.client.CoreV1().RESTClient().Post().Resource("pods").Namespace(p.ns).Name(podName).SubResource("exec").VersionedParams(&v1.PodExecOptions{
		Stdin:     stdin != nil,
		Stdout:    stdout != nil,
		Stderr:    true,
		TTY:       false,
		Container: containerName,
		Command:   command,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.restConfig, "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("error while creating executor: %v", err)
	}

	var stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: &stderr,
		Tty:    false,
	})

	if err != nil {
		return stderr.Bytes(), fmt.Errorf("error in Stream: %v", err)
	}

	return stderr.Bytes(), nil
}

func (p *pods) CopyToPod(ctx context.Context, podName, containerName string, sourceFile io.Reader, targetFile string) ([]byte, error) {
	commands := []string{"/bin/bash", "-c"}

	commands = append(commands, fmt.Sprintf("cp -f /dev/stdin %[1]s", targetFile))
	var stdout bytes.Buffer
	stderr, err := p.Exec(ctx, podName, containerName, commands, sourceFile, &stdout)
	if err != nil {
		return nil, err
	}

	return stderr, nil
}
