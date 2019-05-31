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
	"sigs.k8s.io/controller-runtime/pkg/client"
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
