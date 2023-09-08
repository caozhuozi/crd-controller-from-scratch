package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BalloonSpec struct {
	ReleaseTime string `json:"releaseTime"`
}

type BalloonStatus struct {
	Status string `json:"status"`
}

type Balloon struct {
	// metav1.TypeMeta 实现了 runtime.Object 的第一个方法 GetObjectKind()
	// 当我们引入后 只需要再实现 DeepCopyObject() 就能实现 runtime.Object 接口了
	// 我们可以认为 TypeMeta 代表 CRD 最基本的 group version kind 信息
	metav1.TypeMeta `json:",inline"`

	// Balloon 作为单体需要实现 metav1.Object 接口的
	// 而 metav1.ObjectMeta 则实现了这个接口
	// 对于 ObjectMeta 可以理解为 CRD 的 name namespace annotation 等信息
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec 就是具体的 yaml 中配置信息了 对应 spec 字段
	Spec BalloonSpec `json:"spec"`

	// status 字段涉及到子资源这一概念
	// 用户可以写入/更改资源的期望状态 spec 但不应更改资源的 status 字段
	// 控制器可以写入/更改资源的实际状态 status 但不应更改资源的 spec 字段
	Status BalloonStatus `json:"status"`
}

type BalloonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Balloon `json:"items"`
}
