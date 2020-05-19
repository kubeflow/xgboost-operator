module github.com/kubeflow/xgboost-operator

go 1.12

require (
	cloud.google.com/go v0.39.0 // indirect
	github.com/Azure/go-autorest v11.1.2+incompatible // indirect
	github.com/appscode/jsonpatch v0.0.0-20190108182946-7c0e3b262f30 // indirect
	github.com/codedellemc/goscaleio v0.0.0-20170830184815-20e2ce2cf885 // indirect
	github.com/d2g/dhcp4 v0.0.0-20170904100407-a1d1b6c41b1c // indirect
	github.com/d2g/dhcp4client v0.0.0-20170829104524-6e570ed0a266 // indirect
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/gobuffalo/envy v1.7.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/jteeuwen/go-bindata v0.0.0-20151023091102-a0ff2567cfb7 // indirect
	github.com/kardianos/osext v0.0.0-20150410034420-8fef92e41e22 // indirect
	github.com/kr/fs v0.0.0-20131111012553-2788f0dbd169 // indirect
	github.com/kubeflow/common v0.3.1
	github.com/kubernetes-sigs/kube-batch v0.4.2
	github.com/markbates/inflect v1.0.4 // indirect
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a // indirect
	github.com/mholt/caddy v0.0.0-20180213163048-2de495001514 // indirect
	github.com/petar/GoLLRB v0.0.0-20190514000832-33fb24c13b99 // indirect
	github.com/pkg/sftp v0.0.0-20160930220758-4d0e916071f6 // indirect
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v0.0.0-20151028001915-10ef21a441db // indirect
	github.com/sigma/go-inotify v0.0.0-20181102212354-c87b6cf5033d // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/vmware/photon-controller-go-sdk v0.0.0-20170310013346-4a435daef6cc // indirect
	github.com/xanzy/go-cloudstack v0.0.0-20160728180336-1e2cbf647e57 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	google.golang.org/appengine v1.6.0 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	k8s.io/api v0.16.9
	k8s.io/apimachinery v0.16.9
	k8s.io/client-go v0.16.9
	k8s.io/kubernetes v1.16.9
	sigs.k8s.io/controller-runtime v0.4.0
	sigs.k8s.io/controller-tools v0.1.8 // indirect
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
