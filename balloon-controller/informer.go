package main

import (
	"context"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"balloon-controller/api"
	"balloon-controller/client"
)

func NewBalloonLocalCache(balloonClient *client.BalloonClient) cache.Store {
	balloonStore, balloonController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
				return balloonClient.List(context.TODO(), opts)
			},
			WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
				return balloonClient.Watch(context.TODO(), opts)
			},
		},
		&api.Balloon{},
		30*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				log.Printf("The balloon: %s is created.\n", obj.(*api.Balloon).Name)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				log.Printf("The balloon: %s is updated.\n", oldObj.(*api.Balloon).Name)
			},
			DeleteFunc: func(obj interface{}) {
				log.Printf("The balloon: %s is deleted.\n", obj.(*api.Balloon).Name)
			},
		},
	)

	go balloonController.Run(wait.NeverStop)
	return balloonStore
}
