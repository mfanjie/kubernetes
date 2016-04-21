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

// If you make changes to this file, you should also make the corresponding change in ReplicationController.

package cluster

import (
	"testing"

	"net/http/httptest"

	"encoding/json"
	"fmt"
	federation_v1alpha1 "k8s.io/kubernetes/federation/apis/federation/v1alpha1"
	federationclientset "k8s.io/kubernetes/federation/client/clientset_generated/release_1_3"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/testapi"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	extensions_v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/util/sets"
	"net/http"
	"strings"
)

func getKey(rs *federation_v1alpha1.SubReplicaSet, t *testing.T) string {
	if key, err := controller.KeyFunc(rs); err != nil {
		t.Errorf("Unexpected error getting key for subreplicaset %v: %v", rs.Name, err)
		return ""
	} else {
		return key
	}
}

func newSubReplicaSet(annotationMap map[string]string, subReplicaSetName string) *federation_v1alpha1.SubReplicaSet {
	var replicas int32 = 1
	rs := &federation_v1alpha1.SubReplicaSet{
		TypeMeta: unversioned.TypeMeta{APIVersion: testapi.Federation.GroupVersion().String()},
		ObjectMeta: v1.ObjectMeta{
			UID:             util.NewUUID(),
			Name:            subReplicaSetName,
			Namespace:       api.NamespaceDefault,
			ResourceVersion: "18",
			Annotations:     annotationMap,
		},
		Spec: extensions_v1beta1.ReplicaSetSpec{
			Replicas: &replicas,
		},
	}
	return rs
}

func newCluster(clusterName string, serverUrl string) *federation_v1alpha1.Cluster {
	cluster := federation_v1alpha1.Cluster{
		TypeMeta: unversioned.TypeMeta{APIVersion: testapi.Federation.GroupVersion().String()},
		ObjectMeta: v1.ObjectMeta{
			UID:  util.NewUUID(),
			Name: clusterName,
		},
		Spec: federation_v1alpha1.ClusterSpec{
			ServerAddressByClientCIDRs: []unversioned.ServerAddressByClientCIDR{
				{
					ClientCIDR:    "0.0.0.0",
					ServerAddress: serverUrl,
				},
			},
		},
	}
	return &cluster
}

func newReplicaSet(rsName string) *extensions_v1beta1.ReplicaSet {
	var replicas int32 = 1
	rs := &extensions_v1beta1.ReplicaSet{
		TypeMeta: unversioned.TypeMeta{APIVersion: testapi.Extensions.GroupVersion().String()},
		ObjectMeta: v1.ObjectMeta{
			UID:             util.NewUUID(),
			Name:            rsName,
			Namespace:       api.NamespaceDefault,
			ResourceVersion: "18",
		},
		Spec: extensions_v1beta1.ReplicaSetSpec{
			Replicas: &replicas,
		},
	}
	return rs
}



var RsNameSetForCreateUseCase = make(sets.String)
var RsNameSetForDeleteUseCase = make(sets.String)
var RsNameSetForUpdateUseCase = make(sets.String)

// init a fake http handler, simulate a cluster apiserver, response the "DELETE" "PUT" "GET" "UPDATE"
// when "canBeGotten" is false, means that user can not get the rs from apiserver
// and subRsNameSet simulate the resource store(rs) of apiserver
func createHttptestFakeHandlerForCluster(rs *extensions_v1beta1.ReplicaSet, canBeGotten bool, subReplicaSetName string) *http.HandlerFunc {
	rs.Name = subReplicaSetName
	rsString, _ := json.Marshal(*rs)
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "DELETE":
			RsNameSetForDeleteUseCase.Delete(rs.Name)
			fmt.Fprintln(w, "")
		case "PUT":
			RsNameSetForUpdateUseCase.Insert(rs.Name)
			fmt.Fprintln(w, string(rsString))
		case "POST":
			RsNameSetForCreateUseCase.Insert(rs.Name)
			fmt.Fprintln(w, string(rsString))
		case "GET":
			if canBeGotten {
				fmt.Fprintln(w, string(rsString))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		default:
			fmt.Fprintln(w, string(rsString))
		}
	})
	return &fakeHandler
}

// init a fake http handler, simulate a federation apiserver, response the "DELETE" "PUT" "GET" "UPDATE"
// when "canBeGotten" is false, means that user can not get the cluster subrs rs from apiserver
func createHttptestFakeHandlerForFederation(
	cluster *federation_v1alpha1.Cluster,
	subRs *federation_v1alpha1.SubReplicaSet,
	canBeGotten bool) *http.HandlerFunc {
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "subreplicaset") {
			subRsString, _ := json.Marshal(*subRs)
			switch r.Method {
			case "DELETE":
				fmt.Fprintln(w, "")
			case "PUT":
				fmt.Fprintln(w, string(subRsString))
			case "GET":
				if canBeGotten {
					fmt.Fprintln(w, string(subRsString))
				} else {
					fmt.Fprintln(w, "")
				}
			default:
				fmt.Fprintln(w, string(subRsString))
			}
		} else {
			clusterString, _ := json.Marshal(*cluster)
			switch r.Method {
			case "DELETE":
				fmt.Fprintln(w, "delete")
			case "PUT":
				fmt.Fprintln(w, "put")
			case "GET":
				if canBeGotten {
					fmt.Fprintln(w, string(clusterString))
				} else {
					fmt.Fprintln(w, "")
				}
			default:
				fmt.Fprintln(w, "")
			}
		}
	})
	return &fakeHandler
}

