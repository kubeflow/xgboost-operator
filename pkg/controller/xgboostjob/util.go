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

package xgboostjob

import (
	"errors"
	"fmt"
	common "github.com/kubeflow/common/job_controller/api/v1"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboostjob/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

// getClientReaderFromClient try to extract client reader from client, client
// reader reads cluster info from api client.
func getClientReaderFromClient(c client.Client) (client.Reader, error) {
	if dr, err := getDelegatingReader(c); err != nil {
		return nil, err
	} else {
		return dr.ClientReader, nil
	}
}

// getDelegatingReader try to extract DelegatingReader from client.
func getDelegatingReader(c client.Client) (*client.DelegatingReader, error) {
	dc, ok := c.(*client.DelegatingClient)
	if !ok {
		return nil, errors.New("cannot convert from Client to DelegatingClient")
	}
	dr, ok := dc.Reader.(*client.DelegatingReader)
	if !ok {
		return nil, errors.New("cannot convert from DelegatingClient.Reader to Delegating Reader")
	}
	return dr, nil
}

func computeMasterAddr(jobName, rtype, index string) string {
	n := jobName + "-" + rtype + "-" + index
	return strings.Replace(n, "/", "-", -1)
}

// GetPortFromXGBoostJob gets the port of xgboost container.
func GetPortFromXGBoostJob(job *v1alpha1.XGBoostJob, rtype v1alpha1.XGBoostJobReplicaType) (int32, error) {
	containers := job.Spec.XGBReplicaSpecs[common.ReplicaType(rtype)].Template.Spec.Containers
	for _, container := range containers {
		if container.Name == v1alpha1.DefaultContainerName {
			ports := container.Ports
			for _, port := range ports {
				if port.Name == v1alpha1.DefaultContainerPortName {
					return port.ContainerPort, nil
				}
			}
		}
	}
	return -1, fmt.Errorf("failed to found the port")
}

func computeTotalReplicas(obj metav1.Object) int32 {
	job := obj.(*v1alpha1.XGBoostJob)
	jobReplicas := int32(0)

	if job.Spec.XGBReplicaSpecs == nil || len(job.Spec.XGBReplicaSpecs) == 0 {
		return jobReplicas
	}
	for _, r := range job.Spec.XGBReplicaSpecs {
		if r.Replicas == nil {
			continue
		} else {
			jobReplicas += *r.Replicas
		}
	}
	return jobReplicas
}

// FakeWorkQueue implements RateLimitingInterface but actually does nothing.
type FakeWorkQueue struct{}

// Add WorkQueue Add method
func (f *FakeWorkQueue) Add(item interface{}) {}

// Len WorkQueue Len method
func (f *FakeWorkQueue) Len() int { return 0 }

// Get WorkQueue Get method
func (f *FakeWorkQueue) Get() (item interface{}, shutdown bool) { return nil, false }

// Done WorkQueue Done method
func (f *FakeWorkQueue) Done(item interface{}) {}

// ShutDown WorkQueue ShutDown method
func (f *FakeWorkQueue) ShutDown() {}

// ShuttingDown WorkQueue ShuttingDown method
func (f *FakeWorkQueue) ShuttingDown() bool { return true }

// AddAfter WorkQueue AddAfter method
func (f *FakeWorkQueue) AddAfter(item interface{}, duration time.Duration) {}

// AddRateLimited WorkQueue AddRateLimited method
func (f *FakeWorkQueue) AddRateLimited(item interface{}) {}

// Forget WorkQueue Forget method
func (f *FakeWorkQueue) Forget(item interface{}) {}

// NumRequeues WorkQueue NumRequeues method
func (f *FakeWorkQueue) NumRequeues(item interface{}) int { return 0 }
