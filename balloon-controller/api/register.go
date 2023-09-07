package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

var GroupVersion = schema.GroupVersion{
	Group:   "book.dong.io",
	Version: "v1",
}

func init() {
	scheme.Scheme.AddKnownTypes(
		GroupVersion,
		&Balloon{},
		&BalloonList{},
	)

	metav1.AddToGroupVersion(scheme.Scheme, GroupVersion)
}
