module github.com/kubeflow/xgboost-operator

go 1.12

require (
	cloud.google.com/go v0.39.0 // indirect
	github.com/Azure/azure-sdk-for-go v32.5.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.9.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.2.0 // indirect
	github.com/GoogleCloudPlatform/k8s-cloud-provider v0.0.0-20190822182118-27a4ced34534 // indirect
	github.com/appscode/jsonpatch v0.0.0-20190108182946-7c0e3b262f30 // indirect
	github.com/bazelbuild/bazel-gazelle v0.19.1-0.20191105222053-70208cbdc798 // indirect
	github.com/checkpoint-restore/go-criu v0.0.0-20190109184317-bdb7599cd87b // indirect
	github.com/containernetworking/cni v0.7.1 // indirect
	github.com/coredns/corefile-migration v1.0.2 // indirect
	github.com/coreos/etcd v3.3.17+incompatible // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/go-bindata/go-bindata v3.1.1+incompatible // indirect
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/go-openapi/validate v0.19.2 // indirect
	github.com/gobuffalo/envy v1.7.0 // indirect
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/google/cadvisor v0.34.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/heketi/heketi v9.0.1-0.20190917153846-c2e2a4ab7ab9+incompatible // indirect
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/kubeflow/common v0.0.0-20190619012831-09e1ac17011c
	github.com/kubernetes-sigs/kube-batch v0.4.2
	github.com/libopenstorage/openstorage v1.0.0 // indirect
	github.com/markbates/inflect v1.0.4 // indirect
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a // indirect
	github.com/mattn/go-shellwords v1.0.5 // indirect
	github.com/miekg/dns v1.1.4 // indirect
	github.com/mistifyio/go-zfs v2.1.1+incompatible // indirect
	github.com/mvdan/xurls v1.1.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v1.0.0-rc2.0.20190611121236-6cc515888830 // indirect
	github.com/opencontainers/selinux v1.2.2 // indirect
	github.com/petar/GoLLRB v0.0.0-20190514000832-33fb24c13b99 // indirect
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/robfig/cron v1.1.0 // indirect
	github.com/seccomp/libseccomp-golang v0.9.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/thecodeteam/goscaleio v0.1.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975 // indirect
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9 // indirect
	google.golang.org/api v0.6.1-0.20190607001116-5213b8090861 // indirect
	google.golang.org/appengine v1.6.0 // indirect
	google.golang.org/grpc v1.23.0 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	gotest.tools/gotestsum v0.3.5 // indirect
	honnef.co/go/tools v0.0.1-2019.2.2 // indirect
	k8s.io/api v0.16.9
	k8s.io/apimachinery v0.16.9
	k8s.io/client-go v0.16.9
	k8s.io/gengo v0.0.0-20190822140433-26a664648505 // indirect
	k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf // indirect
	k8s.io/kubectl v0.0.0 // indirect
	k8s.io/kubernetes v1.15.10
	k8s.io/repo-infra v0.0.1-alpha.1 // indirect
	k8s.io/utils v0.0.0-20190801114015-581e00157fb1 // indirect
	sigs.k8s.io/controller-runtime v0.3.0
	sigs.k8s.io/controller-tools v0.1.8 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.15.10
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.15.10
	k8s.io/apimachinery => k8s.io/apimachinery v0.15.12-beta.0
	k8s.io/apiserver => k8s.io/apiserver v0.15.10
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.15.10
	k8s.io/client-go => k8s.io/client-go v0.15.10
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.15.10
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.15.10
	k8s.io/code-generator => k8s.io/code-generator v0.15.13-beta.0
	k8s.io/component-base => k8s.io/component-base v0.15.10
	k8s.io/cri-api => k8s.io/cri-api v0.15.13-beta.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.15.10
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.15.10
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.15.10
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.15.10
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.15.10
	k8s.io/kubectl => k8s.io/kubectl v0.15.13-beta.0
	k8s.io/kubelet => k8s.io/kubelet v0.15.10
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.15.10
	k8s.io/metrics => k8s.io/metrics v0.15.10
	k8s.io/node-api => k8s.io/node-api v0.15.10
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.15.10
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.15.10
	k8s.io/sample-controller => k8s.io/sample-controller v0.15.10
)
