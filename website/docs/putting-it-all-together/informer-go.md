---
sidebar_position: 3
---

# informer.go
```go
package main

import (
	"context"
	"fmt"
	"github.com/caozhuozi/balloon-controller/api"
	"github.com/caozhuozi/balloon-controller/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"time"
)

func NewBalloonLocalCache(balloonClient *client.BalloonClient) cache.Store {
    // 🤖️ (1)
	balloonStore, balloonController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(opts metav1.ListOptions) (result runtime.Object, err error) {
				return balloonClient.List(context.TODO(), opts)
			},
			WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
				return balloonClient.Watch(context.TODO(), opts)
			},
		},
		&api.Balloon{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				fmt.Printf("The ballon: %s is created.", obj.(*api.Balloon).Name)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				fmt.Printf("The ballon: %s is updated.", oldObj.(*api.Balloon).Name)
			},
			DeleteFunc: func(obj interface{}) {
				fmt.Printf("The ballon: %s is deleted.", obj.(*api.Balloon).Name)
			},
		},
	)

	go balloonController.Run(wait.NeverStop)
	return balloonStore
}
```

1. [Informer](../client-go/informer-is-all-you-need)

