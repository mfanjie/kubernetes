/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

// DO NOT EDIT. THIS FILE IS AUTO-GENERATED BY $KUBEROOT/hack/update-generated-conversions.sh

package v1alpha1

import (
	reflect "reflect"

	api "k8s.io/kubernetes/pkg/api"
	resource "k8s.io/kubernetes/pkg/api/resource"
	controlplane "k8s.io/kubernetes/pkg/apis/controlplane"
	conversion "k8s.io/kubernetes/pkg/conversion"
)

func autoConvert_controlplane_Cluster_To_v1alpha1_Cluster(in *controlplane.Cluster, out *Cluster, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*controlplane.Cluster))(in)
	}
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := Convert_controlplane_ClusterSpec_To_v1alpha1_ClusterSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_controlplane_ClusterStatus_To_v1alpha1_ClusterStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

func Convert_controlplane_Cluster_To_v1alpha1_Cluster(in *controlplane.Cluster, out *Cluster, s conversion.Scope) error {
	return autoConvert_controlplane_Cluster_To_v1alpha1_Cluster(in, out, s)
}

func autoConvert_controlplane_ClusterAddress_To_v1alpha1_ClusterAddress(in *controlplane.ClusterAddress, out *ClusterAddress, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*controlplane.ClusterAddress))(in)
	}
	out.Url = in.Url
	return nil
}

func Convert_controlplane_ClusterAddress_To_v1alpha1_ClusterAddress(in *controlplane.ClusterAddress, out *ClusterAddress, s conversion.Scope) error {
	return autoConvert_controlplane_ClusterAddress_To_v1alpha1_ClusterAddress(in, out, s)
}

func autoConvert_controlplane_ClusterList_To_v1alpha1_ClusterList(in *controlplane.ClusterList, out *ClusterList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*controlplane.ClusterList))(in)
	}
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := api.Convert_unversioned_ListMeta_To_unversioned_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]Cluster, len(in.Items))
		for i := range in.Items {
			if err := Convert_controlplane_Cluster_To_v1alpha1_Cluster(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func Convert_controlplane_ClusterList_To_v1alpha1_ClusterList(in *controlplane.ClusterList, out *ClusterList, s conversion.Scope) error {
	return autoConvert_controlplane_ClusterList_To_v1alpha1_ClusterList(in, out, s)
}

func autoConvert_controlplane_ClusterSelectionSpec_To_v1alpha1_ClusterSelectionSpec(in *controlplane.ClusterSelectionSpec, out *ClusterSelectionSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*controlplane.ClusterSelectionSpec))(in)
	}
	out.Name = in.Name
	if in.Selector != nil {
		out.Selector = make(map[string]string)
		for key, val := range in.Selector {
			out.Selector[key] = val
		}
	} else {
		out.Selector = nil
	}
	return nil
}

func Convert_controlplane_ClusterSelectionSpec_To_v1alpha1_ClusterSelectionSpec(in *controlplane.ClusterSelectionSpec, out *ClusterSelectionSpec, s conversion.Scope) error {
	return autoConvert_controlplane_ClusterSelectionSpec_To_v1alpha1_ClusterSelectionSpec(in, out, s)
}

func autoConvert_controlplane_ClusterSpec_To_v1alpha1_ClusterSpec(in *controlplane.ClusterSpec, out *ClusterSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*controlplane.ClusterSpec))(in)
	}
	if err := Convert_controlplane_ClusterAddress_To_v1alpha1_ClusterAddress(&in.Address, &out.Address, s); err != nil {
		return err
	}
	out.Credential = in.Credential
	return nil
}

func Convert_controlplane_ClusterSpec_To_v1alpha1_ClusterSpec(in *controlplane.ClusterSpec, out *ClusterSpec, s conversion.Scope) error {
	return autoConvert_controlplane_ClusterSpec_To_v1alpha1_ClusterSpec(in, out, s)
}

