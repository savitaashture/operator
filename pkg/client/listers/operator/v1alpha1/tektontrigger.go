/*
Copyright 2020 The Tekton Authors

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// TektonTriggerLister helps list TektonTriggers.
type TektonTriggerLister interface {
	// List lists all TektonTriggers in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.TektonTrigger, err error)
	// Get retrieves the TektonTrigger from the index for a given name.
	Get(name string) (*v1alpha1.TektonTrigger, error)
	TektonTriggerListerExpansion
}

// tektonTriggerLister implements the TektonTriggerLister interface.
type tektonTriggerLister struct {
	indexer cache.Indexer
}

// NewTektonTriggerLister returns a new TektonTriggerLister.
func NewTektonTriggerLister(indexer cache.Indexer) TektonTriggerLister {
	return &tektonTriggerLister{indexer: indexer}
}

// List lists all TektonTriggers in the indexer.
func (s *tektonTriggerLister) List(selector labels.Selector) (ret []*v1alpha1.TektonTrigger, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TektonTrigger))
	})
	return ret, err
}

// Get retrieves the TektonTrigger from the index for a given name.
func (s *tektonTriggerLister) Get(name string) (*v1alpha1.TektonTrigger, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("tektontrigger"), name)
	}
	return obj.(*v1alpha1.TektonTrigger), nil
}
