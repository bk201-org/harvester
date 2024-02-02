/*
Copyright 2024 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/generic"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ControllerRevisionHandler func(string, *v1.ControllerRevision) (*v1.ControllerRevision, error)

type ControllerRevisionController interface {
	generic.ControllerMeta
	ControllerRevisionClient

	OnChange(ctx context.Context, name string, sync ControllerRevisionHandler)
	OnRemove(ctx context.Context, name string, sync ControllerRevisionHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ControllerRevisionCache
}

type ControllerRevisionClient interface {
	Create(*v1.ControllerRevision) (*v1.ControllerRevision, error)
	Update(*v1.ControllerRevision) (*v1.ControllerRevision, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.ControllerRevision, error)
	List(namespace string, opts metav1.ListOptions) (*v1.ControllerRevisionList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ControllerRevision, err error)
}

type ControllerRevisionCache interface {
	Get(namespace, name string) (*v1.ControllerRevision, error)
	List(namespace string, selector labels.Selector) ([]*v1.ControllerRevision, error)

	AddIndexer(indexName string, indexer ControllerRevisionIndexer)
	GetByIndex(indexName, key string) ([]*v1.ControllerRevision, error)
}

type ControllerRevisionIndexer func(obj *v1.ControllerRevision) ([]string, error)

type controllerRevisionController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewControllerRevisionController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ControllerRevisionController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &controllerRevisionController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromControllerRevisionHandlerToHandler(sync ControllerRevisionHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.ControllerRevision
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.ControllerRevision))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *controllerRevisionController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.ControllerRevision))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateControllerRevisionDeepCopyOnChange(client ControllerRevisionClient, obj *v1.ControllerRevision, handler func(obj *v1.ControllerRevision) (*v1.ControllerRevision, error)) (*v1.ControllerRevision, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *controllerRevisionController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *controllerRevisionController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *controllerRevisionController) OnChange(ctx context.Context, name string, sync ControllerRevisionHandler) {
	c.AddGenericHandler(ctx, name, FromControllerRevisionHandlerToHandler(sync))
}

func (c *controllerRevisionController) OnRemove(ctx context.Context, name string, sync ControllerRevisionHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromControllerRevisionHandlerToHandler(sync)))
}

func (c *controllerRevisionController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *controllerRevisionController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *controllerRevisionController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *controllerRevisionController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *controllerRevisionController) Cache() ControllerRevisionCache {
	return &controllerRevisionCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *controllerRevisionController) Create(obj *v1.ControllerRevision) (*v1.ControllerRevision, error) {
	result := &v1.ControllerRevision{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *controllerRevisionController) Update(obj *v1.ControllerRevision) (*v1.ControllerRevision, error) {
	result := &v1.ControllerRevision{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *controllerRevisionController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *controllerRevisionController) Get(namespace, name string, options metav1.GetOptions) (*v1.ControllerRevision, error) {
	result := &v1.ControllerRevision{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *controllerRevisionController) List(namespace string, opts metav1.ListOptions) (*v1.ControllerRevisionList, error) {
	result := &v1.ControllerRevisionList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *controllerRevisionController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *controllerRevisionController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.ControllerRevision, error) {
	result := &v1.ControllerRevision{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type controllerRevisionCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *controllerRevisionCache) Get(namespace, name string) (*v1.ControllerRevision, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.ControllerRevision), nil
}

func (c *controllerRevisionCache) List(namespace string, selector labels.Selector) (ret []*v1.ControllerRevision, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ControllerRevision))
	})

	return ret, err
}

func (c *controllerRevisionCache) AddIndexer(indexName string, indexer ControllerRevisionIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.ControllerRevision))
		},
	}))
}

func (c *controllerRevisionCache) GetByIndex(indexName, key string) (result []*v1.ControllerRevision, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.ControllerRevision, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.ControllerRevision))
	}
	return result, nil
}