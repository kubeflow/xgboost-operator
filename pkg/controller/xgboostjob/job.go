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
	v1 "github.com/kubeflow/common/job_controller/api/v1"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboostjob/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Reasons for job events.
const (
	FailedDeleteJobReason     = "FailedDeleteJob"
	SuccessfulDeleteJobReason = "SuccessfulDeleteJob"
)

// DeleteJob deletes the job
func (r *ReconcileXGBoostJob) DeleteJob(job interface{}) error {
	xgbjob, ok := job.(*v1alpha1.XGBoostJob)
	if !ok {
		return fmt.Errorf("%+v is not a type of XGBoostJob", xgbjob)
	}
	if err := r.Delete(context.Background(), xgbjob); err != nil {
		r.recorder.Eventf(xgbjob, corev1.EventTypeWarning, FailedDeleteJobReason, "Error deleting: %v", err)
		log.Error(err, "failed to delete job", "namespace", xgbjob.Namespace, "name", xgbjob.Name)
		return err
	}
	r.recorder.Eventf(xgbjob, corev1.EventTypeNormal, SuccessfulDeleteJobReason, "Deleted job: %v", xgbjob.Name)
	log.Info("job deleted", "namespace", xgbjob.Namespace, "name", xgbjob.Name)
	return nil
}

// GetJobFromInformerCache returns the Job from Informer Cache
func (r *ReconcileXGBoostJob) GetJobFromInformerCache(namespace, name string) (metav1.Object, error) {
	job := &v1alpha1.XGBoostJob{}
	// Default reader for XGBoostJob is cache reader.
	err := r.Get(context.Background(), types.NamespacedName{Namespace: namespace, Name: name}, job)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Error(err, "xgboost job not found", "namespace", namespace, "name", name)
		} else {
			log.Error(err, "failed to get job from api-server", "namespace", namespace, "name", name)
		}
		return nil, err
	}
	return job, nil
}

// GetJobFromAPIClient returns the Job from API server
func (r *ReconcileXGBoostJob) GetJobFromAPIClient(namespace, name string) (metav1.Object, error) {
	job := &v1alpha1.XGBoostJob{}

	clientReader, err := getClientReaderFromClient(r.Client)
	if err != nil {
		return nil, err
	}
	err = clientReader.Get(context.Background(), types.NamespacedName{Namespace: namespace, Name: name}, job)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Error(err, "xgboost job not found", "namespace", namespace, "name", name)
		} else {
			log.Error(err, "failed to get job from api-server", "namespace", namespace, "name", name)
		}
		return nil, err
	}
	return job, nil
}

// UpdateJobStatus updates the job status and job conditions
func (r *ReconcileXGBoostJob) UpdateJobStatus(job interface{}, replicas map[v1.ReplicaType]*v1.ReplicaSpec, jobStatus v1.JobStatus) error {
	xgboostJob, ok := job.(*v1alpha1.XGBoostJob)
	if !ok {
		return fmt.Errorf("%+v is not a type of xgboostJob", xgboostJob)
	}
	//TODO implement this for updating job statue
	return nil
}
