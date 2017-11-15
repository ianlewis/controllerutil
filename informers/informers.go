// Copyright 2017 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package informers

import (
	"context"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// SharedInformers allows controllers to register informers that watch the Kubernetes API for changes in various API types and maintain a local cache of those objects.
type SharedInformers interface {
	// InformerFor registers an informer for the given type. If an informer for the given type has not been
	// registered before, the given constructor function is called to construct the informer.
	InformerFor(string, metav1.GroupVersionKind, func() cache.SharedIndexInformer) cache.SharedIndexInformer
	// Run starts all registered informers.
	Run(context.Context) error
}

type sharedInformers struct {
	lock             sync.Mutex
	informers        map[string]map[metav1.GroupVersionKind]cache.SharedIndexInformer
	startedInformers map[string]map[metav1.GroupVersionKind]bool
}

// NewSharedInformers creates an empty SharedInformers instance.
func NewSharedInformers() SharedInformers {
	return &sharedInformers{
		informers:        make(map[string]map[metav1.GroupVersionKind]cache.SharedIndexInformer),
		startedInformers: make(map[string]map[metav1.GroupVersionKind]bool),
	}
}

// InformerFor registers an informer for the given type. If an informer for the given type has not been registered before, the given constructor function is called to construct the informer.
func (s *sharedInformers) InformerFor(namespace string, gvk metav1.GroupVersionKind, f func() cache.SharedIndexInformer) cache.SharedIndexInformer {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.informers == nil {
		s.informers = make(map[string]map[metav1.GroupVersionKind]cache.SharedIndexInformer)
	}
	informerMap, nsExists := s.informers[namespace]
	if !nsExists {
		i := f()
		s.informers[namespace] = map[metav1.GroupVersionKind]cache.SharedIndexInformer{
			gvk: i,
		}
		return i
	}
	informer, exists := informerMap[gvk]
	if !exists {
		i := f()
		informerMap[gvk] = i
		return i
	}
	return informer
}

// Run runs all registered informers.
func (f *sharedInformers) Run(ctx context.Context) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	for namespace, informerMap := range f.informers {
		startedMap, exists := f.startedInformers[namespace]
		if !exists {
			startedMap = make(map[metav1.GroupVersionKind]bool)
			f.startedInformers[namespace] = startedMap
		}
		for informerType, informer := range informerMap {
			if !startedMap[informerType] {
				go informer.Run(ctx.Done())
				startedMap[informerType] = true
			}
		}
	}

	<-ctx.Done()

	return nil
}
