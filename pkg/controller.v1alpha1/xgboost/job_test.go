package xgboost

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha1 "github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1beta2"
)

func NewXGBoostJob(worker int) *v1alpha1.XGBoostJob {

	job := &v1alpha1.XGBoostJob{
		TypeMeta: metav1.TypeMeta{
			Kind: v1alpha1.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      TestPyTorchJobName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: v1alpha1.XGBoostJobSpec{
			XGBoostReplicaSpecs: make(map[v1alpha1.XGBoostReplicaType]*common.ReplicaSpec),
		},
	}

	if worker > 0 {
		worker := int32(worker)
		workerReplicaSpec := &common.ReplicaSpec{
			Replicas: &worker,
			Template: NewPyTorchReplicaSpecTemplate(),
		}
		job.Spec.XGBoostReplicaSpecs[v1alpha1.XGBoostReplicaTypeWorker] = workerReplicaSpec
	}

	return job
}