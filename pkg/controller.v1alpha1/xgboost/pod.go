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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/kubeflow/common/util/k8sutil"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	common "github.com/kubeflow/common/operator/v1"
	"github.com/kubeflow/tf-operator/pkg/common/jobcontroller"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	pylogger "github.com/kubeflow/tf-operator/pkg/logger"

)

const (

	// gang scheduler name.
	gangSchedulerName = "kube-batch"

	// podTemplateRestartPolicyReason is the warning reason when the restart
	// policy is set in pod template.
	podTemplateRestartPolicyReason = "SettedPodTemplateRestartPolicy"

	exitedWithCodeReason           = "ExitedWithCode"

	// podTemplateSchedulerNameReason is the warning reason when other scheduler name is set
	// in pod templates with gang-scheduling enabled
	podTemplateSchedulerNameReason = "SettedPodTemplateSchedulerName"
)

var (
	errPortNotFound = fmt.Errorf("failed to found the port")
)


func (xc *XGBoostController) CreateService(job interface{}, service *corev1.Service) error {
	xgbJob := job.(*v1alpha1.XGBoostJob)
	controllerRef := xc.GenOwnerReference(xgbJob)
	return xc.ServiceControl.CreateServicesWithControllerRef(xgbJob.Namespace, service, xgbJob, controllerRef)
}

func (xc *XGBoostController) DeleteService(job interface{}, name string, namespace string) error {
	log.Info("Deleting service " + name)
	return xc.ServiceControl.DeleteService(namespace, name, job.(*v1alpha1.XGBoostJob))
}

