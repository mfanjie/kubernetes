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

package validation

import (
	"k8s.io/kubernetes/pkg/api/validation"
	"k8s.io/kubernetes/pkg/apis/controlplane"
	"k8s.io/kubernetes/pkg/util/validation/field"
)

func ValidateClusterName(name string, prefix bool) (bool, string) {
	return validation.NameIsDNSSubdomain(name, prefix)
}

func ValidateCluster(cluster *controlplane.Cluster) field.ErrorList {
	allErrs := validation.ValidateObjectMeta(&cluster.ObjectMeta, false, ValidateClusterName, field.NewPath("metadata"))

	// Only validate spec. All status fields are optional and can be updated later.
	// address is required.
	if len(cluster.Spec.Address.Url) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "address.url"), ""))
	}

	return allErrs
}

func ValidateClusterUpdate(cluster, oldCluster *controlplane.Cluster) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&cluster.ObjectMeta, &oldCluster.ObjectMeta, field.NewPath("metadata"))

	return allErrs
}

func ValidateClusterSelectorSpec(cluster *controlplane.ClusterSelectionSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if ok, qualifier := ValidateClusterName(cluster.Name, false); !ok {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), cluster.Name, qualifier))
	}
	allErrs = append(allErrs, validation.ValidateLabels(cluster.Selector, fldPath.Child("selector"))...)
	return allErrs
}

func ValidateSubReplicationController(controller *controlplane.SubReplicationController) field.ErrorList {
	allErrs := validation.ValidateObjectMeta(&controller.ObjectMeta, true, validation.ValidateReplicationControllerName, field.NewPath("metadata"))
	allErrs = append(allErrs, validation.ValidateReplicationControllerSpec(&controller.Spec.ReplicationControllerSpec, field.NewPath("spec"))...)
	allErrs = append(allErrs, ValidateClusterSelectorSpec(&controller.Spec.Cluster, field.NewPath("spec", "cluster"))...)
	return allErrs
}

func ValidateSubReplicationControllerUpdate(controller, oldController *controlplane.SubReplicationController) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&controller.ObjectMeta, &oldController.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, validation.ValidateReplicationControllerSpec(&controller.Spec.ReplicationControllerSpec, field.NewPath("spec"))...)
	allErrs = append(allErrs, ValidateClusterSelectorSpec(&controller.Spec.Cluster, field.NewPath("spec", "cluster"))...)
	return allErrs
}

func ValidateSubReplicationControllerStatusUpdate(controller, oldController *controlplane.SubReplicationController) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&controller.ObjectMeta, &oldController.ObjectMeta, field.NewPath("metadata"))
	statusPath := field.NewPath("status")
	allErrs = append(allErrs, validation.ValidateNonnegativeField(int64(controller.Status.Replicas), statusPath.Child("replicas"))...)
	allErrs = append(allErrs, validation.ValidateNonnegativeField(int64(controller.Status.FullyLabeledReplicas), statusPath.Child("fullyLabeledReplicas"))...)
	allErrs = append(allErrs, validation.ValidateNonnegativeField(int64(controller.Status.ObservedGeneration), statusPath.Child("observedGeneration"))...)
	return allErrs
}
