---
sidebar_position: 4
---
import PodPNG from '@site/static/img/pod.png';

# kind案例探究

在本节我们将通过`k8s.io/api`库的一个具体的*kind*案例——`Pod`及`PodList`来进一步诠释*kind*在实现上是如何具体需要遵循的约定的。
探究原生*kind*的具体实现有助于我们实现*自定义资源*的Go类型。


## Pod

在Go语言中，如果一个类型实现了某接口，将此类型嵌入到另一个结构类型中，那么被嵌入的结构类型也等同于实现了该接口。

在`k8s.io/api`库中定义了所有Kubernetes原生资源类型对应的kind。
这些kind正是通过**类型嵌入**的方式来达到之前[kind与runtime.Object](./runtime.Object#kind与runtimeobject)小节所提及的实现基础接口`runtime.Object`同等的效果。

例如， 我们以`Pod`类型为例，
```go title="k8s.io/api/core/v1/types.go"
type Pod struct {
	metav1.TypeMeta 
	metav1.ObjectMeta 

	Spec PodSpec 
	Status PodStatus 
}
```
1. 通过嵌入`metav1.TypeMeta`类型，使得`Pod`也实现了`runtime.Object`接口的第一个方法`GetObjectKind() schema.ObjectKind`。
2. 通过嵌入`metav1.ObjectMeta`类型， 使得`Pod`也实现了`metav1.Object`接口。

3. 除此以外，`Pod`类型**未以类型嵌入的方式**实现了`runtime.Object`接口的第二个方法`DeepCopyObject() Object`：
   ```go title="k8s.io/api/core/v1/zz_generated.deepcopy.go"
   func (in *Pod) DeepCopyObject() runtime.Object {
       if c := in.DeepCopy(); c != nil {
           return c
       }
       return nil
   }
   ```

综合以上三点，`Pod`类型同时实现了`runtime.Object`接口以及`metav1.Object`接口。根据定义，`Pod`类型应属于kind第一种类（即资源类型对应的单体种类）。

<img src={PodPNG} width="95%" />


## PodList

我们再以`PodList`为例：
```go title="k8s.io/api/core/v1/types.go"
type PodList struct {
	metav1.TypeMeta
	metav1.ListMeta 

	Items []Pod
}
```
1. 通过嵌入`metav1.TypeMeta`结构类型，使得`PodList`类型实现了`runtime.Object`接口的第一个方法`GetObjectKind() schema.ObjectKind`。
2. 通过嵌入`metav1.ListMeta`结构类型， 使得`PodList`类型实现了`metav1.ListInterface`接口。
3. 除此以外，`PodList`类型**未以类型嵌入的方式**实现了`runtime.Object`接口的第二个方法`DeepCopyObject() Object`：
   ```go title="k8s.io/api/core/v1/zz_generated.deepcopy.go"
   func (in *PodList) DeepCopyObject() runtime.Object {
	   if c := in.DeepCopy(); c != nil {
		   return c
	   }
	   return nil
   }
   ```
4. `PodList`类型含有名为`Items`的slice字段。

综合以上四点，`PodList`类型同时实现了`runtime.Object`接口和`metav1.ListInterface`接口，并且含有名为`Items`的slice字段。根据定义，`PodList`类型应属于kind第二种类（即集合种类）。


## 小结
:::tip 小结
本节通过对`Pod`及`PodList`这两个*kind*的探究，旨在为读者建立起一个编写*自定义资源*Go类型对应的模板。
:::
