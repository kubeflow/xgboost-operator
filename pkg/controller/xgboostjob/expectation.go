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
	"github.com/kubeflow/common/job_controller"
	"github.com/kubeflow/common/job_controller/api/v1"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboostjob/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// satisfiedExpectations returns true if the required adds/dels for the given job have been observed.
// Add/del counts are established by the controller at sync time, and updated as controllees are observed by the controller
// manager.
func (r *ReconcileXGBoostJob) satisfiedExpectations(xgbJob *v1alpha1.XGBoostJob) bool {
	satisfied := false
	key, err := job_controller.KeyFunc(xgbJob)
	if err != nil {
		return false
	}
	for rtype := range xgbJob.Spec.XGBReplicaSpecs {
		// Check the expectations of the pods.
		expectationPodsKey := job_controller.GenExpectationPodsKey(key, string(rtype))
		satisfied = satisfied || r.xgbJobController.Expectations.SatisfiedExpectations(expectationPodsKey)

		// Check the expectations of the services.
		expectationServicesKey := job_controller.GenExpectationServicesKey(key, string(rtype))
		satisfied = satisfied || r.xgbJobController.Expectations.SatisfiedExpectations(expectationServicesKey)
	}
	return satisfied
}

// onDependentCreateFunc modify expectations when dependent (pod/service) creation observed.
func onDependentCreateFunc(r reconcile.Reconciler) func(event.CreateEvent) bool {
	return func(e event.CreateEvent) bool {
		xgbr, ok := r.(*ReconcileXGBoostJob)
		if !ok {
			return true
		}
		rtype := e.Meta.GetLabels()[v1.ReplicaTypeLabel]
		key, err := job_controller.KeyFunc(e.Meta)
		if err != nil {
			return false
		}
		expectKey := job_controller.GenExpectationPodsKey(key, rtype)
		xgbr.xgbJobController.Expectations.CreationObserved(expectKey)
		return true
	}
}

// onDependentDeleteFunc modify expectations when dependent (pod/service) deletion observed.
func onDependentDeleteFunc(r reconcile.Reconciler) func(event.DeleteEvent) bool {
	return func(e event.DeleteEvent) bool {
		xgbr, ok := r.(*ReconcileXGBoostJob)
		if !ok {
			return true
		}
		rtype := e.Meta.GetLabels()[v1.ReplicaTypeLabel]
		key, err := job_controller.KeyFunc(e.Meta)
		if err != nil {
			return false
		}
		expectKey := job_controller.GenExpectationPodsKey(key, rtype)
		xgbr.xgbJobController.Expectations.DeletionObserved(expectKey)
		return true
	}
}
