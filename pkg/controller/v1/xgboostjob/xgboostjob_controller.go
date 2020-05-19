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
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
	"github.com/kubeflow/common/pkg/controller.v1/common"
	"github.com/kubeflow/common/pkg/controller.v1/control"
	"github.com/kubeflow/common/pkg/controller.v1/expectation"
	v1xgboost "github.com/kubeflow/xgboost-operator/pkg/apis/xgboostjob/v1"
	corev1 "k8s.io/api/core/v1"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
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
	// gang scheduler name.
	gangSchedulerName = "kube-batch"
)

var (
	defaultTTLseconds     = int32(100)
	defaultCleanPodPolicy = commonv1.CleanPodPolicyNone
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

	r.recorder = mgr.GetEventRecorderFor(r.ControllerName())

	var mode string
	var kubeconfig *string
	var kcfg *rest.Config
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig_", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig_", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	mode = flag.Lookup("mode").Value.(flag.Getter).Get().(string)
	if mode == "local" {
		log.Info("Running controller in local mode, using kubeconfig file")
		/// TODO, add the master url and kubeconfigpath with user input
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Info("Error building kubeconfig: %s", err.Error())
			panic(err.Error())
		}
		kcfg = config
	} else if mode == "in-cluster" {
		log.Info("Running controller in in-cluster mode")
		/// TODO, add the master url and kubeconfigpath with user input
		kcfg, err := rest.InClusterConfig()
		if err != nil {
			log.Info("Error getting in-cluster kubeconfig")
			panic(err.Error())
		}
		_ = kcfg
	} else {
		log.Info("Given mode is not valid: ", "mode", mode)
		panic("-mode should be either local or in-cluster")
	}

	// Create clients.
	kubeClientSet, _, volcanoClientSet, err := createClientSets(kcfg)
	if err != nil {
		log.Info("Error building kubeclientset: %s", err.Error())
	}

	// Create Informer factory
	xgboostjob := &v1xgboost.XGBoostJob{}

	gangScheduling := isGangSchedulerSet(xgboostjob.Spec.XGBReplicaSpecs)

	log.Info("gang scheduling is set: ", "gangscheduling", gangScheduling)

	// Initialize common job controller with components we only need.
	r.JobController = common.JobController{
		Controller:       r,
		Expectations:     expectation.NewControllerExpectations(),
		Config:           common.JobControllerConfiguration{EnableGangScheduling: gangScheduling},
		WorkQueue:        &FakeWorkQueue{},
		Recorder:         r.recorder,
		KubeClientSet:    kubeClientSet,
		VolcanoClientSet: volcanoClientSet,
		PodControl:       control.RealPodControl{KubeClient: kubeClientSet, Recorder: r.recorder},
		ServiceControl:   control.RealServiceControl{KubeClient: kubeClientSet, Recorder: r.recorder},
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
	err = c.Watch(&source.Kind{Type: &v1xgboost.XGBoostJob{}}, &handler.EnqueueRequestForObject{},
		predicate.Funcs{CreateFunc: onOwnerCreateFunc(r)},
	)
	if err != nil {
		return err
	}

	//inject watching for  xgboostjob related pod
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1xgboost.XGBoostJob{},
	},
		predicate.Funcs{CreateFunc: onDependentCreateFunc(r), DeleteFunc: onDependentDeleteFunc(r)},
	)
	if err != nil {
		return err
	}

	//inject watching for xgboostjob related service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1xgboost.XGBoostJob{},
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
	common.JobController
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
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
	xgboostjob := &v1xgboost.XGBoostJob{}
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
	// Set default priorities for xgboost job
	scheme.Scheme.Default(xgboostjob)

	// Use common to reconcile the job related pod and service
	err = r.ReconcileJobs(xgboostjob, xgboostjob.Spec.XGBReplicaSpecs, xgboostjob.Status.JobStatus, &xgboostjob.Spec.RunPolicy)

	if err != nil {
		logrus.Warnf("Reconcile XGBoost Job error %v", err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, err
}

func (r *ReconcileXGBoostJob) ControllerName() string {
	return controllerName
}

func (r *ReconcileXGBoostJob) GetAPIGroupVersionKind() schema.GroupVersionKind {
	return v1xgboost.SchemeGroupVersionKind
}

func (r *ReconcileXGBoostJob) GetAPIGroupVersion() schema.GroupVersion {
	return v1xgboost.SchemeGroupVersion
}

func (r *ReconcileXGBoostJob) GetGroupNameLabelValue() string {
	return v1xgboost.GroupName
}

func (r *ReconcileXGBoostJob) GetDefaultContainerName() string {
	return v1xgboost.DefaultContainerName
}

func (r *ReconcileXGBoostJob) GetDefaultContainerPortName() string {
	return v1xgboost.DefaultContainerPortName
}

func (r *ReconcileXGBoostJob) GetJobRoleKey() string {
	return labelXGBoostJobRole
}

func (r *ReconcileXGBoostJob) IsMasterRole(replicas map[commonv1.ReplicaType]*commonv1.ReplicaSpec,
	rtype commonv1.ReplicaType, index int) bool {
	return string(rtype) == string(v1xgboost.XGBoostReplicaTypeMaster)
}

// SetClusterSpec sets the cluster spec for the pod
func (r *ReconcileXGBoostJob) SetClusterSpec(job interface{}, podTemplate *corev1.PodTemplateSpec, rtype, index string) error {
	return SetPodEnv(job, podTemplate, rtype, index)
}