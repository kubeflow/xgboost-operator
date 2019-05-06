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
	"strings"

	common "github.com/kubeflow/common/operator/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Int32 is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it.
func Int32(v int32) *int32 {
	return &v
}

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// setDefaultPort sets the default ports for container.
func setDefaultPort(spec *v1.PodSpec) {
	index := 0
	for i, container := range spec.Containers {
		if container.Name == DefaultContainerName {
			index = i
			break
		}
	}

	hasJobPort := false
	for _, port := range spec.Containers[index].Ports {
		if port.Name == DefaultPortName {
			hasJobPort = true
			break
		}
	}
	if !hasJobPort {
		spec.Containers[index].Ports = append(spec.Containers[index].Ports, v1.ContainerPort{
			Name:          DefaultPortName,
			ContainerPort: DefaultPort,
		})
	}
}

func setDefaultReplicas(spec *common.ReplicaSpec) {
	if spec.Replicas == nil {
		spec.Replicas = Int32(1)
	}
	if spec.RestartPolicy == "" {
		spec.RestartPolicy = DefaultRestartPolicy
	}
}

// setTypeNamesToCamelCase sets the name of all replica types from any case to correct case.
func setTypeNamesToCamelCase(xgboostJob *XGBoostJob) {
	setTypeNameToCamelCase(xgboostJob, XGBoostReplicaTypeWorker)
	setTypeNameToCamelCase(xgboostJob, XGBoostReplicaTypeMaster)
}

// setTypeNameToCamelCase sets the name of the replica type from any case to correct case.
func setTypeNameToCamelCase(xgboostJob *XGBoostJob, typ XGBoostReplicaType) {
	for t := range xgboostJob.Spec.XGBoostReplicaSpecs {
		if strings.EqualFold(string(t), string(typ)) && t != typ {
			spec := xgboostJob.Spec.XGBoostReplicaSpecs[t]
			delete(xgboostJob.Spec.XGBoostReplicaSpecs, t)
			xgboostJob.Spec.XGBoostReplicaSpecs[typ] = spec
			return
		}
	}
}

// SetDefaults_XGBoostJob sets any unspecified values to defaults.
func SetDefaults_XGBoostJob(xgboostjob *XGBoostJob) {
	// Set default cleanpod policy to Running.
	if xgboostjob.Spec.RunPolicy.CleanPodPolicy == nil {
		running := common.CleanPodPolicyRunning
		xgboostjob.Spec.RunPolicy.CleanPodPolicy = &running
	}

	// Update the key of XGBoostReplicaSpecs to camel case.
	setTypeNamesToCamelCase(xgboostjob)

	for _, spec := range xgboostjob.Spec.XGBoostReplicaSpecs {
		// Set default replicas to 1.
		setDefaultReplicas(spec)
		// Set default port to the container.
		setDefaultPort(&spec.Template.Spec)
	}
}
