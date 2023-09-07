package main

import (
	"context"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"balloon-controller/api"
	"balloon-controller/client"
)

func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("get in-cluster config error: %s", err)
	}

	balloonClient, err := client.NewBalloonClient(config, "default")
	if err != nil {
		log.Fatalf("create balloon client error: %s", err)
	}

	balloonLocalCache := NewBalloonLocalCache(balloonClient)
	log.Println("balloon controller started successfully!")

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			balloons := balloonLocalCache.List()

			for _, item := range balloons {
				balloon, _ := item.(*api.Balloon)
				releaseTime, _ := time.Parse(time.RFC3339, balloon.Spec.ReleaseTime)

				if releaseTime.Before(time.Now()) {
					log.Printf("releasing balloon: %s\n", balloon.GetName())

					balloon.Status.Status = "Released"
					result, err := balloonClient.UpdateStatus(
						context.TODO(),
						balloon,
						metav1.UpdateOptions{},
					)
					if err != nil {
						log.Printf("error of update status for balloon %s: %s", balloon.GetName(), err)
					} else {
						log.Printf("balloon: %s is released, status is %s\n", result.GetName(), result.Status.Status)
					}
				}
			}
		}
	}
}
