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
	commonv1 "github.com/kubeflow/common/operator/v1"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
)

// uncomment Reasons temporarily to pass golang ci-lint, since they are unused for now.
/*
const (
	// xgboostJobCreatedReason is added in a job when it is created.
	xgboostJobCreatedReason = "XGBoostJobCreated"
	// xgboostJobSucceededReason is added in a job when it is succeeded.
	xgboostJobSucceededReason = "XGBoostJobSucceeded"
	// xgboostJobRunningReason is added in a job when it is running.
	xgboostJobRunningReason = "XGBoostJobRunning"
	// xgboostJobFailedReason is added in a job when it is failed.
	xgboostJobFailedReason = "XGBoostJobFailed"
	// xgboostJobRestarting is added in a job when it is restarting.
	xgboostJobRestartingReason = "XGBoostJobRestarting"
)
*/

func (xc *XGBoostController) UpdateJobStatus(job interface{},
	replicas map[commonv1.ReplicaType]*commonv1.ReplicaSpec,
	jobStatus *commonv1.JobStatus) error {
	//	TODO: Implement this method
	return nil
}

func (xc *XGBoostController) UpdateJobStatusInApiServer(job interface{}, jobStatus *commonv1.JobStatus) error {
	xgbJob := job.(*v1alpha1.XGBoostJob)
	jobStatus.DeepCopyInto(&xgbJob.Status)
	_, err := xc.jobClientSet.KubeflowV1alpha1().XGBoostJobs(xgbJob.Namespace).UpdateStatus(xgbJob)
	return err
}
