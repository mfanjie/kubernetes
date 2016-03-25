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

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/rest"
	"k8s.io/kubernetes/pkg/apis/controlplane"
	"k8s.io/kubernetes/pkg/watch"
)

type Registry interface {
	ListSubControllers(ctx api.Context, options *api.ListOptions) (*controlplane.SubReplicationControllerList, error)
	WatchSubControllers(ctx api.Context, options *api.ListOptions) (watch.Interface, error)
	GetSubController(ctx api.Context, controllerID string) (*controlplane.SubReplicationController, error)
	CreateSubController(ctx api.Context, controller *controlplane.SubReplicationController) (*controlplane.SubReplicationController, error)
	UpdateSubController(ctx api.Context, controller *controlplane.SubReplicationController) (*controlplane.SubReplicationController, error)
	DeleteSubController(ctx api.Context, controllerID string) error
}

type storage struct {
	rest.StandardStorage
}

func NewRegistry(s rest.StandardStorage) Registry {
	return &storage{s}
}

func (s *storage) ListSubControllers(ctx api.Context, options *api.ListOptions) (*controlplane.SubReplicationControllerList, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*controlplane.SubReplicationControllerList), err
}

func (s *storage) WatchSubControllers(ctx api.Context, options *api.ListOptions) (watch.Interface, error) {
	return s.Watch(ctx, options)
}

func (s *storage) GetSubController(ctx api.Context, controllerID string) (*controlplane.SubReplicationController, error) {
	obj, err := s.Get(ctx, controllerID)
	if err != nil {
		return nil, err
	}
	return obj.(*controlplane.SubReplicationController), nil
}

func (s *storage) CreateSubController(ctx api.Context, controller *controlplane.SubReplicationController) (*controlplane.SubReplicationController, error) {
	obj, err := s.Create(ctx, controller)
	if err != nil {
		return nil, err
	}
	return obj.(*controlplane.SubReplicationController), nil
}

func (s *storage) UpdateSubController(ctx api.Context, controller *controlplane.SubReplicationController) (*controlplane.SubReplicationController, error) {
	obj, _, err := s.Update(ctx, controller)
	if err != nil {
		return nil, err
	}
	return obj.(*controlplane.SubReplicationController), nil
}

func (s *storage) DeleteSubController(ctx api.Context, controllerID string) error {
	_, err := s.Delete(ctx, controllerID, nil)
	return err
}
