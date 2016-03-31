// +build !ignore_autogenerated

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package federation

import (
	api "k8s.io/kubernetes/pkg/api"
	resource "k8s.io/kubernetes/pkg/api/resource"
	unversioned "k8s.io/kubernetes/pkg/api/unversioned"
	conversion "k8s.io/kubernetes/pkg/conversion"
)

func init() {
	if err := api.Scheme.AddGeneratedDeepCopyFuncs(
		DeepCopy_federation_Cluster,
		DeepCopy_federation_ClusterAddress,
		DeepCopy_federation_ClusterList,
		DeepCopy_federation_ClusterMeta,
		DeepCopy_federation_ClusterSpec,
		DeepCopy_federation_ClusterStatus,
	); err != nil {
		// if one of the deep copy functions is malformed, detect it immediately.
		panic(err)
	}
}

func DeepCopy_federation_Cluster(in Cluster, out *Cluster, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if err := DeepCopy_federation_ClusterSpec(in.Spec, &out.Spec, c); err != nil {
		return err
	}
	if err := DeepCopy_federation_ClusterStatus(in.Status, &out.Status, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_federation_ClusterAddress(in ClusterAddress, out *ClusterAddress, c *conversion.Cloner) error {
	out.Url = in.Url
	return nil
}

func DeepCopy_federation_ClusterList(in ClusterList, out *ClusterList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]Cluster, len(in))
		for i := range in {
			if err := DeepCopy_federation_Cluster(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_federation_ClusterMeta(in ClusterMeta, out *ClusterMeta, c *conversion.Cloner) error {
	out.Version = in.Version
	return nil
}

func DeepCopy_federation_ClusterSpec(in ClusterSpec, out *ClusterSpec, c *conversion.Cloner) error {
	if in.ServerAddressByClientCIDRs != nil {
		in, out := in.ServerAddressByClientCIDRs, &out.ServerAddressByClientCIDRs
		*out = make([]unversioned.ServerAddressByClientCIDR, len(in))
		for i := range in {
			if err := unversioned.DeepCopy_unversioned_ServerAddressByClientCIDR(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.ServerAddressByClientCIDRs = nil
	}
	out.Credential = in.Credential
	return nil
}

func DeepCopy_federation_ClusterStatus(in ClusterStatus, out *ClusterStatus, c *conversion.Cloner) error {
	out.Phase = in.Phase
	if in.Capacity != nil {
		in, out := in.Capacity, &out.Capacity
		*out = make(api.ResourceList)
		for key, val := range in {
			newVal := new(resource.Quantity)
			if err := resource.DeepCopy_resource_Quantity(val, newVal, c); err != nil {
				return err
			}
			(*out)[key] = *newVal
		}
	} else {
		out.Capacity = nil
	}
	if in.Allocatable != nil {
		in, out := in.Allocatable, &out.Allocatable
		*out = make(api.ResourceList)
		for key, val := range in {
			newVal := new(resource.Quantity)
			if err := resource.DeepCopy_resource_Quantity(val, newVal, c); err != nil {
				return err
			}
			(*out)[key] = *newVal
		}
	} else {
		out.Allocatable = nil
	}
	if err := DeepCopy_federation_ClusterMeta(in.ClusterMeta, &out.ClusterMeta, c); err != nil {
		return err
	}
	return nil
}
