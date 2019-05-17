package xgboost

import (
	"encoding/json"
	"k8s.io/client-go/tools/cache"
	"testing"
	"time"
	"fmt"
	"strings"

	"github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1beta1"
	kubebatchclient "github.com/kubernetes-sigs/kube-batch/pkg/client/clientset/versioned"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/controller"

	"github.com/kubeflow/xgboost-operator/cmd/xgboost-operator.v1alpha1/app/options"
	v1alpha1 "github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
	jobclientset "github.com/kubeflow/xgboost-operator/pkg/client/clientset/versioned"
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1beta1"
)

var (
	AlwaysReady = func() bool { return true }

	GroupName = v1beta1.GroupName
)

const (
	SleepInterval = 500 * time.Millisecond
	ThreadCount   = 1
	LabelGroupName      = "group-name"
	LabelPyTorchJobName = "pytorch-job-name"
	LabelMaster        = "master"
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
			GroupVersion: &v1beta1.SchemeGroupVersion,
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
		ctr.Run(ThreadCount, stopCh)
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
		t.Errorf("Failed to enqueue the PyTorchJob %s: expected %s, got %s", job.Name, GetKey(job, t), key)
	}
	close(stopCh)
}

func NewXGBoostJobWithMaster(worker int) *v1alpha1.XGBoostJob {
	job := newXGboostJob(worker)
	job.Spec.PyTorchReplicaSpecs[v1beta1.PyTorchReplicaTypeMaster] = &common.ReplicaSpec{
		Template: NewPyTorchReplicaSpecTemplate(),
	}
	return job
}

// ConvertXGBoostJob uses JSON to convert PyTorchJob to Unstructured.
func ConvertXGBoostJobToUnstructured(job *v1alpha1.XGBoostJob) (*unstructured.Unstructured, error) {
	var unstructured unstructured.Unstructured
	b, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &unstructured); err != nil {
		return nil, err
	}
	return &unstructured, nil
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

func GenLabels(jobName string) map[string]string {
	return map[string]string{
		LabelGroupName:      GroupName,
		LabelPyTorchJobName: strings.Replace(jobName, "/", "-", -1),
	}
}

func GetKey(job *v1alpha1.XGBoostJob, t *testing.T) string {
	key, err := KeyFunc(job)
	if err != nil {
		t.Errorf("Unexpected error getting key for job %v: %v", job.Name, err)
		return ""
	}
	return key
}