func TestSyncReplicaSetUpdate(t *testing.T) {
	clusterName := "foobarCluster"
	replicaSetName := "foobarRs"
	subReplicaSetName := "foobarSubRs"
	clusterRS := newReplicaSet(replicaSetName)

	// create dummy httpserver
	testClusterServer := httptest.NewServer(createHttptestFakeHandlerForCluster(clusterRS, true, subReplicaSetName))
	defer testClusterServer.Close()
	federationCluster := newCluster(clusterName, testClusterServer.URL)

	annotationMap := map[string]string{AnnotationKeyOfTargetCluster: clusterName, AnnotationKeyOfFederationReplicaSet: replicaSetName}
	subRs := newSubReplicaSet(annotationMap, subReplicaSetName)


	testFederationServer := httptest.NewServer(createHttptestFakeHandlerForFederation(federationCluster, subRs, true))
        defer testFederationServer.Close()

	kubeClientSet := clientset.NewForConfigOrDie(
		&restclient.Config{
			Host:          testFederationServer.URL,
			ContentConfig: restclient.ContentConfig{GroupVersion: testapi.Default.GroupVersion()},
		},
	)

	federationClientSet := federationclientset.NewForConfigOrDie(
		&restclient.Config{
			Host:          testFederationServer.URL,
			ContentConfig: restclient.ContentConfig{GroupVersion: testapi.Federation.GroupVersion()},
		},
	)
	manager := NewclusterController(kubeClientSet, federationClientSet, 5)

	manager.knownClusterSet.Insert(federationCluster.Name)
	// create the clusterclientset for cluster
	clusterClientSet, err := NewClusterClientSet(federationCluster)
	if err != nil {
		t.Errorf("Failed to create cluster clientset: %v", err)
	}
	manager.clusterKubeClientMap[federationCluster.Name] = *clusterClientSet
	manager.subReplicaSetStore.Add(subRs)

	// sync subreplicaSet
	err = manager.syncSubReplicaSet(getKey(subRs, t))
	if err != nil {
		t.Errorf("Failed to syncSubReplicaSet: %v", err)
	}

	if !RsNameSetForUpdateUseCase.Has(subRs.Name) {
		t.Errorf("Expected rs update in cluster, but not success")
	}
}



func TestSyncReplicaSetCreate(t *testing.T) {
	clusterName := "foobarCluster"
	replicaSetName := "foobarRs"
	subReplicaSetName := "foobarSubRs"
	clusterRS := newReplicaSet(replicaSetName)

	// create dummy httpserver
	testClusterServer := httptest.NewServer(createHttptestFakeHandlerForCluster(clusterRS, false, subReplicaSetName))
	defer testClusterServer.Close()
	federationCluster := newCluster(clusterName, testClusterServer.URL)

	annotationMap := map[string]string{AnnotationKeyOfTargetCluster: clusterName, AnnotationKeyOfFederationReplicaSet: replicaSetName}
	subRs := newSubReplicaSet(annotationMap, subReplicaSetName)

	testFederationServer := httptest.NewServer(createHttptestFakeHandlerForFederation(federationCluster, subRs, true))
        defer testFederationServer.Close()

	kubeClientSet := clientset.NewForConfigOrDie(
		&restclient.Config{
			Host:          testFederationServer.URL,
			ContentConfig: restclient.ContentConfig{GroupVersion: testapi.Default.GroupVersion()},
		},
	)

	federationClientSet := federationclientset.NewForConfigOrDie(
		&restclient.Config{
			Host:          testFederationServer.URL,
			ContentConfig: restclient.ContentConfig{GroupVersion: testapi.Federation.GroupVersion()},
		},
	)
	manager := NewclusterController(kubeClientSet, federationClientSet, 5)

	manager.knownClusterSet.Insert(federationCluster.Name)
	// create the clusterclientset for cluster
	clusterClientSet, err := NewClusterClientSet(federationCluster)
	if err != nil {
		t.Errorf("Failed to create cluster clientset: %v", err)
	}
	manager.clusterKubeClientMap[federationCluster.Name] = *clusterClientSet
	manager.subReplicaSetStore.Add(subRs)

	// sync subreplicaSet
	err = manager.syncSubReplicaSet(getKey(subRs, t))
	if err != nil {
		t.Errorf("Failed to syncSubReplicaSet: %v", err)
	}

	fmt.Println(RsNameSetForCreateUseCase.List())
	if !RsNameSetForCreateUseCase.Has(subRs.Name) {
		t.Errorf("Expected rs create in cluster, but not success")
	}
}

