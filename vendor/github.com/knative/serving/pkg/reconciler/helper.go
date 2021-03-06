/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package reconciler

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Callback func(interface{})

// TODO(mattmoor): Move this into knative/pkg/controller
func EnsureTypeMeta(f Callback, gvk schema.GroupVersionKind) Callback {
	apiVersion, kind := gvk.ToAPIVersionAndKind()

	return func(untyped interface{}) {
		// TODO(mattmoor): Use a kmeta.Accessor based on getObject
		// in pkg/controller.go first, to support deletes.
		typed, ok := untyped.(runtime.Object)
		if !ok {
			return
		}
		// We need to populated TypeMeta, but cannot trample
		// the informer's copy.
		copy := typed.DeepCopyObject()

		accessor, err := meta.TypeAccessor(copy)
		if err != nil {
			return
		}

		accessor.SetAPIVersion(apiVersion)
		accessor.SetKind(kind)
		f(copy)
	}
}
