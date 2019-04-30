// Copyright 2018 The Kubeflow Authors
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

package v1beta1

import (
	commmonv1 "github.com/kubeflow/common/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=xgboostjob

// XGBoostJob represents the configuration of XGBoostJob
type XGBoostJob struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the XGBoostJob.
	Spec XGBoostJobSpec `json:"spec,omitempty"`

	// Most recently observed status of the XGBoostJob.
	// This data may not be up to date.
	// Populated by the system.
	// Read-only.
	Status commmonv1.JobStatus `json:"status,omitempty"`

	commmonv1.CleanPodPolicy
}

// XGBoostJobSpec is a desired state description of the XGBoostJob.
type XGBoostJobSpec struct {
	// RunPolicy encapsulates various runtime policies of the distributed training
	// job, for example how to clean up resources and how long the job can stay
	// active.
	RunPolicy *commmonv1.RunPolicy `json:"runPolicy,omitempty"`

	// XGBReplicaSpec specifies the PyTorch replicas to run.
	XGBReplicaSpec *commmonv1.ReplicaSpec `json:"xgbReplicaSpec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=xgboostjobs

// XGBoostJobList is a list of XGBoostJobs.
type XGBoostJobList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of XGBoostJobs.
	Items []XGBoostJob `json:"items"`
}
