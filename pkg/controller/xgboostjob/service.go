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
	"fmt"
	"github.com/kubeflow/common/job_controller"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboostjob/v1alpha1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateService creates the service
func (r *ReconcileXGBoostJob) CreateService(job interface{}, service *corev1.Service) error {
	xgboostjob, ok := job.(*v1alpha1.XGBoostJob)
	if !ok {
		return fmt.Errorf("%+v is not a type of XGBoostJob", xgboostjob)
	}

	logrus.Info("Creating service ", " Controller name ", xgboostjob.GetName(), " Service name ", service.Namespace+"/"+service.Name)

	//service, err := r.xgbJobController.KubeClientSet.CoreV1().Services(xgboostjob.Namespace).Create(service)
	err := r.Create(context.Background(), service)

	if err != nil {
		logrus.Warnf("Create service error %s", xgboostjob.Name)
	}

	return err
}

// DeleteService deletes the service
func (r *ReconcileXGBoostJob) DeleteService(job interface{}, name string, namespace string) error {
	xgboostjob, ok := job.(*v1alpha1.XGBoostJob)
	if !ok {
		return fmt.Errorf("%+v is not a type of XGBoostJob", xgboostjob)
	}

	service := &corev1.Service{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name}}

	logrus.Info("Deleting service ", " Controller name ", xgboostjob.GetName(), " Service name ", service.Namespace+"/"+service.Name)

	if err := r.Delete(context.Background(), service); err != nil {
		r.recorder.Eventf(xgboostjob, corev1.EventTypeWarning, job_controller.FailedDeleteServiceReason, "Error deleting: %v", err)
		return fmt.Errorf("unable to delete service: %v", err)
	}

	r.recorder.Eventf(xgboostjob, corev1.EventTypeNormal, job_controller.SuccessfulDeleteServiceReason, "Deleted service: %v", name)

	return nil

}

// GetServicesForJob returns the services managed by the job. This can be achieved by selecting services using label key "job-name"
// i.e. all services created by the job will come with label "job-name" = <this_job_name>
func (r *ReconcileXGBoostJob) GetServicesForJob(obj interface{}) ([]*corev1.Service, error) {
	job, err := meta.Accessor(obj)
	if err != nil {
		return nil, fmt.Errorf("%+v is not a type of XGBoostJob", job)
	}
	// List all pods to include those that don't match the selector anymore
	// but have a ControllerRef pointing to this controller.
	serviceList := &corev1.ServiceList{}
	err = r.List(context.Background(), client.MatchingLabels(r.xgbJobController.GenLabels(job.GetName())), serviceList)
	if err != nil {
		return nil, err
	}
	//TODO support adopting/orphaning
	return convertServiceList(serviceList.Items), nil
}

// convertServiceList convert service list to service point list
func convertServiceList(list []corev1.Service) []*corev1.Service {
	if list == nil {
		return nil
	}
	ret := make([]*corev1.Service, 0, len(list))
	for i := range list {
		ret = append(ret, &list[i])
	}
	return ret
}
