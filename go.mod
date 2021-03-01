module github.com/kubeflow/xgboost-operator

go 1.12

require (
	cloud.google.com/go v0.39.0 // indirect
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/kubeflow/common v0.3.1
	github.com/sirupsen/logrus v1.4.2
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	google.golang.org/appengine v1.6.0 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	k8s.io/api v0.16.9
	k8s.io/apimachinery v0.16.9
	k8s.io/client-go v0.16.9
	sigs.k8s.io/controller-runtime v0.4.0
	volcano.sh/volcano v0.4.0
)

replace (
	k8s.io/api => k8s.io/api v0.16.9
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.9
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.10-beta.0
	k8s.io/apiserver => k8s.io/apiserver v0.16.9
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.9
	k8s.io/client-go => k8s.io/client-go v0.16.9
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.9
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.9
	k8s.io/code-generator => k8s.io/code-generator v0.16.10-beta.0
	k8s.io/component-base => k8s.io/component-base v0.16.9
	k8s.io/cri-api => k8s.io/cri-api v0.16.10-beta.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.9
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.9
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.9
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.9
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.9
	k8s.io/kubectl => k8s.io/kubectl v0.16.9
	k8s.io/kubelet => k8s.io/kubelet v0.16.9
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.9
	k8s.io/metrics => k8s.io/metrics v0.16.9
	k8s.io/node-api => k8s.io/node-api v0.16.9
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.9
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.16.9
	k8s.io/sample-controller => k8s.io/sample-controller v0.16.9
)