func autoConvert_controlplane_ClusterStatus_To_v1alpha1_ClusterStatus(in *controlplane.ClusterStatus, out *ClusterStatus, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*controlplane.ClusterStatus))(in)
	}
	out.Phase = ClusterPhase(in.Phase)
	if in.Capacity != nil {
		out.Capacity = make(api.ResourceList)
		for key, val := range in.Capacity {
			newVal := resource.Quantity{}
			if err := s.Convert(&val, &newVal, 0); err != nil {
				return err
			}
			out.Capacity[key] = newVal
		}
	} else {
		out.Capacity = nil
	}
	out.ClusterMeta = in.ClusterMeta
	return nil
}

func Convert_controlplane_ClusterStatus_To_v1alpha1_ClusterStatus(in *controlplane.ClusterStatus, out *ClusterStatus, s conversion.Scope) error {
	return autoConvert_controlplane_ClusterStatus_To_v1alpha1_ClusterStatus(in, out, s)
}

func autoConvert_controlplane_SubReplicationController_To_v1alpha1_SubReplicationController(in *controlplane.SubReplicationController, out *SubReplicationController, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*controlplane.SubReplicationController))(in)
	}
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := Convert_controlplane_SubReplicationControllerSpec_To_v1alpha1_SubReplicationControllerSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := s.Convert(&in.Status, &out.Status, 0); err != nil {
		return err
	}
	return nil
}

func Convert_controlplane_SubReplicationController_To_v1alpha1_SubReplicationController(in *controlplane.SubReplicationController, out *SubReplicationController, s conversion.Scope) error {
	return autoConvert_controlplane_SubReplicationController_To_v1alpha1_SubReplicationController(in, out, s)
}

func autoConvert_controlplane_SubReplicationControllerList_To_v1alpha1_SubReplicationControllerList(in *controlplane.SubReplicationControllerList, out *SubReplicationControllerList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*controlplane.SubReplicationControllerList))(in)
	}
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := api.Convert_unversioned_ListMeta_To_unversioned_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]SubReplicationController, len(in.Items))
		for i := range in.Items {
			if err := Convert_controlplane_SubReplicationController_To_v1alpha1_SubReplicationController(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func Convert_controlplane_SubReplicationControllerList_To_v1alpha1_SubReplicationControllerList(in *controlplane.SubReplicationControllerList, out *SubReplicationControllerList, s conversion.Scope) error {
	return autoConvert_controlplane_SubReplicationControllerList_To_v1alpha1_SubReplicationControllerList(in, out, s)
}

func autoConvert_controlplane_SubReplicationControllerSpec_To_v1alpha1_SubReplicationControllerSpec(in *controlplane.SubReplicationControllerSpec, out *SubReplicationControllerSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*controlplane.SubReplicationControllerSpec))(in)
	}
	if err := s.Convert(&in.ReplicationControllerSpec, &out.ReplicationControllerSpec, 0); err != nil {
		return err
	}
	if err := Convert_controlplane_ClusterSelectionSpec_To_v1alpha1_ClusterSelectionSpec(&in.Cluster, &out.Cluster, s); err != nil {
		return err
	}
	return nil
}

func Convert_controlplane_SubReplicationControllerSpec_To_v1alpha1_SubReplicationControllerSpec(in *controlplane.SubReplicationControllerSpec, out *SubReplicationControllerSpec, s conversion.Scope) error {
	return autoConvert_controlplane_SubReplicationControllerSpec_To_v1alpha1_SubReplicationControllerSpec(in, out, s)
}

func autoConvert_v1alpha1_Cluster_To_controlplane_Cluster(in *Cluster, out *controlplane.Cluster, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*Cluster))(in)
	}
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ClusterSpec_To_controlplane_ClusterSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ClusterStatus_To_controlplane_ClusterStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_Cluster_To_controlplane_Cluster(in *Cluster, out *controlplane.Cluster, s conversion.Scope) error {
	return autoConvert_v1alpha1_Cluster_To_controlplane_Cluster(in, out, s)
}

