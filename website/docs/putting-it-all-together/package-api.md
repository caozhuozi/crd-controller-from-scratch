---
sidebar_position: 1
---

# package api

## types.go
```go
package api

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type BalloonSpec struct {
	ReleaseTime string `json:"release_time"`
}

type Balloon struct {
	metav1.TypeMeta   `json:",inline"`  // 🤖️ (1)
	metav1.ObjectMeta `json:"metadata,omitempty"`  // 🤖️ (2)
    
    // 🤖️ (3)
	Spec BalloonSpec `json:"spec"`
	Status BalloonStatus `json:"status"`
}

type BalloonStatus struct {
	Status string `json:"status"`
}

type BalloonList struct {
	metav1.TypeMeta `json:",inline"`  // 🤖️ (1)
	metav1.ListMeta `json:"metadata,omitempty"` // 🤖️ (4)

	Items []Balloon `json:"items"` // 🤖️ (4)
}
```

1. [kind与runtime.Object](../apimachinery/runtime.Object#kind与runtimeobject)
2. [kind的单体种类](../apimachinery/runtime.Object#单体类型)
3. [Kubernetes对象的期望状态与实际状态](../client-go/controller#Kubernetes对象的期望状态与实际状态)
4. [kind的集合种类](../apimachinery/runtime.Object#集合类型)


## deepcopy.go
```go
package api

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// 🤖️ (1)
func (in *Balloon) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Balloon)
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)

	return out
}

// 🤖️ (1)
func (in *BalloonList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}

	out := new(BalloonList)

	*out = *in
	in.ListMeta.DeepCopyInto(&out.ListMeta)

	if in.Items != nil {
		in, out := &in.Items, &out.Items
		for i := range *in {
			c := (*in)[i].DeepCopyObject().(*Balloon)
			*out = append(*out, *c)
		}
	}

	return out
}

```
1. [kind与runtime.Object](../apimachinery/runtime.Object#kind与runtimeobject)

## register.go
```go
package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

var GroupVersion = schema.GroupVersion{Group: "book.dong.io", Version: "v1"}

func init() {
    // 🤖️ (1) (2)
	scheme.Scheme.AddKnownTypes(GroupVersion,
		&Balloon{},
		&BalloonList{},
	)
	// 🤖️ (2) (3)
	metav1.AddToGroupVersion(scheme.Scheme, GroupVersion)
}
```
1. [初识序列化器](../apimachinery/runtime.Object#初识序列化器)
2. [kind的注册原理](../apimachinery/client-go-and-apimachinery#client-go中的全局scheme对象)
3. [kind中的特殊种类](../apimachinery/gvk#再识kind)