func (xc *XGBoostController) CreatePod(job *v1alpha1.XGBoostJob, rtype v1alpha1.XGBoostReplicaType, index string, spec *common.ReplicaSpec, masterRole bool) error {

	rt := strings.ToLower(string(rtype))
	jobKey, err := KeyFunc(job)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for job object %#v: %v", job, err))
		return err
	}
	expectationPodsKey := jobcontroller.GenExpectationPodsKey(jobKey, rt)
	err = xc.Expectations.ExpectCreations(expectationPodsKey, 1)
	if err != nil {
		return err
	}
	logger := pylogger.LoggerForReplica(job, rt)
	// Create OwnerReference.
	controllerRef := xc.GenOwnerReference(job)

	// Set type and index for the worker.
	labels := xc.GenLabels(job.Name)
	labels[replicaTypeLabel] = rt
	labels[replicaIndexLabel] = index

	if masterRole {
		labels[labelXGBoostJobRole] = "master"
	}
	podTemplate := spec.Template.DeepCopy()
	totalReplicas := k8sutil.GetTotalReplicas(job.Spec.XGBoostReplicaSpecs)
	// Set name for the template.
	podTemplate.Name = jobcontroller.GenGeneralName(job.Name, rt, index)

	if podTemplate.Labels == nil {
		podTemplate.Labels = make(map[string]string)
	}

	for key, value := range labels {
		podTemplate.Labels[key] = value
	}

	if err := setClusterSpec(podTemplate, job, totalReplicas, index, rtype); err != nil {
		return err
	}

	// Submit a warning event if the user specifies restart policy for
	// the pod template. We recommend to set it from the replica level.
	if podTemplate.Spec.RestartPolicy != corev1.RestartPolicy("") {
		errMsg := "Restart policy in pod template will be overwritten by restart policy in replica spec"
		logger.Warning(errMsg)
		xc.Recorder.Event(job, corev1.EventTypeWarning, podTemplateRestartPolicyReason, errMsg)
	}
	setRestartPolicy(podTemplate, spec)

	// if gang-scheduling is enabled:
	// 1. if user has specified other scheduler, we report a warning without overriding any fields.
	// 2. if no SchedulerName is set for pods, then we set the SchedulerName to "kube-batch".
	if xc.Config.EnableGangScheduling {
		if isNonGangSchedulerSet(job) {
			errMsg := "Another scheduler is specified when gang-scheduling is enabled and it will not be overwritten"
			logger.Warning(errMsg)
			xc.Recorder.Event(job, corev1.EventTypeWarning, podTemplateSchedulerNameReason, errMsg)
		} else {
			podTemplate.Spec.SchedulerName = gangSchedulerName
		}
	}

	err = xc.PodControl.CreatePodsWithControllerRef(job.Namespace, podTemplate, job, controllerRef)
	if err != nil && k8serrors.IsTimeout(err) {
		// Pod is created but its initialization has timed out.
		// If the initialization is successful eventually, the
		// controller will observe the creation via the informer.
		// If the initialization fails, or if the pod keeps
		// uninitialized for a long time, the informer will not
		// receive any update, and the controller will create a new
		// pod when the expectation expires.
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func (xc *XGBoostController) DeletePod(job interface{}, pod *corev1.Pod) error {
	log.Info("Deleting pod " + pod.Name)
	return xc.PodControl.DeletePod(pod.Namespace, pod.Name, job.(*v1alpha1.XGBoostJob))
}

func setClusterSpec(podTemplateSpec *corev1.PodTemplateSpec, job *v1alpha1.XGBoostJob, totalReplicas int32, index string, rtype v1alpha1.XGBoostReplicaType) error {
	rank, err := strconv.Atoi(index)
	if err != nil {
		return err
	}

	masterPort, err := GetPortFromXGBoostJob(job, v1alpha1.XGBoostReplicaTypeMaster)
	if err != nil {
		return err
	}

	masterAddr := jobcontroller.GenGeneralName(job.Name, strings.ToLower(string(v1alpha1.XGBoostReplicaTypeMaster)), strconv.Itoa(0))
	if rtype == v1alpha1.XGBoostReplicaTypeMaster {
		if rank != 0 {
			return errors.New("invalid config: There should be only a single master with index=0")
		}
		///TODO: add the container IP later
		masterAddr = "localhost"
	} else {
		rank = rank + 1
	}

	for i := range podTemplateSpec.Spec.Containers {
		if len(podTemplateSpec.Spec.Containers[i].Env) == 0 {
			podTemplateSpec.Spec.Containers[i].Env = make([]corev1.EnvVar, 0)
		}
		podTemplateSpec.Spec.Containers[i].Env = append(podTemplateSpec.Spec.Containers[i].Env, corev1.EnvVar{
			Name:  "MASTER_PORT",
			Value: strconv.Itoa(int(masterPort)),
		})
		podTemplateSpec.Spec.Containers[i].Env = append(podTemplateSpec.Spec.Containers[i].Env, corev1.EnvVar{
			Name:  "MASTER_ADDR",
			Value: masterAddr,
		})
		podTemplateSpec.Spec.Containers[i].Env = append(podTemplateSpec.Spec.Containers[i].Env, corev1.EnvVar{
			Name:  "WORLD_SIZE",
			Value: strconv.Itoa(int(totalReplicas)),
		})
		podTemplateSpec.Spec.Containers[i].Env = append(podTemplateSpec.Spec.Containers[i].Env, corev1.EnvVar{
			Name:  "RANK",
			Value: strconv.Itoa(rank),
		})
		podTemplateSpec.Spec.Containers[i].Env = append(podTemplateSpec.Spec.Containers[i].Env, corev1.EnvVar{
			Name:  "PYTHONUNBUFFERED",
			Value: "0",
		})
	}
	return nil
}

func setRestartPolicy(podTemplateSpec *corev1.PodTemplateSpec, spec *common.ReplicaSpec) {
	if spec.RestartPolicy == common.RestartPolicyExitCode {
		podTemplateSpec.Spec.RestartPolicy = corev1.RestartPolicyNever
	} else {
		podTemplateSpec.Spec.RestartPolicy = corev1.RestartPolicy(spec.RestartPolicy)
	}
}

func isNonGangSchedulerSet(xgjob *v1alpha1.XGBoostJob) bool {
	for _, spec := range xgjob.Spec.XGBoostReplicaSpecs {
		if spec.Template.Spec.SchedulerName != "" && spec.Template.Spec.SchedulerName != gangSchedulerName {
			return true
		}
	}
	return false
}

// GetPortFromPyTorchJob gets the port of pytorch container.
func GetPortFromXGBoostJob(job *v1alpha1.XGBoostJob, rtype v1alpha1.XGBoostReplicaType) (int32, error) {
	containers := job.Spec.XGBoostReplicaSpecs[rtype].Template.Spec.Containers
	for _, container := range containers {
		if container.Name == v1alpha1.DefaultContainerName {
			ports := container.Ports
			for _, port := range ports {
				if port.Name == v1alpha1.DefaultPortName {
					return port.ContainerPort, nil
				}
			}
		}
	}
	return -1, errPortNotFound
}


