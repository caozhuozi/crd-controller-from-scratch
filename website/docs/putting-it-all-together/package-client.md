---
sidebar_position: 2
---

# package client

## client.go

```client
package client

import (
	"context"
	"github.com/caozhuozi/balloon-controller/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type BalloonClient struct {
	restClient rest.Interface
	ns         string
}

func setConfigDefaults(config *rest.Config) error {
    // 🤖️ (1)
	gv := api.GroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	// 🤖️ (2)
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	return nil

}

func NewBalloonClient(c *rest.Config, namespace string) (*BalloonClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	// 🤖️ (1)
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &BalloonClient{restClient: client, ns: namespace}, nil
}

// 🤖️ (1)
func (c *BalloonClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*api.Balloon, error) {
	result := api.Balloon{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("balloons").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}
// 🤖️ (1)
func (c *BalloonClient) List(ctx context.Context, opts metav1.ListOptions) (*api.BalloonList, error) {
	result := api.BalloonList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("balloons").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

// 🤖️ (1)
func (c *BalloonClient) Create(ctx context.Context, balloon *api.Balloon) (*api.Balloon, error) {
	result := &api.Balloon{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("balloons").
		Body(balloon).
		Do(ctx).
		Into(result)

	return result, err
}

// 🤖️ (1) (2)
func (c *BalloonClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("balloons").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

// 🤖️ (1) (3)
func (c *BalloonClient) UpdateStatus(ctx context.Context, balloon *api.Balloon, opts metav1.UpdateOptions) (result *api.Balloon, err error) {
	result = &api.Balloon{}
	err = c.restClient.Put().
		Namespace(c.ns).
		Resource("balloons").
		Name(balloon.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(balloon).
		Do(ctx).
		Into(result)
	return result, err
}
```
1. [RESTClient基本用法](../client-go/restclient#restclient基本用法)
2. [watch机制](../client-go/controller#被遗忘的监听watch机制)
3. [status子资源](../client-go/controller#kubernetes对象子资源status)
