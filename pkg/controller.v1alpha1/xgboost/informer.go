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
	commonutil "github.com/kubeflow/common/util"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/validation"
	log "github.com/sirupsen/logrus"
	metav1unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

const (
	resyncPeriod     = 30 * time.Second
	failedMarshalMsg = "Failed to marshal the object to XGBoostJob: %v"
)

var (
	errGetFromKey    = fmt.Errorf("failed to get XGBoostJob from key")
	errNotExists     = fmt.Errorf("the object is not found")
	errFailedMarshal = fmt.Errorf("failed to marshal the object to XGBoostJob")
)

func (xc *XGBoostController) getXGBoostJobFromName(namespace, name string) (*v1alpha1.XGBoostJob, error) {
	key := fmt.Sprintf("%s/%s", namespace, name)
	return xc.getXGBoostJobFromKey(key)
}

func (xc *XGBoostController) getXGBoostJobFromKey(key string) (*v1alpha1.XGBoostJob, error) {
	// Check if the key exists.
	obj, exists, err := xc.jobInformer.GetIndexer().GetByKey(key)
	logger := commonutil.LoggerForKey(key)
	if err != nil {
		logger.Errorf("Failed to get XGBoostJob '%s' from informer index: %+v", key, err)
		return nil, errGetFromKey
	}
	if !exists {
		// This happens after a job was deleted, but the work queue still had an entry for it.
		return nil, errNotExists
	}

	return jobFromUnstructured(obj)
}

func jobFromUnstructured(obj interface{}) (*v1alpha1.XGBoostJob, error) {
	// Check if the spec is valid.
	un, ok := obj.(*metav1unstructured.Unstructured)
	if !ok {
		log.Errorf("The object in index is not an unstructured; %+v", obj)
		return nil, errGetFromKey
	}
	var job v1alpha1.XGBoostJob
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(un.Object, &job)
	logger := commonutil.LoggerForUnstructured(un, v1alpha1.Kind)
	if err != nil {
		logger.Errorf(failedMarshalMsg, err)
		return nil, errFailedMarshal
	}

	err = validation.ValidateAlphaOneXGBoostJobSpec(&job.Spec)
	if err != nil {
		logger.Errorf(failedMarshalMsg, err)
		return nil, errFailedMarshal
	}
	return &job, nil
}

func unstructuredFromJob(obj interface{}, job *v1alpha1.XGBoostJob) error {
	un, ok := obj.(*metav1unstructured.Unstructured)
	logger := commonutil.LoggerForJob(job)
	if !ok {
		logger.Warn("The object in index isn't type Unstructured")
		return errGetFromKey
	}

	var err error
	un.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(job)
	if err != nil {
		logger.Error("The XGBoostJob convert failed")
		return err
	}
	return nil
}
