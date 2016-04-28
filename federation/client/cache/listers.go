/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

package cache

import (
	federation_v1alpha1 "k8s.io/kubernetes/federation/apis/federation/v1alpha1"
	kubeCache "k8s.io/kubernetes/pkg/client/cache"
)

// StoreToClusterLister makes a Store have the List method of the unversioned.ClusterInterface
// The Store must contain (only) Nodes.
type StoreToClusterLister struct {
	kubeCache.Store
}

func (s *StoreToClusterLister) List() (clusters federation_v1alpha1.ClusterList, err error) {
	for _, m := range s.Store.List() {
		clusters.Items = append(clusters.Items, *(m.(*federation_v1alpha1.Cluster)))
	}
	return clusters, nil
}

// StoreToReplicationControllerLister gives a store List and Exists methods. The store must contain only ReplicationControllers.
type StoreToSubReplicaSetLister struct {
	kubeCache.Store
}

// StoreToSubReplicaSetLister lists all replicaSets in the store.
func (s *StoreToSubReplicaSetLister) List() (replicaSets []federation_v1alpha1.SubReplicaSet, err error) {
	for _, c := range s.Store.List() {
		replicaSets = append(replicaSets, *(c.(*federation_v1alpha1.SubReplicaSet)))
	}
	return replicaSets, nil
}