func autoConvert_v1alpha1_ClusterAddress_To_controlplane_ClusterAddress(in *ClusterAddress, out *controlplane.ClusterAddress, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*ClusterAddress))(in)
	}
	out.Url = in.Url
	return nil
}

func Convert_v1alpha1_ClusterAddress_To_controlplane_ClusterAddress(in *ClusterAddress, out *controlplane.ClusterAddress, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterAddress_To_controlplane_ClusterAddress(in, out, s)
}

func autoConvert_v1alpha1_ClusterList_To_controlplane_ClusterList(in *ClusterList, out *controlplane.ClusterList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*ClusterList))(in)
	}
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := api.Convert_unversioned_ListMeta_To_unversioned_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]controlplane.Cluster, len(in.Items))
		for i := range in.Items {
			if err := Convert_v1alpha1_Cluster_To_controlplane_Cluster(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func Convert_v1alpha1_ClusterList_To_controlplane_ClusterList(in *ClusterList, out *controlplane.ClusterList, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterList_To_controlplane_ClusterList(in, out, s)
}

func autoConvert_v1alpha1_ClusterSelectionSpec_To_controlplane_ClusterSelectionSpec(in *ClusterSelectionSpec, out *controlplane.ClusterSelectionSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*ClusterSelectionSpec))(in)
	}
	out.Name = in.Name
	if in.Selector != nil {
		out.Selector = make(map[string]string)
		for key, val := range in.Selector {
			out.Selector[key] = val
		}
	} else {
		out.Selector = nil
	}
	return nil
}

func Convert_v1alpha1_ClusterSelectionSpec_To_controlplane_ClusterSelectionSpec(in *ClusterSelectionSpec, out *controlplane.ClusterSelectionSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterSelectionSpec_To_controlplane_ClusterSelectionSpec(in, out, s)
}

func autoConvert_v1alpha1_ClusterSpec_To_controlplane_ClusterSpec(in *ClusterSpec, out *controlplane.ClusterSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*ClusterSpec))(in)
	}
	if err := Convert_v1alpha1_ClusterAddress_To_controlplane_ClusterAddress(&in.Address, &out.Address, s); err != nil {
		return err
	}
	out.Credential = in.Credential
	return nil
}

func Convert_v1alpha1_ClusterSpec_To_controlplane_ClusterSpec(in *ClusterSpec, out *controlplane.ClusterSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterSpec_To_controlplane_ClusterSpec(in, out, s)
}

func autoConvert_v1alpha1_ClusterStatus_To_controlplane_ClusterStatus(in *ClusterStatus, out *controlplane.ClusterStatus, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*ClusterStatus))(in)
	}
	out.Phase = controlplane.ClusterPhase(in.Phase)
	if in.Capacity != nil {
		out.Capacity = make(api.ResourceList)
		for key, val := range in.Capacity {
			newVal := resource.Quantity{}
			if err := s.Convert(&val, &newVal, 0); err != nil {
				return err
			}
			out.Capacity[key] = newVal
		}
	} else {
		out.Capacity = nil
	}
	out.ClusterMeta = in.ClusterMeta
	return nil
}

func Convert_v1alpha1_ClusterStatus_To_controlplane_ClusterStatus(in *ClusterStatus, out *controlplane.ClusterStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterStatus_To_controlplane_ClusterStatus(in, out, s)
}

func autoConvert_v1alpha1_SubReplicationController_To_controlplane_SubReplicationController(in *SubReplicationController, out *controlplane.SubReplicationController, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*SubReplicationController))(in)
	}
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := Convert_v1alpha1_SubReplicationControllerSpec_To_controlplane_SubReplicationControllerSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := s.Convert(&in.Status, &out.Status, 0); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_SubReplicationController_To_controlplane_SubReplicationController(in *SubReplicationController, out *controlplane.SubReplicationController, s conversion.Scope) error {
	return autoConvert_v1alpha1_SubReplicationController_To_controlplane_SubReplicationController(in, out, s)
}

