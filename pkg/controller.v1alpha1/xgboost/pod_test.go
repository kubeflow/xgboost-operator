package xgboost

import (

	"testing"
	"fmt"

	kubebatchclient "github.com/kubernetes-sigs/kube-batch/pkg/client/clientset/versioned"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/controller"

	"github.com/kubeflow/xgboost-operator/cmd/xgboost-operator.v1alpha1/app/options"
	v1alpha1 "github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
	jobclientset "github.com/kubeflow/xgboost-operator/pkg/client/clientset/versioned"
)


func TestAddPod(t *testing.T) {
	// Prepare the clientset and controller for the test.
	kubeClientSet := kubeclientset.NewForConfigOrDie(&rest.Config{
		Host: "",
		ContentConfig: rest.ContentConfig{
			GroupVersion: &v1.SchemeGroupVersion,
		},
	},
	)

	// Prepare the kube-batch clientset and controller for the test.
	kubeBatchClientSet := kubebatchclient.NewForConfigOrDie(&rest.Config{
		Host: "",
		ContentConfig: rest.ContentConfig{
			GroupVersion: &v1.SchemeGroupVersion,
		},
	},
	)

	config := &rest.Config{
		Host: "",
		ContentConfig: rest.ContentConfig{
			GroupVersion: &v1alpha1.SchemeGroupVersion,
		},
	}
	jobClientSet := jobclientset.NewForConfigOrDie(config)

	ctr, _, _ :=  newXGBoostController(config, kubeClientSet, kubeBatchClientSet, jobClientSet, controller.NoResyncPeriodFunc, options.ServerOption{})
	ctr.jobInformerSynced = AlwaysReady
	ctr.PodInformerSynced = AlwaysReady
	ctr.ServiceInformerSynced = AlwaysReady
	jobIndexer := ctr.jobInformer.GetIndexer()

	stopCh := make(chan struct{})
	run := func(<-chan struct{}) {
		err := ctr.Run(ThreadCount, stopCh)
		if err != nil {
			t.Errorf("Failed to run thread count: %v", err)
		}

	}
	go run(stopCh)

	var key string
	syncChan := make(chan string)
	ctr.syncHandler = func(jobKey string) (bool, error) {
		key = jobKey
		<-syncChan
		return true, nil
	}

	job := NewXGBoostJobWithMaster(1)
	unstructured, err := ConvertXGBoostJobToUnstructured(job)
	if err != nil {
		t.Errorf("Failed to convert the job to Unstructured: %v", err)
	}

	if err := jobIndexer.Add(unstructured); err != nil {
		t.Errorf("Failed to add job to jobIndexer: %v", err)
	}
	pod := NewPod(job, LabelMaster, 0, t)
	ctr.AddPod(pod)

	syncChan <- "sync"
	if key != GetKey(job, t) {
		t.Errorf("Failed to enqueue the XGBoostJob %s: expected %s, got %s", job.Name, GetKey(job, t), key)
	}
	close(stopCh)
}


func NewPod(job *v1alpha1.XGBoostJob, typ string, index int, t *testing.T) *v1.Pod {
	pod := NewBasePod(fmt.Sprintf("%s-%d", typ, index), job, t)
	pod.Labels[replicaTypeLabel] = typ
	pod.Labels[replicaIndexLabel] = fmt.Sprintf("%d", index)
	return pod
}

func NewBasePod(name string, job *v1alpha1.XGBoostJob, t *testing.T) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Labels:          GenLabels(job.Name),
			Namespace:       job.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(job, controllerKind)},
		},
	}
}

