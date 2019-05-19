// Copyright 2019 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xgboost

import (
	"fmt"
	jobcontroller "github.com/kubeflow/common/job_controller"
	common "github.com/kubeflow/common/operator/v1"
	commonutil "github.com/kubeflow/common/util"
	"github.com/kubeflow/common/util/k8sutil"

	"github.com/kubeflow/xgboost-operator/cmd/xgboost-operator.v1alpha1/app/options"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
	jobclientset "github.com/kubeflow/xgboost-operator/pkg/client/clientset/versioned"
	jobscheme "github.com/kubeflow/xgboost-operator/pkg/client/clientset/versioned/scheme"
	jobinformersv1alpha1 "github.com/kubeflow/xgboost-operator/pkg/client/informers/externalversions/xgboost/v1alpha1"
	kubebatchclient "github.com/kubernetes-sigs/kube-batch/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"time"
)

const (
	controllerName = "xgboost-operator"

	// labels for pods and servers.
	replicaTypeLabel    = "xgboost-replica-type"
	replicaIndexLabel   = "xgboost-replica-index"
	labelXGBoostJobRole = "xgboost-job-role"
)

type XGBoostController struct {
	jobcontroller.JobController

	// jobClientSet is a clientset for CRD XGBoostJOb.
	jobClientSet jobclientset.Interface

	// To allow injection of sync functions for testing.
	syncHandler func(string) (bool, error)

	// jobInformer is a temporary field for unstructured informer support.
	jobInformer cache.SharedIndexInformer

	// jobInformerSynced returns true if the job store has been synced at least once.
	jobInformerSynced cache.InformerSynced
}


