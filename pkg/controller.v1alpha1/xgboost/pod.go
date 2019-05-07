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
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

func (xc *XGBoostController) SetClusterSpec(job interface{}, podTemplate *corev1.PodTemplateSpec,
	rtype, index string) error {
	// TODO: Implement this method
	return nil
}

func (xc *XGBoostController) CreateService(job interface{}, service *corev1.Service) error {
	xgbJob := job.(*v1alpha1.XGBoostJob)
	controllerRef := xc.GenOwnerReference(xgbJob)
	return xc.ServiceControl.CreateServicesWithControllerRef(xgbJob.Namespace, service, xgbJob, controllerRef)
}

func (xc *XGBoostController) DeleteService(job interface{}, name string, namespace string) error {
	log.Info("Deleting service " + name)
	return xc.ServiceControl.DeleteService(namespace, name, job.(*v1alpha1.XGBoostJob))
}

func (xc *XGBoostController) CreatePod(job interface{}, podTemplate *corev1.PodTemplateSpec) error {
	xgbJob := job.(*v1alpha1.XGBoostJob)
	controllerRef := xc.GenOwnerReference(xgbJob)
	return xc.PodControl.CreatePodsWithControllerRef(xgbJob.Namespace, podTemplate, xgbJob, controllerRef)
}

func (xc *XGBoostController) DeletePod(job interface{}, pod *corev1.Pod) error {
	log.Info("Deleting pod " + pod.Name)
	return xc.PodControl.DeletePod(pod.Namespace, pod.Name, job.(*v1alpha1.XGBoostJob))
}
