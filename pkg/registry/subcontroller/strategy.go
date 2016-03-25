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

package subcontroller

import (
	"fmt"
	"reflect"
	"strconv"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/controlplane"
	"k8s.io/kubernetes/pkg/apis/controlplane/validation"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/generic"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/validation/field"
)

type subRcStrategy struct {
	runtime.ObjectTyper
	api.NameGenerator
}

var Strategy = subRcStrategy{api.Scheme, api.SimpleNameGenerator}

func (subRcStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears the status of a replication controller before creation.
func (subRcStrategy) PrepareForCreate(obj runtime.Object) {
	controller := obj.(*controlplane.SubReplicationController)
	controller.Status = api.ReplicationControllerStatus{}

	controller.Generation = 1
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (subRcStrategy) PrepareForUpdate(obj, old runtime.Object) {
	newController := obj.(*controlplane.SubReplicationController)
	oldController := old.(*controlplane.SubReplicationController)
	// update is not allowed to set status
	newController.Status = oldController.Status

	// Any changes to the spec increment the generation number, any changes to the
	// status should reflect the generation number of the corresponding object. We push
	// the burden of managing the status onto the clients because we can't (in general)
	// know here what version of spec the writer of the status has seen. It may seem like
	// we can at first -- since obj contains spec -- but in the future we will probably make
	// status its own object, and even if we don't, writes may be the result of a
	// read-update-write loop, so the contents of spec may not actually be the spec that
	// the controller has *seen*.
	if !reflect.DeepEqual(oldController.Spec, newController.Spec) {
		newController.Generation = oldController.Generation + 1
	}
}

// Validate validates a new replication controller.
func (subRcStrategy) Validate(ctx api.Context, obj runtime.Object) field.ErrorList {
	controller := obj.(*controlplane.SubReplicationController)
	return validation.ValidateSubReplicationController(controller)
}

// Canonicalize normalizes the object after validation.
func (subRcStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for replication controllers; this means a POST is
// needed to create one.
func (subRcStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (subRcStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	validationErrorList := validation.ValidateSubReplicationController(obj.(*controlplane.SubReplicationController))
	updateErrorList := validation.ValidateSubReplicationControllerUpdate(obj.(*controlplane.SubReplicationController), old.(*controlplane.SubReplicationController))
	return append(validationErrorList, updateErrorList...)
}

func (subRcStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// ControllerToSelectableFields returns a field set that represents the object.
func ControllerToSelectableFields(controller *controlplane.SubReplicationController) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(controller.ObjectMeta, true)
	controllerSpecificFieldsSet := fields.Set{
		"status.replicas":   strconv.Itoa(controller.Status.Replicas),
		"spec.cluster.name": controller.Spec.Cluster.Name,
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, controllerSpecificFieldsSet)
}

// MatchController is the filter used by the generic etcd backend to route
// watch events from etcd to clients of the apiserver only interested in specific
// labels/fields.
func MatchController(label labels.Selector, field fields.Selector) generic.Matcher {
	return &generic.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			rc, ok := obj.(*controlplane.SubReplicationController)
			if !ok {
				return nil, nil, fmt.Errorf("Given object is not a replication controller.")
			}
			return labels.Set(rc.ObjectMeta.Labels), ControllerToSelectableFields(rc), nil
		},
	}
}

type rcStatusStrategy struct {
	subRcStrategy
}

var StatusStrategy = rcStatusStrategy{Strategy}

func (rcStatusStrategy) PrepareForUpdate(obj, old runtime.Object) {
	newRc := obj.(*controlplane.SubReplicationController)
	oldRc := old.(*controlplane.SubReplicationController)
	// update is not allowed to set spec
	newRc.Spec = oldRc.Spec
}

func (rcStatusStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateSubReplicationControllerStatusUpdate(obj.(*controlplane.SubReplicationController), old.(*controlplane.SubReplicationController))
}
