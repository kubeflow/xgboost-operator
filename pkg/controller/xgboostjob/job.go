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
	"reflect"

	v1 "github.com/kubeflow/common/job_controller/api/v1"
	commonutil "github.com/kubeflow/common/util"
	"k8s.io/client-go/kubernetes/scheme"
	logger "github.com/kubeflow/common/util"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboostjob/v1alpha1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reasons for job events.
const (
	FailedDeleteJobReason     = "FailedDeleteJob"
	SuccessfulDeleteJobReason = "SuccessfulDeleteJob"
	// xgboostJobCreatedReason is added in a job when it is created.
	xgboostJobCreatedReason = "XGBoostJobCreated"

	xgboostJobSucceededReason  = "XGBoostJobSucceeded"
	xgboostJobRunningReason    = "XGBoostJobRunning"
	xgboostJobFailedReason     = "XGBoostJobFailed"
	xgboostJobRestartingReason = "XGBoostJobRestarting"
)

// DeleteJob deletes the job
func (r *ReconcileXGBoostJob) DeleteJob(job interface{}) error {
	xgboostjob, ok := job.(*v1alpha1.XGBoostJob)
	if !ok {
		return fmt.Errorf("%+v is not a type of XGBoostJob", xgboostjob)
	}
	if err := r.Delete(context.Background(), xgboostjob); err != nil {
		r.recorder.Eventf(xgboostjob, corev1.EventTypeWarning, FailedDeleteJobReason, "Error deleting: %v", err)
		log.Error(err, "failed to delete job", "namespace", xgboostjob.Namespace, "name", xgboostjob.Name)
		return err
	}
	r.recorder.Eventf(xgboostjob, corev1.EventTypeNormal, SuccessfulDeleteJobReason, "Deleted job: %v", xgboostjob.Name)
	log.Info("job deleted", "namespace", xgboostjob.Namespace, "name", xgboostjob.Name)
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
func (r *ReconcileXGBoostJob) UpdateJobStatus(job interface{}, replicas map[v1.ReplicaType]*v1.ReplicaSpec, jobStatus *v1.JobStatus) error {
	xgboostJob, ok := job.(*v1alpha1.XGBoostJob)
	if !ok {
		return fmt.Errorf("%+v is not a type of xgboostJob", xgboostJob)
	}

	logrus.Info("job status and time", jobStatus.StartTime)

	for rtype, spec := range replicas {
		status := jobStatus.ReplicaStatuses[rtype]

		expected := *(spec.Replicas) - status.Succeeded
		running := status.Active
		failed := status.Failed

		logrus.Infof("XGboostJob=%s, ReplicaType=%s expected=%d, running=%d, failed=%d",
			xgboostJob.Name, rtype, expected, running, failed)

		if rtype == v1.ReplicaType(v1alpha1.XGBoostReplicaTypeMaster) {
			if running > 0 {
				msg := fmt.Sprintf("XGBoostJob %s is running.", xgboostJob.Name)
				err := commonutil.UpdateJobConditions(jobStatus, v1.JobRunning, xgboostJobRunningReason, msg)
				if err != nil {
					logger.LoggerForJob(xgboostJob).Infof("Append job condition error: %v", err)
					return err
				}
			}
			if expected == 0 {
				msg := fmt.Sprintf("XGBoostJob %s is successfully completed.", xgboostJob.Name)
				r.xgbJobController.Recorder.Event(xgboostJob, k8sv1.EventTypeNormal, xgboostJobSucceededReason, msg)
				if jobStatus.CompletionTime == nil {
					now := metav1.Now()
					jobStatus.CompletionTime = &now
				}
				err := commonutil.UpdateJobConditions(jobStatus, v1.JobSucceeded, xgboostJobSucceededReason, msg)
				if err != nil {
					logger.LoggerForJob(xgboostJob).Infof("Append job condition error: %v", err)
					return err
				}
			}
		}
		if failed > 0 {
			if spec.RestartPolicy == v1.RestartPolicyExitCode {
				msg := fmt.Sprintf("XGBoostJob %s is restarting because %d %s replica(s) failed.", xgboostJob.Name, failed, rtype)
				r.xgbJobController.Recorder.Event(xgboostJob, k8sv1.EventTypeWarning, xgboostJobRestartingReason, msg)
				err := commonutil.UpdateJobConditions(jobStatus, v1.JobRestarting, xgboostJobRestartingReason, msg)
				if err != nil {
					logger.LoggerForJob(xgboostJob).Infof("Append job condition error: %v", err)
					return err
				}
			} else {
				msg := fmt.Sprintf("XGBoostJob %s is failed because %d %s replica(s) failed.", xgboostJob.Name, failed, rtype)
				r.xgbJobController.Recorder.Event(xgboostJob, k8sv1.EventTypeNormal, xgboostJobFailedReason, msg)
				if xgboostJob.Status.CompletionTime == nil {
					now := metav1.Now()
					xgboostJob.Status.CompletionTime = &now
				}
				err := commonutil.UpdateJobConditions(jobStatus, v1.JobFailed, xgboostJobFailedReason, msg)
				if err != nil {
					logger.LoggerForJob(xgboostJob).Infof("Append job condition error: %v", err)
					return err
				}
			}
		}
	}

	// Some workers are still running, leave a running condition.
	msg := fmt.Sprintf("XGBoostJob %s is running.", xgboostJob.Name)
	logger.LoggerForJob(xgboostJob).Infof(msg)

	if err := commonutil.UpdateJobConditions(jobStatus, v1.JobRunning, xgboostJobSucceededReason, msg); err != nil {
		logger.LoggerForJob(xgboostJob).Error(err, "failed to update XGBoost Job conditions")
		return err
	}

	return nil
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

	result := r.Update(context.Background(), xgboostjob)

	if result != nil{
		logger.LoggerForJob(xgboostjob).Error(result, "failed to update XGBoost Job conditions in the API server")
		return result
	}

	return nil
}

// onOwnerCreateFunc modify creation condition.
func onOwnerCreateFunc(r reconcile.Reconciler) func(event.CreateEvent) bool {
	return func(e event.CreateEvent) bool {
		xgboostJob, ok := e.Meta.(*v1alpha1.XGBoostJob)
		if !ok {
			return true
		}
		scheme.Scheme.Default(xgboostJob)
		msg := fmt.Sprintf("xgboostJob %s is created.", e.Meta.GetName())
		if err := commonutil.UpdateJobConditions(&xgboostJob.Status.JobStatus, v1.JobCreated, xgboostJobCreatedReason, msg); err != nil {
			log.Error(err, "append job condition error")
			return false
		}
		return true
	}
}

