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

package v1alpha1

import (
	common "github.com/kubeflow/common/operator/v1"
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
	Status common.JobStatus `json:"status,omitempty"`
}

// XGBoostJobSpec is a desired state description of the XGBoostJob.
type XGBoostJobSpec struct {
	// RunPolicy encapsulates various runtime policies of the distributed training
	// job, for example how to clean up resources and how long the job can stay
	// active.
	RunPolicy *common.RunPolicy `json:"runPolicy,omitempty"`

	// XGBoostReplicaSpecs specifies the XGBoost replicas to run.
	XGBoostReplicaSpecs map[common.ReplicaType]*common.ReplicaSpec `json:"xgboostReplicaSpec"`
}

// XGBoostReplicaType is the type for XGBoostReplica.
type XGBoostReplicaType common.ReplicaType

const (
	// XGBoostReplicaTypeMaster is the type for master worker of distributed XGBoost Job.
	// Rank:0 will be assigned to master worker during AllReduce communication.
	// This is also used as only worker of non-distributed XGBoost Job.
	XGBoostReplicaTypeMaster common.ReplicaType = "Master"


	// XGBoostReplicaTypeWorker is the type for workers of distributed XGBoost Job.
	XGBoostReplicaTypeWorker common.ReplicaType = "Worker"
)

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
