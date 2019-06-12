/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package xgboostjob

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/kubeflow/common/job_controller"
	"github.com/kubeflow/common/job_controller/api/v1"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboostjob/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	k8scontroller "k8s.io/kubernetes/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	controllerName      = "xgboostjob-operator"
	labelXGBoostJobRole = "xgboostjob-job-role"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new XGBoostJob Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

const RecommendedKubeConfigPathEnv = "KUBECONFIG"

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {

	r := &ReconcileXGBoostJob{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}

	r.recorder = mgr.GetRecorder(r.ControllerName())

	var kubeconfig *string

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig_tmp", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig_tmp", "", "absolute path to the kubeconfig file")
	}

	/// TODO, add the master url and kubeconfigpath with user input
	kcfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Info("Error building kubeconfig: %s", err.Error())
	}

	// Create clients.
	kubeClientSet, _, kubeBatchClientSet, err := createClientSets(kcfg)
	if err != nil {
		log.Info("Error building kubeclientset: %s", err.Error())
	}

	// Initialize common job controller with components we only need.
	r.xgbJobController = job_controller.JobController{
		Controller:         r,
		Expectations:       k8scontroller.NewControllerExpectations(),
		Config:             v1.JobControllerConfiguration{EnableGangScheduling: true},
		WorkQueue:          &FakeWorkQueue{},
		Recorder:           r.recorder,
		KubeClientSet:      kubeClientSet,
		KubeBatchClientSet: kubeBatchClientSet,
	}

	return r
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("xgboostjob-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to XGBoostJob
	err = c.Watch(&source.Kind{Type: &v1alpha1.XGBoostJob{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.XGBoostJob{},
	})
	if err != nil {
		return err
	}

	//inject watching for  xgboostjob related pod
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.XGBoostJob{},
	},
		predicate.Funcs{CreateFunc: onDependentCreateFunc(r), DeleteFunc: onDependentDeleteFunc(r)},
	)
	if err != nil {
		return err
	}

	//inject watching for xgboostjob related service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.XGBoostJob{},
	},
		&predicate.Funcs{CreateFunc: onDependentCreateFunc(r), DeleteFunc: onDependentDeleteFunc(r)},
	)
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileXGBoostJob{}

// ReconcileXGBoostJob reconciles a XGBoostJob object
type ReconcileXGBoostJob struct {
	client.Client
	scheme           *runtime.Scheme
	xgbJobController job_controller.JobController
	recorder         record.EventRecorder
}

// Reconcile reads that state of the cluster for a XGBoostJob object and makes changes based on the state read
// and what is in the XGBoostJob.Spec
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=xgboostjob.kubeflow.org,resources=xgboostjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=xgboostjob.kubeflow.org,resources=xgboostjobs/status,verbs=get;update;patch
func (r *ReconcileXGBoostJob) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the XGBoostJob instance
	xgboostjob := &v1alpha1.XGBoostJob{}
	err := r.Get(context.Background(), request.NamespacedName, xgboostjob)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check reconcile is required.
	needSync := r.satisfiedExpectations(xgboostjob)

	if !needSync || xgboostjob.DeletionTimestamp != nil {
		log.Info("reconcile cancelled, job does not need to do reconcile or has been deleted",
			"sync", needSync, "deleted", xgboostjob.DeletionTimestamp != nil)
		return reconcile.Result{}, nil
	}
	oldStatus := xgboostjob.Status.DeepCopy()
	// Set default priorities for xgboost job
	scheme.Scheme.Default(xgboostjob)

	// Use common to reconcile the job related pod and service
	err = r.xgbJobController.ReconcileJobs(xgboostjob, xgboostjob.Spec.XGBReplicaSpecs, xgboostjob.Status.JobStatus, &xgboostjob.Spec.RunPolicy)

	if err != nil {
		return reconcile.Result{}, err
	}
	if !reflect.DeepEqual(oldStatus, &xgboostjob.Status.JobStatus) {
		err = r.UpdateJobStatusInApiServer(xgboostjob, &xgboostjob.Status.JobStatus)
	}
	return reconcile.Result{}, err
}

// UpdateJobStatusInApiServer updates the job status in to cluster.
func (r *ReconcileXGBoostJob) UpdateJobStatusInApiServer(job interface{}, jobStatus *v1.JobStatus) error {
	xgboostjob, ok := job.(*v1alpha1.XGBoostJob)
	if !ok {
		return fmt.Errorf("%+v is not a type of XGBoostJob", xgboostjob)
	}
	// Job status passed in differs with status in job, update in basis of the passed in one.
	if !reflect.DeepEqual(&xgboostjob.Status.JobStatus, jobStatus) {
		xgboostjob = xgboostjob.DeepCopy()
		xgboostjob.Status.JobStatus = *jobStatus.DeepCopy()
	}
	return r.Status().Update(context.Background(), xgboostjob)
}

func (r *ReconcileXGBoostJob) ControllerName() string {
	return controllerName
}

func (r *ReconcileXGBoostJob) GetAPIGroupVersionKind() schema.GroupVersionKind {
	return v1alpha1.SchemeGroupVersionKind
}

func (r *ReconcileXGBoostJob) GetAPIGroupVersion() schema.GroupVersion {
	return v1alpha1.SchemeGroupVersion
}

func (r *ReconcileXGBoostJob) GetGroupNameLabelValue() string {
	return v1alpha1.GroupName
}

func (r *ReconcileXGBoostJob) GetDefaultContainerName() string {
	return v1alpha1.DefaultContainerName
}

func (r *ReconcileXGBoostJob) GetDefaultContainerPortNumber() int32 {
	return v1alpha1.DefaultPort
}

func (r *ReconcileXGBoostJob) GetDefaultContainerPortName() string {
	return v1alpha1.DefaultContainerPortName
}

func (r *ReconcileXGBoostJob) GetJobRoleKey() string {
	return labelXGBoostJobRole
}

func (r *ReconcileXGBoostJob) IsMasterRole(replicas map[v1.ReplicaType]*v1.ReplicaSpec,
	rtype v1.ReplicaType, index int) bool {
	return string(rtype) == string(v1alpha1.XGBoostReplicaTypeMaster)
}

// SetClusterSpec sets the cluster spec for the pod
func (r *ReconcileXGBoostJob) SetClusterSpec(job interface{}, podTemplate *corev1.PodTemplateSpec, rtype, index string) error {
	return SetPodEnv(job, podTemplate, index)
}