func NewXGBoostController(
	jobInformer jobinformersv1alpha1.XGBoostJobInformer,
	kubeClientSet kubeclientset.Interface,
	kubeBatchClientSet kubebatchclient.Interface,
	jobClientSet jobclientset.Interface,
	kubeInformerFactory kubeinformers.SharedInformerFactory,
	option options.ServerOption) *XGBoostController {

	jobscheme.AddToScheme(scheme.Scheme)

	log.Info("Creating XGBoostJob controller")
	// Create new XGBoostController.
	xc := &XGBoostController{
		jobClientSet: jobClientSet,
	}

	// Create base controller
	log.Info("Creating Job controller")
	jc := jobcontroller.NewJobController(xc,
		metav1.Duration{Duration: 15 * time.Second},
		option.EnableGangScheduling,
		kubeClientSet,
		kubeBatchClientSet,
		kubeInformerFactory,
		v1alpha1.Plural)

	xc.JobController = jc

	xc.syncHandler = xc.syncXGBoostJob

	jobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    xc.addXGBoostJob,
		UpdateFunc: xc.updateXGBoostJob,
		DeleteFunc: xc.enqueueXGBoostJob,
	})

	xc.jobInformer = jobInformer.Informer()
	xc.jobInformerSynced = jobInformer.Informer().HasSynced

	// Create pod informer.
	podInformer := kubeInformerFactory.Core().V1().Pods()

	// Set up an event handler for when pod resources change
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    jc.AddPod,
		UpdateFunc: jc.UpdatePod,
		DeleteFunc: jc.DeletePod,
	})

	xc.PodLister = podInformer.Lister()
	xc.PodInformerSynced = podInformer.Informer().HasSynced

	// Create service informer.
	serviceInformer := kubeInformerFactory.Core().V1().Services()

	// Set up an event handler for when service resources change.
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    jc.AddService,
		UpdateFunc: jc.UpdateService,
		DeleteFunc: jc.DeleteService,
	})

	xc.ServiceLister = serviceInformer.Lister()
	xc.ServiceInformerSynced = serviceInformer.Informer().HasSynced

	return xc
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (xc *XGBoostController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer xc.WorkQueue.ShutDown()

	// Start the informer factories to begin populating the informer caches.
	log.Info("Starting XGBoostJob controller")

	// Wait for the caches to be synced before starting workers.
	log.Info("Waiting for informer caches to sync")

	if ok := cache.WaitForCacheSync(stopCh, xc.jobInformerSynced,
		xc.PodInformerSynced, xc.ServiceInformerSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	log.Infof("Starting %v workers", threadiness)
	// Launch workers to process XGBoostJob resources.
	for i := 0; i < threadiness; i++ {
		go wait.Until(xc.runWorker, time.Second, stopCh)
	}

	log.Info("Started workers")
	<-stopCh
	log.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (xc *XGBoostController) runWorker() {
	for xc.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (xc *XGBoostController) processNextWorkItem() bool {
	obj, quit := xc.WorkQueue.Get()
	if quit {
		return false
	}
	defer xc.WorkQueue.Done(obj)

	var key string
	var ok bool
	if key, ok = obj.(string); !ok {
		// As the item in the workqueue is actually invalid, we call
		// Forget here else we'd go into a loop of attempting to
		// process a work item that is invalid.
		xc.WorkQueue.Forget(obj)
		utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
		return true
	}

	logger := commonutil.LoggerForKey(key)

	xgboostJob, err := xc.getXGBoostJobFromKey(key)

	if err != nil {
		if err == errNotExists {
			logger.Infof("XGBoostJob has been deleted: %v", key)
			return true
		}

		// Log the failure to conditions.
		logger.Errorf("Failed to get XGBoostJob from key %s: %v", key, err)
		if err == errFailedMarshal {
			errMsg := fmt.Sprintf("Failed to unmarshal the object to XGBoostJob object: %v", err)
			commonutil.LoggerForJob(xgboostJob).Warn(errMsg)
			xc.Recorder.Event(xgboostJob, corev1.EventTypeWarning, failedMarshalXGBoostJobReason, errMsg)
		}

		return true
	}

	// Sync XGBoostJob to match the actual state to this desired state.
	forget, err := xc.syncHandler(key)
	if err == nil {
		if forget {
			xc.WorkQueue.Forget(key)
		}
		return true
	}

	xc.WorkQueue.AddRateLimited(key)

	return true
}

// syncXGBoostJob syncs the job with the given key if it has had its expectations fulfilled, meaning
// it did not expect to see any more of its pods/services created or deleted.
// This function is not meant to be invoked concurrently with the same key.
func (xc *XGBoostController) syncXGBoostJob(key string) (bool, error) {
	startTime := time.Now()
	logger := commonutil.LoggerForKey(key)
	defer func() {
		logger.Infof("Finished syncing job %q (%v)", key, time.Since(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return false, err
	}
	if len(namespace) == 0 || len(name) == 0 {
		return false, fmt.Errorf("invalid job key %q: either namespace or name is missing", key)
	}

	sharedJob, err := xc.getXGBoostJobFromName(namespace, name)
	if err != nil {
		if err == errNotExists {
			logger.Infof("XGBoostJob has been deleted: %v", key)
			// jm.expectations.DeleteExpectations(key)
			return true, nil
		}
		return false, err
	}

	job := sharedJob.DeepCopy()
	jobNeedsSync := xc.satisfiedExpectations(job)

	if xc.Config.EnableGangScheduling {
		minAvailableReplicas := k8sutil.GetTotalReplicas(job.Spec.XGBoostReplicaSpecs)
		_, err := xc.SyncPodGroup(job, minAvailableReplicas)
		if err != nil {
			logger.Warnf("Sync PodGroup %v: %v", job.Name, err)
		}
	}

	// Set default for the new job.
	scheme.Scheme.Default(job)

	var reconcileXGBoostJobsErr error
	if jobNeedsSync && job.DeletionTimestamp == nil {
		reconcileXGBoostJobsErr = xc.reconcileXGBoostJobs(
			job, job.Spec.XGBoostReplicaSpecs, job.Status, job.Spec.RunPolicy)
	}

	if reconcileXGBoostJobsErr != nil {
		return false, reconcileXGBoostJobsErr
	}

	return true, err
}

// satisfiedExpectations returns true if the required adds/dels for the given job have been observed.
// Add/del counts are established by the controller at sync time, and updated as controllees are observed by the controller
// manager.
func (xc *XGBoostController) satisfiedExpectations(job *v1alpha1.XGBoostJob) bool {
	jobKey, err := jobcontroller.KeyFunc(job)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for job object %#v: %v", job, err))
		return false
	}

	for rtype := range job.Spec.XGBoostReplicaSpecs {
		// Check the expectations of the pods.
		expectationPodsKey := jobcontroller.GenExpectationPodsKey(jobKey, string(rtype))
		if xc.Expectations.SatisfiedExpectations(expectationPodsKey) {
			return true
		}

		// Check the expectations of the services.
		expectationServicesKey := jobcontroller.GenExpectationServicesKey(jobKey, string(rtype))
		if xc.Expectations.SatisfiedExpectations(expectationServicesKey) {
			return true
		}
	}

	return false
}

// reconcilePyTorchJobs checks and updates replicas for each given PyTorchReplicaSpec.
// It will requeue the job in case of an error while creating/deleting pods/services.
func (xc *XGBoostController) reconcileXGBoostJobs(
	job *v1alpha1.XGBoostJob,
	replicas map[common.ReplicaType]*common.ReplicaSpec,
	jobStatus common.JobStatus,
	runPolicy *common.RunPolicy) error {
	///TODO
	return nil
}

func (XGBoostController) ControllerName() string {
	return controllerName
}

func (XGBoostController) GetAPIGroupVersionKind() schema.GroupVersionKind {
	return v1alpha1.SchemeGroupVersionKind
}

func (XGBoostController) GetAPIGroupVersion() schema.GroupVersion {
	return v1alpha1.SchemeGroupVersion
}

func (XGBoostController) GetGroupNameLabelValue() string {
	return v1alpha1.GroupName
}

func (XGBoostController) GetDefaultContainerName() string {
	return v1alpha1.DefaultContainerName
}

func (XGBoostController) GetDefaultContainerPortNumber() string {
	return string(v1alpha1.DefaultPort)
}

func (XGBoostController) GetJobRoleKey() string {
	return ""
}

func (XGBoostController) IsMasterRole(replicas map[common.ReplicaType]*common.ReplicaSpec,
	rtype common.ReplicaType, index int) bool {
	return string(rtype) == string(v1alpha1.XGBoostReplicaTypeMaster)
}


func (xc *XGBoostController) GetServicesForJob(job *v1alpha1.XGBoostJob) ([]*corev1.Service, error) {
	///TODO
	return nil, nil
}

func (xc *XGBoostController) SetClusterSpec(job interface {}, podTemplate *corev1.PodTemplateSpec, rtype, index string) error {
	///TODO
	return nil
}