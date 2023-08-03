---
sidebar_position: 6
---


# client-go与apimachinery

我们知道Kubernetes客户端向Kubernetes API Server发送HTTP请求同样涉及到对资源的序列化/反序列化。
只是不同于服务端，客户端对资源编解码时不需要再进行版本的转化。
在本节中，我们将介绍`client-go`是如何具体使用`apimachinery`中提供的序列化工具的。

## client-go中的全局Scheme对象
在之前的小节中，我们介绍了`apimachinery`库的编解码器需要检查`scheme`中是否注册了相应的*kind*。 而`client-go`中的资源客户端支持所有原生Kubernetes资源类型，我们有理由猜测`client-go`中应该存在一个`scheme`注册了所有Kubernetes原生*kind*。

事实也的确如此，`client-go`存在一个全局的`Scheme`类型的对象`Secheme`，它被定义在`kubernetes/kubernetes/scheme`包中，
它注册了Kubernetes中所有原生*kind*（包括[GVK](./gvk.mdx)总结的**所有**三个种类）[🎈](../intro#约定)：
```go title="client-go/kubernetes/scheme/register.go"
var Scheme = runtime.NewScheme()
```

不过一个值得注意的事情是我们在[初识kind](./kubernetes-api#初识kind)已经介绍了Kubernetes所有原生的*kind*被定义在`k8s.io/api`库中。
那么`client-go`中这个全局的`Scheme`对象是怎么注册上原生*kind*的呢？

难道是从`k8s.io/api`库中导入所有版本的*kind*吗？显然这并不是一个好的办法。
其实官方开发者已经提供了相关的基础代码和组件，旨在为我们提供一种便捷的方式来完成**跨库**的*kind*注册。

### addKnownTypes

首先，在`k8s.io/api`库中，以API分组为单位，开发者为每个分组预先定义了一个注册函数，这个函数的签名如下所示：

```go
func addKnownTypes(scheme *runtime.Scheme) error
```

此函数用于将此分组下的所有*kind*注册进给定的`scheme`中。我们以`core/v1`这个API分组为例：

<details>
<summary>addKnownTypes</summary>

```go title="k8s.io/api/core/v1/register.go"
// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Pod{},
		&PodList{},
		&PodStatusResult{},
		&PodTemplate{},
		&PodTemplateList{},
		&ReplicationController{},
		&ReplicationControllerList{},
		&Service{},
		&ServiceProxyOptions{},
		&ServiceList{},
		&Endpoints{},
		&EndpointsList{},
		&Node{},
		&NodeList{},
		&NodeProxyOptions{},
		&Binding{},
		&Event{},
		&EventList{},
		&List{},
		&LimitRange{},
		&LimitRangeList{},
		&ResourceQuota{},
		&ResourceQuotaList{},
		&Namespace{},
		&NamespaceList{},
		&Secret{},
		&SecretList{},
		&ServiceAccount{},
		&ServiceAccountList{},
		&PersistentVolume{},
		&PersistentVolumeList{},
		&PersistentVolumeClaim{},
		&PersistentVolumeClaimList{},
		&PodAttachOptions{},
		&PodLogOptions{},
		&PodExecOptions{},
		&PodPortForwardOptions{},
		&PodProxyOptions{},
		&ComponentStatus{},
		&ComponentStatusList{},
		&SerializedReference{},
		&RangeAllocation{},
		&ConfigMap{},
		&ConfigMapList{},
	)
```
</details>
此注册函数将`core/v1`下所有的*kind*（例如`Pod`、`PodList`）注册进给定的`scheme`中。

也就是说在`client-go`中，我们只要导入`k8s.io/api`库中所有分组下的这个注册函数`addKnownTypes`，并将全局的`Scheme`对象传入其中执行即可——相比于调用`Scheme`类型的`AddKnownTypes`方法一个一个地注册`kind`，这种**按组批量**注册的方式确实方便许多。

### SchemeBuilder

虽然官方为每个API分组预定义的`addKnownTypes`函数减轻了我们注册的工作量，但是这种方式仍然需要我们一遍遍地去执行所有导入的注册函数，类似于：
```go
import (
    // highlight-start
    admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
    admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
    internalv1alpha1 "k8s.io/api/apiserverinternal/v1alpha1"
    appsv1 "k8s.io/api/apps/v1"
    appsv1beta1 "k8s.io/api/apps/v1beta1"
    appsv1beta2 "k8s.io/api/apps/v1beta2"
    authenticationv1 "k8s.io/api/authentication/v1"
    authenticationv1beta1 "k8s.io/api/authentication/v1beta1"
    // ...
    // highlight-end
)

var Scheme = runtime.NewScheme()

// highlight-start
admissionregistrationv1.AddKnownTypes(Scheme)
admissionregistrationv1beta1.AddKnownTypes(Scheme)
internalv1alpha1.AddKnownTypes(Scheme)
appsv1.AddKnownTypes(Scheme)
appsv1beta1.AddKnownTypes(Scheme)
appsv1beta2.AddKnownTypes(Scheme)
authenticationv1.AddKnownTypes(Scheme)
authenticationv1beta1.AddKnownTypes(Scheme)
//...
// highlight-end

```

而这在官方开发者看来仍然不够优雅。为了解决这个问题，官方开发者在`apimachinery`库中特地提供了`runtime.SchemeBuilder`类。我们先来看看这个类具体的使用方法，我们以`k8s.io/api/core/v1/register.go`中的用法为例：
```go title="k8s.io/api/core/v1/register.go"
var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)
```
其中：
* `NewSchemeBuilder()`"吸收"一个注册函数`addKnownTypes`并创建出一个`SchemeBuilder`对象；
* `SchemeBuilder`的`AddToScheme`成员用于返回刚刚"吸收"的`addKnownTypes`，也就是返回的`AddToScheme`就是`addKnownTypes`函数。

看起来似乎`SchemeBuilder`类型只是将注册函数"左手倒右手"，那它存在的意义又是什么呢？
其实`NewSchemeBuilder()`函数支持同时"吸收"**多个**注册函数：
```go
func NewSchemeBuilder(funcs ...func(*Scheme) error) SchemeBuilder {
    // ...
}
```
`SchemeBuilder`的`AddToScheme`成员其实将"吸收"的**多个**注册函数在逻辑上封装成**一个**。
这样， 仅通过调用一次`AddToSchme(scheme)`就可以一次性地执行多个注册函数。

当然，如果在创建`SchemeBuilder`对象时只传入一个注册函数，就会造成"左手倒右手"的现象。

### 向全局Scheme注册原生kind

我们现在知道了官方开发者已经在`k8s.io/api`库中为我们事先准备了各个API分组的注册函数，并且在`k8s.io/apimachinery`库中也为我们提供了`SchemeBuilder`类型用于"优雅"地执行注册函数，我们现在来看看`client-go`中的全局`Scheme`对象是如何注册上所有原生*kind*的：

<details>
<summary>向全局Scheme对象注册所有原生kind</summary>

```go title="client-go/kubernetes/scheme/register.go"

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	admissionregistrationv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	internalv1alpha1 "k8s.io/api/apiserverinternal/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	authenticationv1 "k8s.io/api/authentication/v1"
	authenticationv1alpha1 "k8s.io/api/authentication/v1alpha1"
	authenticationv1beta1 "k8s.io/api/authentication/v1beta1"
	authorizationv1 "k8s.io/api/authorization/v1"
	authorizationv1beta1 "k8s.io/api/authorization/v1beta1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	certificatesv1 "k8s.io/api/certificates/v1"
	certificatesv1beta1 "k8s.io/api/certificates/v1beta1"
	coordinationv1 "k8s.io/api/coordination/v1"
	coordinationv1beta1 "k8s.io/api/coordination/v1beta1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	discoveryv1beta1 "k8s.io/api/discovery/v1beta1"
	eventsv1 "k8s.io/api/events/v1"
	eventsv1beta1 "k8s.io/api/events/v1beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	flowcontrolv1alpha1 "k8s.io/api/flowcontrol/v1alpha1"
	flowcontrolv1beta1 "k8s.io/api/flowcontrol/v1beta1"
	flowcontrolv1beta2 "k8s.io/api/flowcontrol/v1beta2"
	flowcontrolv1beta3 "k8s.io/api/flowcontrol/v1beta3"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1alpha1 "k8s.io/api/networking/v1alpha1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	nodev1 "k8s.io/api/node/v1"
	nodev1alpha1 "k8s.io/api/node/v1alpha1"
	nodev1beta1 "k8s.io/api/node/v1beta1"
	policyv1 "k8s.io/api/policy/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1alpha1 "k8s.io/api/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	resourcev1alpha1 "k8s.io/api/resource/v1alpha1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	schedulingv1alpha1 "k8s.io/api/scheduling/v1alpha1"
	schedulingv1beta1 "k8s.io/api/scheduling/v1beta1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1alpha1 "k8s.io/api/storage/v1alpha1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var Scheme = runtime.NewScheme()

var localSchemeBuilder = runtime.SchemeBuilder{
	admissionregistrationv1.AddToScheme,
	admissionregistrationv1alpha1.AddToScheme,
	admissionregistrationv1beta1.AddToScheme,
	internalv1alpha1.AddToScheme,
	appsv1.AddToScheme,
	appsv1beta1.AddToScheme,
	appsv1beta2.AddToScheme,
	authenticationv1.AddToScheme,
	authenticationv1alpha1.AddToScheme,
	authenticationv1beta1.AddToScheme,
	authorizationv1.AddToScheme,
	authorizationv1beta1.AddToScheme,
	autoscalingv1.AddToScheme,
	autoscalingv2.AddToScheme,
	autoscalingv2beta1.AddToScheme,
	autoscalingv2beta2.AddToScheme,
	batchv1.AddToScheme,
	batchv1beta1.AddToScheme,
	certificatesv1.AddToScheme,
	certificatesv1beta1.AddToScheme,
	coordinationv1beta1.AddToScheme,
	coordinationv1.AddToScheme,
	corev1.AddToScheme,
	discoveryv1.AddToScheme,
	discoveryv1beta1.AddToScheme,
	eventsv1.AddToScheme,
	eventsv1beta1.AddToScheme,
	extensionsv1beta1.AddToScheme,
	flowcontrolv1alpha1.AddToScheme,
	flowcontrolv1beta1.AddToScheme,
	flowcontrolv1beta2.AddToScheme,
	flowcontrolv1beta3.AddToScheme,
	networkingv1.AddToScheme,
	networkingv1alpha1.AddToScheme,
	networkingv1beta1.AddToScheme,
	nodev1.AddToScheme,
	nodev1alpha1.AddToScheme,
	nodev1beta1.AddToScheme,
	policyv1.AddToScheme,
	policyv1beta1.AddToScheme,
	rbacv1.AddToScheme,
	rbacv1beta1.AddToScheme,
	rbacv1alpha1.AddToScheme,
	resourcev1alpha1.AddToScheme,
	schedulingv1alpha1.AddToScheme,
	schedulingv1beta1.AddToScheme,
	schedulingv1.AddToScheme,
	storagev1beta1.AddToScheme,
	storagev1.AddToScheme,
	storagev1alpha1.AddToScheme,
}

var AddToScheme = localSchemeBuilder.AddToScheme

func init() {
	
	utilruntime.Must(AddToScheme(Scheme))
}

```

</details>

在已经掌握了我们铺垫的预备知识的情况下，`client-go`中这段向全局`Scheme`对象中注册原生`kind`的逻辑就显得十分清晰了：
1. 从`k8s.io/api`导入所有的API分组
2. 调用`runtime.SchemeBuilder`"吸收"所有分组的预注册函数并生成一个`SchemeBuilder`对象：`localSchemeBuilder`
3. 调用`localSchemeBuilder`的`AddToScheme`成员以获取一个逻辑上包括所有预注册函数的函数：`AddToScheme`
4. 将全局`Scheme`对象传入`AddToScheme()`函数：即执行所有预注册函数完成所有原生*kind*的注册



### 向全局Scheme注册特殊kind

到目前为止，我们仅介绍了`client-go`的全局`Scheme`注册**`k8s.io/api`**库中定义的原生*kind*的过程。
而`k8s.io/api`库中的*kind*仅包括*单体种类*以及*集合种类*。对于*kind*的第三种类（通用及特殊类型），它们被定义在`apimachinery`中。
接下来我们将介绍`client-go`的全局`Scheme`注册*kind*第三种类的过程。

#### AddToGroupVersion

像`k8s.io/api`库中提供的预注册函数`addKnownTypes()`一样，`apimachinery`库中也提供了这些特殊*kind*的预注册函数，不过相比于`k8s.io/api`库中每个API分组中都存在一个注册函数，由于特殊*kind*本身就很少，
`apimachinery`库中仅有一个注册函数叫做`AddToGroupVersion()`用于注册所有通用及特殊的*kind*，它被定义在了`metav1`包中：
<details>
<summary>apimachinery库中特殊kind的预注册函数</summary>

```go title="apimachinery/pkg/apis/meta/v1/register.go"
// AddToGroupVersion registers common meta types into schemas.
func AddToGroupVersion(scheme *runtime.Scheme, groupVersion schema.GroupVersion) {
	scheme.AddKnownTypeWithName(groupVersion.WithKind(WatchEventKind), &WatchEvent{})
	scheme.AddKnownTypeWithName(
		schema.GroupVersion{Group: groupVersion.Group, Version: runtime.APIVersionInternal}.WithKind(WatchEventKind),
		&InternalEvent{},
	)
	// Supports legacy code paths, most callers should use metav1.ParameterCodec for now
	scheme.AddKnownTypes(groupVersion, optionsTypes...)
	// Register Unversioned types under their own special group
	scheme.AddUnversionedTypes(Unversioned,
		&Status{},
		&APIVersions{},
		&APIGroupList{},
		&APIGroup{},
		&APIResourceList{},
	)
	
	// ...
}
```

</details>


在`client-go/kubernetes/scheme/register.go`文件中，全局`Scheme`不仅注册了所有`k8s.io/api`中所有的*kind*，也注册了`k8s.io/apimachinery`中通用及特殊的*kind*：
```go title="client-go/kubernetes/scheme/register.go"
import (
    // ...
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

var Scheme = runtime.NewScheme()

// ...

func init() {
	v1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	// ...
}
```







## client-go中的全局序列化器工厂

在`client-go`中，也存在一个全局的序列化器工厂对象`Codecs`，`client-go`使用的正是全局`Scheme`而创建的它：

```go title="client-go/kubernetes/scheme/register.go"
var Codecs = serializer.NewCodecFactory(Scheme)
```
这个全局序列化器工厂负责`client-go`中所有与`kube-apiserver`通信的编/解码工作。[🎈](../intro#约定)

## client-go中的全局URL参数"序列化器"
上述`Codec`是用于请求/返回**体**的编解码。`client-go`中用于将Go对象转化为Kubernetes API URL参数（Query Parameter）的全局URL参数"序列化器"为`ParameterCodec`[🎈](../intro#约定):
```go title="client-go/kubernetes/scheme/register.go"
ParameterCodec = runtime.NewParameterCodec(Scheme)
```

同样它使用的也是`client-go`中的全局`Scheme`对象。


## 小结
:::tip 小结
其实我们使用`client-go`的资源客户端`clientset`时并不需要了解任何与序列化/反序列化有关的细节，这些细节被封装在`clientset`内。
我们之所以在本节中探究`client-go`中使用的序列化器是因为我们在编写*自定义资源*资源控制器时，封装完备的原生资源客户端`clientset`对我们来说已经没有用处了。
因此我们需要了解探究`client-go`背后与`kube-apiserver`通信的细节，而序列化/反序列化就是其中一部分。
:::
    




