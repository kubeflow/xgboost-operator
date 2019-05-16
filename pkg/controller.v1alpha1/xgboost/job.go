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
	"github.com/kubeflow/common/util/k8sutil"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"time"

	jobcontroller "github.com/kubeflow/common/job_controller"
	commonv1 "github.com/kubeflow/common/operator/v1"
	commonutil "github.com/kubeflow/common/util"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
)

const failedMarshalXGBoostJobReason = "InvalidXGBoostJobSpec"

// When a pod is added, set the defaults and enqueue the current xgboostjob.
func (xc *XGBoostController) addXGBoostJob(obj interface{}) {
	// Convert from unstructured object.
	xgboostJob, err := jobFromUnstructured(obj)
	if err != nil {
		un, ok := obj.(*metav1unstructured.Unstructured)
		logger := &log.Entry{}
		if ok {
			logger = commonutil.LoggerForUnstructured(un, v1alpha1.Kind)
		}
		logger.Errorf("Failed to convert the XGBoostJob: %v", err)
		// Log the failure to conditions.
		if err == errFailedMarshal {
			errMsg := fmt.Sprintf("Failed to marshal the object to XGBoostJob; the spec is invalid: %v", err)
			logger.Warn(errMsg)
			// TODO(jlewi): v1 doesn't appear to define an error type.
			xc.Recorder.Event(un, v1.EventTypeWarning, failedMarshalXGBoostJobReason, errMsg)

			status := commonv1.JobStatus{
				Conditions: []commonv1.JobCondition{
					commonv1.JobCondition{
						Type:               commonv1.JobFailed,
						Status:             v1.ConditionTrue,
						LastUpdateTime:     metav1.Now(),
						LastTransitionTime: metav1.Now(),
						Reason:             failedMarshalXGBoostJobReason,
						Message:            errMsg,
					},
				},
			}

			statusMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&status)

			if err != nil {
				logger.Errorf("Could not covert the XGBoostJobStatus to unstructured; %v", err)
				return
			}

			client, err := k8sutil.NewCRDRestClient(&v1alpha1.SchemeGroupVersion)

			if err == nil {
				if err1 := metav1unstructured.SetNestedField(un.Object, statusMap, "status"); err1 != nil {
					logger.Errorf("Could not set nested field: %v", err1)
				}
				logger.Infof("Updating the job to: %+v", un.Object)
				err = client.UpdateStatus(un, v1alpha1.Plural)
				if err != nil {
					logger.Errorf("Could not update the XGBoostJob: %v", err)
				}
			} else {
				logger.Errorf("Could not create a REST client to update the XGBoostJob")
			}
		}
		return
	}

	// Set default for the new xgboostjob.
	scheme.Scheme.Default(xgboostJob)

	msg := fmt.Sprintf("XGBoostJob %s is created.", xgboostJob.Name)
	logger := commonutil.LoggerForJob(xgboostJob)
	logger.Info(msg)

	// Add a created condition.
	err = commonutil.UpdateJobConditions(
		&xgboostJob.Status, commonv1.JobCreated, commonutil.JobCreatedReason, msg)
	if err != nil {
		logger.Errorf("Append xgboostJob condition error: %v", err)
		return
	}

	// Convert from xgboostjob object
	err = unstructuredFromJob(obj, xgboostJob)
	if err != nil {
		logger.Errorf("Failed to convert the obj: %v", err)
		return
	}
	xc.enqueueXGBoostJob(obj)
}

// When a pod is updated, enqueue the current xgboostjob.
func (xc *XGBoostController) updateXGBoostJob(old, cur interface{}) {
	oldXGBoostJob, err := jobFromUnstructured(old)
	if err != nil {
		return
	}
	curXGBoostJob, err := jobFromUnstructured(cur)
	if err != nil {
		return
	}

	// never return error
	key, err := jobcontroller.KeyFunc(curXGBoostJob)
	if err != nil {
		return
	}

	log.Infof("Updating xgboostjob: %s", oldXGBoostJob.Name)
	xc.enqueueXGBoostJob(cur)

	// check if need to add a new rsync for ActiveDeadlineSeconds
	if curXGBoostJob.Status.StartTime != nil {
		curXGBoostJobADS := curXGBoostJob.Spec.RunPolicy.ActiveDeadlineSeconds
		if curXGBoostJobADS == nil {
			return
		}
		oldXGBoostJobADS := oldXGBoostJob.Spec.RunPolicy.ActiveDeadlineSeconds
		if oldXGBoostJobADS == nil || *oldXGBoostJobADS != *curXGBoostJobADS {
			now := metav1.Now()
			start := curXGBoostJob.Status.StartTime.Time
			passed := now.Time.Sub(start)
			total := time.Duration(*curXGBoostJobADS) * time.Second
			// AddAfter will handle total < passed
			xc.WorkQueue.AddAfter(key, total-passed)
			log.Infof("job ActiveDeadlineSeconds updated, will rsync after %d seconds", total-passed)
		}
	}
}

func (xc *XGBoostController) GetJobFromInformerCache(namespace, name string) (metav1.Object, error) {
	return xc.getXGBoostJobFromName(namespace, name)
}

func (xc *XGBoostController) GetJobFromAPIClient(namespace, name string) (metav1.Object, error) {
	return xc.jobClientSet.KubeflowV1alpha1().XGBoostJobs(namespace).Get(name, metav1.GetOptions{})
}

func (xc *XGBoostController) DeleteJob(job interface{}) error {
	log.Info("Deleting job")
	xgbJob := job.(*v1alpha1.XGBoostJob)
	return xc.jobClientSet.KubeflowV1alpha1().XGBoostJobs(xgbJob.Namespace).Delete(xgbJob.Name, &metav1.DeleteOptions{})
}
