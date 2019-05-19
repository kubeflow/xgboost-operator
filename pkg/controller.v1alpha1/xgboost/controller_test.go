package xgboost

import (
	kubebatchclient "github.com/kubernetes-sigs/kube-batch/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/controller"

	control "github.com/kubeflow/common/job_controller"
	"github.com/kubeflow/xgboost-operator/cmd/xgboost-operator.v1alpha1/app/options"
	jobclientset "github.com/kubeflow/xgboost-operator/pkg/client/clientset/versioned"
	jobinformers "github.com/kubeflow/xgboost-operator/pkg/client/informers/externalversions"
)

func newXGBoostController(
	config *rest.Config,
	kubeClientSet kubeclientset.Interface,
	kubeBatchClientSet kubebatchclient.Interface,
	jobClientSet jobclientset.Interface,
	resyncPeriod controller.ResyncPeriodFunc,
	option options.ServerOption,
) (
	*XGBoostController,
	kubeinformers.SharedInformerFactory, jobinformers.SharedInformerFactory,
) {
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClientSet, resyncPeriod())
	jobInformerFactory := jobinformers.NewSharedInformerFactory(jobClientSet, resyncPeriod())

	jobInformer := NewUnstructuredXGBoostJobInformer(config, metav1.NamespaceAll)

	ctr := NewXGBoostController(jobInformer, kubeClientSet, kubeBatchClientSet, jobClientSet, kubeInformerFactory, option)
	ctr.PodControl = &controller.FakePodControl{}
	ctr.ServiceControl = &control.FakeServiceControl{}
	return ctr, kubeInformerFactory, jobInformerFactory
}