func autoConvert_v1alpha1_SubReplicationControllerList_To_controlplane_SubReplicationControllerList(in *SubReplicationControllerList, out *controlplane.SubReplicationControllerList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*SubReplicationControllerList))(in)
	}
	if err := api.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := api.Convert_unversioned_ListMeta_To_unversioned_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]controlplane.SubReplicationController, len(in.Items))
		for i := range in.Items {
			if err := Convert_v1alpha1_SubReplicationController_To_controlplane_SubReplicationController(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func Convert_v1alpha1_SubReplicationControllerList_To_controlplane_SubReplicationControllerList(in *SubReplicationControllerList, out *controlplane.SubReplicationControllerList, s conversion.Scope) error {
	return autoConvert_v1alpha1_SubReplicationControllerList_To_controlplane_SubReplicationControllerList(in, out, s)
}

func autoConvert_v1alpha1_SubReplicationControllerSpec_To_controlplane_SubReplicationControllerSpec(in *SubReplicationControllerSpec, out *controlplane.SubReplicationControllerSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*SubReplicationControllerSpec))(in)
	}
	if err := s.Convert(&in.ReplicationControllerSpec, &out.ReplicationControllerSpec, 0); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ClusterSelectionSpec_To_controlplane_ClusterSelectionSpec(&in.Cluster, &out.Cluster, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_SubReplicationControllerSpec_To_controlplane_SubReplicationControllerSpec(in *SubReplicationControllerSpec, out *controlplane.SubReplicationControllerSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_SubReplicationControllerSpec_To_controlplane_SubReplicationControllerSpec(in, out, s)
}

func init() {
	err := api.Scheme.AddGeneratedConversionFuncs(
		autoConvert_controlplane_ClusterAddress_To_v1alpha1_ClusterAddress,
		autoConvert_controlplane_ClusterList_To_v1alpha1_ClusterList,
		autoConvert_controlplane_ClusterSelectionSpec_To_v1alpha1_ClusterSelectionSpec,
		autoConvert_controlplane_ClusterSpec_To_v1alpha1_ClusterSpec,
		autoConvert_controlplane_ClusterStatus_To_v1alpha1_ClusterStatus,
		autoConvert_controlplane_Cluster_To_v1alpha1_Cluster,
		autoConvert_controlplane_SubReplicationControllerList_To_v1alpha1_SubReplicationControllerList,
		autoConvert_controlplane_SubReplicationControllerSpec_To_v1alpha1_SubReplicationControllerSpec,
		autoConvert_controlplane_SubReplicationController_To_v1alpha1_SubReplicationController,
		autoConvert_v1alpha1_ClusterAddress_To_controlplane_ClusterAddress,
		autoConvert_v1alpha1_ClusterList_To_controlplane_ClusterList,
		autoConvert_v1alpha1_ClusterSelectionSpec_To_controlplane_ClusterSelectionSpec,
		autoConvert_v1alpha1_ClusterSpec_To_controlplane_ClusterSpec,
		autoConvert_v1alpha1_ClusterStatus_To_controlplane_ClusterStatus,
		autoConvert_v1alpha1_Cluster_To_controlplane_Cluster,
		autoConvert_v1alpha1_SubReplicationControllerList_To_controlplane_SubReplicationControllerList,
		autoConvert_v1alpha1_SubReplicationControllerSpec_To_controlplane_SubReplicationControllerSpec,
		autoConvert_v1alpha1_SubReplicationController_To_controlplane_SubReplicationController,
	)
	if err != nil {
		// If one of the conversion functions is malformed, detect it immediately.
		panic(err)
	}
}
