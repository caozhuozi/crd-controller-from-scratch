---
sidebar_position: 4
---

# main.go
```go
package main

import (
	"context"
	"fmt"
	"github.com/caozhuozi/balloon-controller/api"
	"github.com/caozhuozi/balloon-controller/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"time"
)

func main() {

	config, _ := rest.InClusterConfig()
	balloonClient, _ := client.NewBalloonClient(config, "default")

	balloonLocalCache := NewBalloonLocalCache(balloonClient)

	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			balloons := balloonLocalCache.List()

			for _, item := range balloons {
				balloon, _ := item.(*api.Balloon)

				releaseTime, _ := time.Parse(time.RFC3339, balloon.Spec.ReleaseTime)
				if releaseTime.After(time.Now()) {

					fmt.Printf("releasing balloon: %s\n", balloon.GetName())

					balloon.Status.Status = "Released"
					_, _ = balloonClient.UpdateStatus(context.TODO(), balloon, metav1.UpdateOptions{})
					fmt.Printf("balloon: %s is released\n", balloon.GetName())
				}
			}
		}
	}

}
```