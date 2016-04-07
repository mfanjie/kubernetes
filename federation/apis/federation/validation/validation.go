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
	"k8s.io/kubernetes/federation/apis/federation"
	"k8s.io/kubernetes/pkg/api/validation"
	"k8s.io/kubernetes/pkg/util/validation/field"
)

func ValidateClusterName(name string, prefix bool) (bool, string) {
	return validation.NameIsDNSSubdomain(name, prefix)
}

func ValidateClusterSpec(spec *federation.ClusterSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	// address is required.
	if len(spec.ServerAddressByClientCIDRs) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("serverAddressByClientCIDRs"), ""))
	}
	return allErrs
}

func ValidateCluster(cluster *federation.Cluster) field.ErrorList {
	allErrs := validation.ValidateObjectMeta(&cluster.ObjectMeta, false, ValidateClusterName, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateClusterSpec(&cluster.Spec, field.NewPath("spec"))...)
	return allErrs
}

func ValidateClusterUpdate(cluster, oldCluster *federation.Cluster) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&cluster.ObjectMeta, &oldCluster.ObjectMeta, field.NewPath("metadata"))
	if cluster.Name != oldCluster.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("meta", "name"),
			cluster.Name+" != "+oldCluster.Name, "cannot change cluster name"))
	}
	return allErrs
}

func phaseTransitionAllowed(from, to federation.ClusterPhase) bool {
	validPhaseTransition := map[federation.ClusterPhase][]federation.ClusterPhase{
		federation.ClusterPending:    {federation.ClusterRunning, federation.ClusterOffline, federation.ClusterTerminated},
		federation.ClusterRunning:    {federation.ClusterPending, federation.ClusterOffline, federation.ClusterTerminated},
		federation.ClusterOffline:    {federation.ClusterRunning, federation.ClusterTerminated},
		federation.ClusterTerminated: {},
	}
	for _, allowedPhase := range validPhaseTransition[from] {
		if to == allowedPhase {
			return true
		}
	}
	return false
}
func ValidateClusterStatusUpdate(cluster, oldCluster *federation.Cluster) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&cluster.ObjectMeta, &oldCluster.ObjectMeta, field.NewPath("metadata"))
	if !phaseTransitionAllowed(oldCluster.Status.Phase, cluster.Status.Phase) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("status", "phase"),
			oldCluster.Status.Phase+" => "+cluster.Status.Phase, "cluster phase transition not allowed"))
	}
	return allErrs
}
