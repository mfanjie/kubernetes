/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package v1alpha1

import (
	extensionsv1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) {
	scheme.AddDefaultingFuncs(
		func(obj *Cluster) {
			if obj.Status.Phase == "" {
				obj.Status.Phase = ClusterPending
			}
		},
		func(obj *SubReplicaSet) {
			labels := obj.Spec.Template.Labels

			// TODO: support templates defined elsewhere when we support them in the API
			if labels != nil {
				if obj.Spec.Selector == nil {
					obj.Spec.Selector = &extensionsv1.LabelSelector{
						MatchLabels: labels,
					}
				}
				if len(obj.Labels) == 0 {
					obj.Labels = labels
				}
			}
			if obj.Spec.Replicas == nil {
				obj.Spec.Replicas = new(int32)
				*obj.Spec.Replicas = 1
			}
		},
	)
}
