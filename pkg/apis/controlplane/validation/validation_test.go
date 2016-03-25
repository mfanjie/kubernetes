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
	"strings"
	"testing"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/controlplane"
)

func TestValidateCluster(t *testing.T) {
	successCases := []controlplane.Cluster{
		{
			ObjectMeta: api.ObjectMeta{Name: "cluster-s"},
			Spec: controlplane.ClusterSpec{
				Address: controlplane.ClusterAddress{
					Url: "http://localhost:8888",
				},
			},
		},
	}
	for _, successCase := range successCases {
		errs := ValidateCluster(&successCase)
		if len(errs) != 0 {
			t.Errorf("expect success: %v", errs)
		}
	}

	errorCases := map[string]controlplane.Cluster{
		"missing cluster address": {
			ObjectMeta: api.ObjectMeta{Name: "cluster-f"},
		},
		"invalid_label": {
			ObjectMeta: api.ObjectMeta{
				Name: "cluster-f",
				Labels: map[string]string{
					"NoUppercaseOrSpecialCharsLike=Equals": "bar",
				},
			},
		},
	}
	for testName, errorCase := range errorCases {
		errs := ValidateCluster(&errorCase)
		if len(errs) == 0 {
			t.Errorf("expected failur for %s", testName)
		}
	}
}

func TestValidateClusterUpdate(t *testing.T) {
	type clusterUpdateTest struct {
		old    controlplane.Cluster
		update controlplane.Cluster
	}
	successCases := []clusterUpdateTest{
		{
			old: controlplane.Cluster{
				ObjectMeta: api.ObjectMeta{Name: "cluster-s"},
				Spec: controlplane.ClusterSpec{
					Address: controlplane.ClusterAddress{
						Url: "http://localhost:8888",
					},
				},
			},
			update: controlplane.Cluster{
				ObjectMeta: api.ObjectMeta{Name: "cluster-s"},
				Spec: controlplane.ClusterSpec{
					Address: controlplane.ClusterAddress{
						Url: "http://127.0.0.1:8888",
					},
				},
			},
		},
	}
	for _, successCase := range successCases {
		successCase.old.ObjectMeta.ResourceVersion = "1"
		successCase.update.ObjectMeta.ResourceVersion = "1"
		errs := ValidateClusterUpdate(&successCase.update, &successCase.old)
		if len(errs) != 0 {
			t.Errorf("expect success: %v", errs)
		}
	}

	errorCases := map[string]clusterUpdateTest{}
	for testName, errorCase := range errorCases {
		errs := ValidateClusterUpdate(&errorCase.update, &errorCase.old)
		if len(errs) == 0 {
			t.Errorf("expected failure: %s", testName)
		}
	}
}

func TestValidateReplicationController(t *testing.T) {
	validSelector := map[string]string{"a": "b"}
	validPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: api.ObjectMeta{
				Labels: validSelector,
			},
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
			},
		},
	}
	readWriteVolumePodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: api.ObjectMeta{
				Labels: validSelector,
			},
			Spec: api.PodSpec{
				Volumes:       []api.Volume{{Name: "gcepd", VolumeSource: api.VolumeSource{GCEPersistentDisk: &api.GCEPersistentDiskVolumeSource{PDName: "my-PD", FSType: "ext4", Partition: 1, ReadOnly: false}}}},
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
			},
		},
	}
	successCases := []controlplane.SubReplicationController{
		{
			ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
			Spec: controlplane.SubReplicationControllerSpec{
				ReplicationControllerSpec: api.ReplicationControllerSpec{
					Selector: validSelector,
					Template: &validPodTemplate.Template,
				},
				Cluster: controlplane.ClusterSelectionSpec{
					Name:     "def",
					Selector: map[string]string{"a": "b"},
				},
			},
		},
		{
			ObjectMeta: api.ObjectMeta{Name: "abc-123", Namespace: api.NamespaceDefault},
			Spec: controlplane.SubReplicationControllerSpec{
				ReplicationControllerSpec: api.ReplicationControllerSpec{
					Replicas: 1,
					Selector: validSelector,
					Template: &readWriteVolumePodTemplate.Template,
				},
				Cluster: controlplane.ClusterSelectionSpec{
					Name:     "def",
					Selector: map[string]string{"a": "b"},
				},
			},
		},
	}
	for _, successCase := range successCases {
		if errs := ValidateSubReplicationController(&successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := map[string]controlplane.SubReplicationController{
		"invalid-cluster-name": {
			ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
			Spec: controlplane.SubReplicationControllerSpec{
				ReplicationControllerSpec: api.ReplicationControllerSpec{
					Selector: validSelector,
					Template: &validPodTemplate.Template,
				},
				Cluster: controlplane.ClusterSelectionSpec{
					Name: "asdf+b",
				},
			},
		},
		"invalid-cluster-selector": {
			ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
			Spec: controlplane.SubReplicationControllerSpec{
				ReplicationControllerSpec: api.ReplicationControllerSpec{
					Selector: validSelector,
					Template: &validPodTemplate.Template,
				},
				Cluster: controlplane.ClusterSelectionSpec{
					Name:     "asdf",
					Selector: map[string]string{"a+.": "b"},
				},
			},
		},
	}
	for k, v := range errorCases {
		errs := ValidateSubReplicationController(&v)
		if len(errs) == 0 {
			t.Errorf("expected failure for %s", k)
		}
		for i := range errs {
			field := errs[i].Field
			if !strings.HasPrefix(field, "spec.template.") &&
				field != "metadata.name" &&
				field != "metadata.namespace" &&
				field != "spec.selector" &&
				field != "spec.template" &&
				field != "GCEPersistentDisk.ReadOnly" &&
				field != "spec.replicas" &&
				field != "spec.template.labels" &&
				field != "metadata.annotations" &&
				field != "metadata.labels" &&
				field != "status.replicas" &&
				field != "spec.cluster.name" &&
				field != "spec.cluster.selector" {
				t.Errorf("%s: missing prefix for: %v", k, errs[i])
			}
		}
	}
}
