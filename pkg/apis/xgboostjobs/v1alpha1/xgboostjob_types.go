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

package v1alpha1

import (
	common "github.com/kubeflow/common/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XGBoostJobSpec defines the desired state of XGBoostJob
type XGBoostJobSpec struct {
	// RunPolicy encapsulates various runtime policies of the distributed training
	RunPolicy *common.RunPolicy `json:",inline"`

	// XGBoostReplicaSpecs specifies the XGBoost replicas to run.
	XGBoostReplicaSpecs map[common.ReplicaType]*common.ReplicaSpec `json:"xgboostReplicaSpec"`
}

// XGBoostJobStatus defines the observed state of XGBoostJob
type XGBoostJobStatus struct {
	//XGBoostJob Status
	common.JobStatus `json:",inline"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XGBoostJob is the Schema for the xgboostjobs API
// +k8s:openapi-gen=true
type XGBoostJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XGBoostJobSpec   `json:"spec,omitempty"`
	Status XGBoostJobStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XGBoostJobList contains a list of XGBoostJob
type XGBoostJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []XGBoostJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&XGBoostJob{}, &XGBoostJobList{})
}
