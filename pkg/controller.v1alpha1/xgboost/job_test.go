package xgboost

import (
	"encoding/json"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	common "github.com/kubeflow/common/operator/v1"
	 "github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
)

func NewXGBoostJob(worker int) *v1alpha1.XGBoostJob {

	job := &v1alpha1.XGBoostJob{
		TypeMeta: metav1.TypeMeta{
			Kind: v1alpha1.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      TestXGBoostJobName,
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
			Template: NewXGBoostReplicaSpecTemplate(),
		}
		job.Spec.XGBoostReplicaSpecs[v1alpha1.XGBoostReplicaTypeWorker] = workerReplicaSpec
	}

	return job
}

func NewXGBoostJobWithMaster(worker int) *v1alpha1.XGBoostJob {
	job := NewXGBoostJob(worker)
	job.Spec.XGBoostReplicaSpecs[v1alpha1.XGBoostReplicaTypeMaster] = &common.ReplicaSpec{
		Template: NewXGBoostReplicaSpecTemplate(),
	}
	return job
}

// ConvertXGBoostJob uses JSON to convert XGBoostJob to Unstructured.
func ConvertXGBoostJobToUnstructured(job *v1alpha1.XGBoostJob) (*unstructured.Unstructured, error) {
	var unstructured unstructured.Unstructured
	b, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &unstructured); err != nil {
		return nil, err
	}
	return &unstructured, nil
}

func NewXGBoostReplicaSpecTemplate() v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{
					Name:  v1alpha1.DefaultContainerName,
					Image: TestImageName,
					Args:  []string{"Fake", "Fake"},
					Ports: []v1.ContainerPort{
						v1.ContainerPort{
							Name:          v1alpha1.DefaultPortName,
							ContainerPort: v1alpha1.DefaultPort,
						},
					},
				},
			},
		},
	}
}