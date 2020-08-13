module github.com/vperson/k8s-client

go 1.14

require (
	github.com/coreos/prometheus-operator v0.41.0
	github.com/imdario/mergo v0.3.9 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace (
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // pinned to release-branch.go1.13
	k8s.io/api => k8s.io/api v0.0.0-20200214081623-ecbd4af0fc33
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20200214081019-7490b3ed6e92
	k8s.io/client-go => k8s.io/client-go v0.0.0-20200214082307-e38a84523341
)
