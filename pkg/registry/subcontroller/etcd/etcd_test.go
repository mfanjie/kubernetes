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

package etcd

import (
	"testing"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/controlplane"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/generic"
	"k8s.io/kubernetes/pkg/registry/registrytest"
	"k8s.io/kubernetes/pkg/runtime"
	etcdtesting "k8s.io/kubernetes/pkg/storage/etcd/testing"
)

const (
	namespace = api.NamespaceDefault
	name      = "foo"
)

func newStorage(t *testing.T) (*REST, *etcdtesting.EtcdTestServer) {
	etcdStorage, server := registrytest.NewEtcdStorage(t, controlplane.GroupName)
	restOptions := generic.RESTOptions{etcdStorage, generic.UndecoratedStorage, 1}
	storage, _ := NewREST(restOptions)
	return storage, server
}

// createController is a helper function that returns a controller with the updated resource version.
func createController(storage *REST, rc controlplane.SubReplicationController, t *testing.T) (controlplane.SubReplicationController, error) {
	ctx := api.WithNamespace(api.NewContext(), rc.Namespace)
	obj, err := storage.Create(ctx, &rc)
	if err != nil {
		t.Errorf("Failed to create controller, %v", err)
	}
	newRc := obj.(*controlplane.SubReplicationController)
	return *newRc, nil
}

func validNewController() *controlplane.SubReplicationController {
	return &controlplane.SubReplicationController{
		ObjectMeta: api.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"a": "b"},
		},
		Spec: controlplane.SubReplicationControllerSpec{
			ReplicationControllerSpec: api.ReplicationControllerSpec{
				Selector: map[string]string{"a": "b"},
				Template: &api.PodTemplateSpec{
					ObjectMeta: api.ObjectMeta{
						Labels: map[string]string{"a": "b"},
					},
					Spec: api.PodSpec{
						Containers: []api.Container{
							{
								Name:            "test",
								Image:           "test_image",
								ImagePullPolicy: api.PullIfNotPresent,
							},
						},
						RestartPolicy: api.RestartPolicyAlways,
						DNSPolicy:     api.DNSClusterFirst,
					},
				},
			},
			Cluster: controlplane.ClusterSelectionSpec{
				Name: "abc",
			},
		},
	}
}

var validController = validNewController()

func TestCreate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Etcd)
	controller := validNewController()
	controller.ObjectMeta = api.ObjectMeta{}
	test.TestCreate(
		// valid
		controller,
		// invalid (invalid selector)
		&controlplane.SubReplicationController{
			Spec: controlplane.SubReplicationControllerSpec{
				ReplicationControllerSpec: api.ReplicationControllerSpec{
					Replicas: 2,
					Selector: map[string]string{},
					Template: validController.Spec.Template,
				},
				Cluster: controlplane.ClusterSelectionSpec{
					Name: "abcd",
				},
			},
		},
	)
}

func TestUpdate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Etcd)
	test.TestUpdate(
		// valid
		validNewController(),
		// valid updateFunc
		func(obj runtime.Object) runtime.Object {
			object := obj.(*controlplane.SubReplicationController)
			object.Spec.Replicas = object.Spec.Replicas + 1
			return object
		},
		// invalid updateFunc
		func(obj runtime.Object) runtime.Object {
			object := obj.(*controlplane.SubReplicationController)
			object.UID = "newUID"
			return object
		},
		func(obj runtime.Object) runtime.Object {
			object := obj.(*controlplane.SubReplicationController)
			object.Name = ""
			return object
		},
		func(obj runtime.Object) runtime.Object {
			object := obj.(*controlplane.SubReplicationController)
			object.Spec.Selector = map[string]string{}
			return object
		},
	)
}

func TestDelete(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Etcd)
	test.TestDelete(validNewController())
}

func TestGenerationNumber(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	modifiedSno := *validNewController()
	modifiedSno.Generation = 100
	modifiedSno.Status.ObservedGeneration = 10
	ctx := api.NewDefaultContext()
	rc, err := createController(storage, modifiedSno, t)
	ctrl, err := storage.Get(ctx, rc.Name)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	controller, _ := ctrl.(*controlplane.SubReplicationController)

	// Generation initialization
	if controller.Generation != 1 && controller.Status.ObservedGeneration != 0 {
		t.Fatalf("Unexpected generation number %v, status generation %v", controller.Generation, controller.Status.ObservedGeneration)
	}

	// Updates to spec should increment the generation number
	controller.Spec.Replicas += 1
	storage.Update(ctx, controller)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	ctrl, err = storage.Get(ctx, rc.Name)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	controller, _ = ctrl.(*controlplane.SubReplicationController)
	if controller.Generation != 2 || controller.Status.ObservedGeneration != 0 {
		t.Fatalf("Unexpected generation, spec: %v, status: %v", controller.Generation, controller.Status.ObservedGeneration)
	}

	// Updates to status should not increment either spec or status generation numbers
	controller.Status.Replicas += 1
	storage.Update(ctx, controller)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	ctrl, err = storage.Get(ctx, rc.Name)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	controller, _ = ctrl.(*controlplane.SubReplicationController)
	if controller.Generation != 2 || controller.Status.ObservedGeneration != 0 {
		t.Fatalf("Unexpected generation number, spec: %v, status: %v", controller.Generation, controller.Status.ObservedGeneration)
	}
}

func TestGet(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Etcd)
	test.TestGet(validNewController())
}

func TestList(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Etcd)
	test.TestList(validNewController())
}

func TestWatch(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	test := registrytest.New(t, storage.Etcd)
	test.TestWatch(
		validController,
		// matching labels
		[]labels.Set{
			{"a": "b"},
		},
		// not matching labels
		[]labels.Set{
			{"a": "c"},
			{"foo": "bar"},
		},
		// matching fields
		[]fields.Set{
			{"status.replicas": "0"},
			{"metadata.name": "foo"},
			{"status.replicas": "0", "metadata.name": "foo"},
		},
		// not matchin fields
		[]fields.Set{
			{"status.replicas": "10"},
			{"metadata.name": "bar"},
			{"name": "foo"},
			{"status.replicas": "10", "metadata.name": "foo"},
			{"status.replicas": "0", "metadata.name": "bar"},
		},
	)
}
