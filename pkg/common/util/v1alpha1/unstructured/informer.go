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

package unstructured

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"time"

	informer "github.com/kubeflow/xgboost-operator/pkg/client/informers/externalversions/xgboost/v1alpha1"
	lister "github.com/kubeflow/xgboost-operator/pkg/client/listers/xgboost/v1alpha1"
)

type UnstructuredInformer struct {
	informer cache.SharedIndexInformer
}

func NewXGBoostJobInformer(resource schema.GroupVersionResource, client dynamic.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) informer.XGBoostJobInformer {
	return &UnstructuredInformer{
		informer: newUnstructuredInformer(resource, client, namespace, resyncPeriod, indexers),
	}
}

func (f *UnstructuredInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

func (f *UnstructuredInformer) Lister() lister.XGBoostJobLister {
	return lister.NewXGBoostJobLister(f.Informer().GetIndexer())
}


// newUnstructuredInformer constructs a new informer for Unstructured type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func newUnstructuredInformer(resource schema.GroupVersionResource, client dynamic.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return newFilteredUnstructuredInformer(resource, client, namespace, resyncPeriod, indexers)
}

// newFilteredUnstructuredInformer constructs a new informer for Unstructured type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func newFilteredUnstructuredInformer(resource schema.GroupVersionResource, client dynamic.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return client.Resource(resource).Namespace(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return client.Resource(resource).Namespace(namespace).Watch(options)
			},
		},
		&unstructured.Unstructured{},
		resyncPeriod,
		indexers,
	)
}
