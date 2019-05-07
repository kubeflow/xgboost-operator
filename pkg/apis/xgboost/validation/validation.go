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

package validation

import (
	"fmt"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
	log "github.com/sirupsen/logrus"
)

func ValidateAlphaOneXGBoostJobSpec(c *v1alpha1.XGBoostJobSpec) error {
	if c.XGBoostReplicaSpecs == nil {
		return fmt.Errorf("XGBoostJobSpec is not valid")
	}
	masterExists := false
	for rType, value := range c.XGBoostReplicaSpecs {
		if value == nil || len(value.Template.Spec.Containers) == 0 {
			return fmt.Errorf("XGBoostJobSpec is not valid")
		}
		// Make sure the replica type is valid.
		validReplicaTypes := []v1alpha1.XGBoostReplicaType{v1alpha1.XGBoostReplicaTypeMaster, v1alpha1.XGBoostReplicaTypeWorker}

		isValidReplicaType := false
		for _, t := range validReplicaTypes {
			if t == rType {
				isValidReplicaType = true
				break
			}
		}

		if !isValidReplicaType {
			return fmt.Errorf("XGBoostReplicaType is %v but must be one of %v", rType, validReplicaTypes)
		}

		//Make sure the image is defined in the container
		defaultContainerPresent := false
		for _, container := range value.Template.Spec.Containers {
			if container.Image == "" {
				log.Warn("Image is undefined in the container")
				return fmt.Errorf("XGBoostJobSpec is not valid")
			}
			if container.Name == v1alpha1.DefaultContainerName {
				defaultContainerPresent = true
			}
		}
		//Make sure there has at least one container named "pytorch"
		if !defaultContainerPresent {
			log.Warnf("There is no container named pytorch in %v", rType)
			return fmt.Errorf("XGBoostJobSpec is not valid")
		}
		if rType == v1alpha1.XGBoostReplicaTypeMaster {
			masterExists = true
			if value.Replicas != nil && int(*value.Replicas) != 1 {
				log.Warnf("There must be only 1 master replica")
				return fmt.Errorf("XGBoostJobSpec is not valid")
			}
		}

	}

	if !masterExists {
		log.Warnf("Master ReplicaSpec must be present")
		return fmt.Errorf("XGBoostJobSpec is not valid")
	}
	return nil
}
