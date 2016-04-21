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

package cluster

import (
	extensions_v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"

	federation_v1alpha1 "k8s.io/kubernetes/federation/apis/federation/v1alpha1"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_3"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/typed/discovery"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"strings"
)

const (
	UserAgentName = "Cluster-Controller"
	KubeAPIQPS    = 20.0
	KubeAPIBurst  = 30
)

type ClusterClient struct {
	clientSet       clientset.Interface
	discoveryClient *discovery.DiscoveryClient
}

func NewClusterClientSet(c *federation_v1alpha1.Cluster) (*ClusterClient, error) {
	//TODO:How to get cluster IP(huangyuqi)
	var clusterClientSet = ClusterClient{}
	clusterConfig, err := clientcmd.BuildConfigFromFlags(c.Spec.ServerAddressByClientCIDRs[0].ServerAddress, "")
	if err != nil {
		return nil, err
	}
	clusterConfig.QPS = KubeAPIQPS
	clusterConfig.Burst = KubeAPIBurst
	clusterClientSet.clientSet = clientset.NewForConfigOrDie(restclient.AddUserAgent(clusterConfig, UserAgentName))
	clusterClientSet.discoveryClient = discovery.NewDiscoveryClientForConfigOrDie((restclient.AddUserAgent(clusterConfig, UserAgentName)))
	if clusterClientSet.clientSet == nil || clusterClientSet.discoveryClient == nil {
		return nil, nil
	}
	return &clusterClientSet, err
}

// GetReplicaSetFromCluster get the replicaset from the kubernetes cluster
func (self *ClusterClient) GetReplicaSetFromCluster(rsName string, rsNameSpace string) (*extensions_v1beta1.ReplicaSet, error) {
	return self.clientSet.Extensions().ReplicaSets(rsNameSpace).Get(rsName)
}

// CreateReplicaSetToCluster create replicaset to the kubernetes cluster
func (self *ClusterClient) CreateReplicaSetToCluster(rs *extensions_v1beta1.ReplicaSet) (*extensions_v1beta1.ReplicaSet, error) {
	rs.ResourceVersion = ""
	return self.clientSet.Extensions().ReplicaSets(rs.Namespace).Create(rs)
}

// UpdateReplicaSetToCluster update replicaset to the kubernetes cluster
func (self *ClusterClient) UpdateReplicaSetToCluster(rs *extensions_v1beta1.ReplicaSet) (*extensions_v1beta1.ReplicaSet, error) {
	return self.clientSet.Extensions().ReplicaSets(rs.Namespace).Update(rs)
}

// DeleteReplicasetFromCluster delete the replicaset from the kubernetes cluster
func (self *ClusterClient) DeleteReplicasetFromCluster(rs *extensions_v1beta1.ReplicaSet) error {
	return self.clientSet.Extensions().ReplicaSets(rs.Namespace).Delete(rs.Name, &api.DeleteOptions{})
}

// GetClusterHealthStatus get the kubernetes cluster health status by requesting "/healthz"
func (self *ClusterClient) GetClusterHealthStatus() *federation_v1alpha1.ClusterStatus {
	clusterStatus := federation_v1alpha1.ClusterStatus{}
	currentTime := unversioned.Now()
	newClusterReadyCondition := federation_v1alpha1.ClusterCondition{
		Type:               federation_v1alpha1.ClusterReady,
		Status:             v1.ConditionTrue,
		Reason:             "federate cluster is ready",
		Message:            "federate cluster response the '/healthz' with 'OK' ",
		LastProbeTime:      currentTime,
		LastTransitionTime: currentTime,
	}
	newClusterNotReadyCondition := federation_v1alpha1.ClusterCondition{
		Type:               federation_v1alpha1.ClusterReady,
		Status:             v1.ConditionFalse,
		Reason:             "federate cluster is not ready",
		Message:            "federate cluster response the '/healthz' request without 'OK' ",
		LastProbeTime:      currentTime,
		LastTransitionTime: currentTime,
	}
	newNodeOfflineCondition := federation_v1alpha1.ClusterCondition{
		Type:               federation_v1alpha1.ClusterOffline,
		Status:             v1.ConditionTrue,
		Reason:             "federate cluster is not reachable",
		Message:            "federate cluster can not response the '/healthz' request",
		LastProbeTime:      currentTime,
		LastTransitionTime: currentTime,
	}
	newNodeNotOfflineCondition := federation_v1alpha1.ClusterCondition{
		Type:               federation_v1alpha1.ClusterOffline,
		Status:             v1.ConditionFalse,
		Reason:             "federate cluster is reachable",
		Message:            "federate cluster can response the '/healthz' request",
		LastProbeTime:      currentTime,
		LastTransitionTime: currentTime,
	}
	body, err := self.discoveryClient.Get().AbsPath("/healthz").Do().Raw()
	if err != nil {
		clusterStatus.Conditions = append(clusterStatus.Conditions, newNodeOfflineCondition)
	} else {
		if !strings.EqualFold(string(body), "ok") {
			clusterStatus.Conditions = append(clusterStatus.Conditions, newClusterNotReadyCondition, newNodeNotOfflineCondition)
		} else {
			clusterStatus.Conditions = append(clusterStatus.Conditions, newClusterReadyCondition)
		}
	}
	return &clusterStatus
}
