package job_controller

import (
	"errors"
	"fmt"
	"strings"

	apiv1 "github.com/kubeflow/common/job_controller/api/v1"
	"github.com/kubernetes-sigs/kube-batch/pkg/apis/scheduling/v1alpha1"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/policy/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/cache"
)

var (
	// KeyFunc is the short name to DeletionHandlingMetaNamespaceKeyFunc.
	// IndexerInformer uses a delta queue, therefore for deletes we have to use this
	// key function but it should be just fine for non delete events.
	KeyFunc = cache.DeletionHandlingMetaNamespaceKeyFunc
)

type JobController apiv1.JobController

func (jc *JobController) GenOwnerReference(obj metav1.Object) *metav1.OwnerReference {
	boolPtr := func(b bool) *bool { return &b }
	controllerRef := &metav1.OwnerReference{
		APIVersion:         jc.Controller.GetAPIGroupVersion().String(),
		Kind:               jc.Controller.GetAPIGroupVersionKind().Kind,
		Name:               obj.GetName(),
		UID:                obj.GetUID(),
		BlockOwnerDeletion: boolPtr(true),
		Controller:         boolPtr(true),
	}

	return controllerRef
}

func (jc *JobController) GenLabels(jobName string) map[string]string {
	labelGroupName := apiv1.GroupNameLabel
	labelJobName := apiv1.JobNameLabel
	groupName := jc.Controller.GetGroupNameLabelValue()
	return map[string]string{
		labelGroupName: groupName,
		labelJobName:   strings.Replace(jobName, "/", "-", -1),
	}
}

func (jc *JobController) SyncPodGroup(job metav1.Object, minAvailableReplicas int32) (*v1alpha1.PodGroup, error) {

	kubeBatchClientInterface := jc.KubeBatchClientSet
	// Check whether podGroup exists or not
	podGroup, err := kubeBatchClientInterface.SchedulingV1alpha1().PodGroups(job.GetNamespace()).Get(job.GetName(), metav1.GetOptions{})
	if err == nil {
		return podGroup, nil
	}

	// create podGroup for gang scheduling by kube-batch
	minAvailable := intstr.FromInt(int(minAvailableReplicas))
	createPodGroup := &v1alpha1.PodGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: job.GetName(),
			OwnerReferences: []metav1.OwnerReference{
				*jc.GenOwnerReference(job),
			},
		},
		Spec: v1alpha1.PodGroupSpec{
			MinMember: minAvailable.IntVal,
		},
	}
	return kubeBatchClientInterface.SchedulingV1alpha1().PodGroups(job.GetNamespace()).Create(createPodGroup)
}

// SyncPdb will create a PDB for gang scheduling by kube-batch.
func (jc *JobController) SyncPdb(job metav1.Object, minAvailableReplicas int32) (*v1beta1.PodDisruptionBudget, error) {
	labelJobName := apiv1.JobNameLabel

	// Check the pdb exist or not
	pdb, err := jc.KubeClientSet.PolicyV1beta1().PodDisruptionBudgets(job.GetNamespace()).Get(job.GetName(), metav1.GetOptions{})
	if err == nil || !k8serrors.IsNotFound(err) {
		if err == nil {
			err = errors.New(string(metav1.StatusReasonAlreadyExists))
		}
		return pdb, err
	}

	// Create pdb for gang scheduling by kube-batch
	minAvailable := intstr.FromInt(int(minAvailableReplicas))
	createPdb := &v1beta1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name: job.GetName(),
			OwnerReferences: []metav1.OwnerReference{
				*jc.GenOwnerReference(job),
			},
		},
		Spec: v1beta1.PodDisruptionBudgetSpec{
			MinAvailable: &minAvailable,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					labelJobName: job.GetName(),
				},
			},
		},
	}
	return jc.KubeClientSet.PolicyV1beta1().PodDisruptionBudgets(job.GetNamespace()).Create(createPdb)
}

func (jc *JobController) DeletePodGroup(job metav1.Object) error {
	kubeBatchClientInterface := jc.KubeBatchClientSet

	//check whether podGroup exists or not
	_, err := kubeBatchClientInterface.SchedulingV1alpha1().PodGroups(job.GetNamespace()).Get(job.GetName(), metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		return nil
	}

	log.Infof("Deleting PodGroup %s", job.GetName())

	//delete podGroup
	err = kubeBatchClientInterface.SchedulingV1alpha1().PodGroups(job.GetNamespace()).Delete(job.GetName(), &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("unable to delete PodGroup: %v", err)
	}
	return nil
}

func (jc *JobController) DeletePdb(job metav1.Object) error {

	// Check the pdb exist or not
	_, err := jc.KubeClientSet.PolicyV1beta1().PodDisruptionBudgets(job.GetNamespace()).Get(job.GetName(), metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		return nil
	}

	msg := fmt.Sprintf("Deleting pdb %s", job.GetName())
	log.Info(msg)

	if err := jc.KubeClientSet.PolicyV1beta1().PodDisruptionBudgets(job.GetNamespace()).Delete(job.GetName(), &metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("unable to delete pdb: %v", err)
	}
	return nil
}

// resolveControllerRef returns the job referenced by a ControllerRef,
// or nil if the ControllerRef could not be resolved to a matching job
// of the correct Kind.
func (jc *JobController) resolveControllerRef(namespace string, controllerRef *metav1.OwnerReference) metav1.Object {
	// We can't look up by UID, so look up by Name and then verify UID.
	// Don't even try to look up by Name if it's the wrong Kind.
	if controllerRef.Kind != jc.Controller.GetAPIGroupVersionKind().Kind {
		return nil
	}
	job, err := jc.Controller.GetJobFromInformerCache(namespace, controllerRef.Name)
	if err != nil {
		return nil
	}
	if job.GetUID() != controllerRef.UID {
		// The controller we found with this Name is not the same one that the
		// ControllerRef points to.
		return nil
	}
	return job
}
