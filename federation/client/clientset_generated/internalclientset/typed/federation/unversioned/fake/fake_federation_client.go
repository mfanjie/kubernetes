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

// This file is generated by client-gen with arguments: --clientset-path=k8s.io/kubernetes/federation/client/clientset_generated --input=[../../federation/apis/federation/]

package fake

import (
	unversioned "k8s.io/kubernetes/federation/client/clientset_generated/internalclientset/typed/federation/unversioned"
	core "k8s.io/kubernetes/pkg/client/testing/core"
)

type FakeFederation struct {
	*core.Fake
}

func (c *FakeFederation) Clusters() unversioned.ClusterInterface {
	return &FakeClusters{c}
}
