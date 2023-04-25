---
sidebar_position: 5
---


# client-go与apimachinery

在本节中，我们将介绍`client-go`是如何使用`apimachinery`中的序列化工具的。

## client-go中的全局Scheme对象
在之前的小节中，我们介绍了`apimachinery`库提供的序列化器需要检查scheme中是否注册了相应的kind。

我们知道`client-go`中的资源客户端支持所有原生Kubernetes资源类型，我们有理由猜测`client-go`中应该存在一个scheme注册了所有kind。
事实的确如此，`client-go`存在一个全局的Scheme对象`Secheme`，它被声明在`client-go/kubernetes/kubernetes/scheme`包中：
```go title="client-go/kubernetes/scheme/register.go"
var Scheme = runtime.NewScheme()
```
它注册了Kubernetes中所有原生的kind（包括[前面小节](TODO[cross-reference]:)总结的**所有**三个种类）[^1][^2]。

## client-go中的全局序列化器工厂

在`client-go`中，还存在一个全局的序列化器工厂`Codecs`，它使用的正是`client-go`中的全局scheme对象`Scheme`。

```go title="client-go/kubernetes/scheme/register.go"
var Codecs = serializer.NewCodecFactory(Scheme)
```
这个全局序列化器工厂"生产的"的序列化器将会被应用到`client-go`封装的各个资源客户端中用于与`kube-apiserver`通信的编/解码工作。

## client-go中的全局URL参数"序列化器"
上面的全局`Codec`是用于请求/返回**体**的编解码。`client-go`中用于将Go对象转化为Kubernetes API URL中查询参数（query parameter）的全局"序列化器"为`ParameterCodec`:
```go title="client-go/kubernetes/scheme/register.go"
ParameterCodec = runtime.NewParameterCodec(Scheme)
```

同样它使用的也是`client-go`中的全局scheme对象`Scheme`。

[^1]: `client-go`中的scheme包并没有**直接**调用`runtime.Scheme`的`AddKnownTypeWithName()`方法进行注册。而是使用了`runtime.SchemeBuilder`类**间接**地注册了所有原生的kind。
      这样做的原因是Kubernetes所有原生的kind并不在`client-go`中，而是被定义在了`k8s.io/api`库中。这种跨库的局限性使得我们只能找出一种间接的方式把`client-go`外的kind注册进`client-go`中的全局`Scheme`中。
      `apimachinery`中的`SchemeBuilder`类旨在优雅方便地解决跨库的注册`Scheme`问题，具体操作如下：
      
      首先，在`k8s.io/api`库中，每个API分组下都提供了一个相应注册函数，例如`core/v1`这个分组：
      ```go title="k8s.io/api/core/v1/register.go"
      func addKnownTypes(scheme *runtime.Scheme) error {
	      scheme.AddKnownTypes(SchemeGroupVersion, &Pod{}, &Service{}, // ...)
          // ...
      }

      var (
          SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
          AddToScheme   = SchemeBuilder.AddToScheme
      )
      ```
      `addKnownTypes`的作用是把此版本分组下的所有kind(`Pod`，`Serive`)注册进传入的任意一个`scheme`对象参数中。
      同时利用`SchemeBuilder`类将注册函数封装成一个函数变量`AddToScheme`供对外导出。
      这相当于使`Pod`这个kind对外提供了注册进任意一个`scheme`对象的能力。
      
      在`client-go`中，通过导入`k8s.io/api`库中所有原生kind的`AddToScheme`注册函数，并将全局的`Scheme`对象传入这些注册函数中并执行，
      这就使得`client-go`的全局`Scheme`获得了所有原生的kind。
      ```go
      import (
      	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
      	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
      	internalv1alpha1 "k8s.io/api/apiserverinternal/v1alpha1"
      	appsv1 "k8s.io/api/apps/v1"
      	appsv1beta1 "k8s.io/api/apps/v1beta1"
      	appsv1beta2 "k8s.io/api/apps/v1beta2"
      	authenticationv1 "k8s.io/api/authentication/v1"
      	authenticationv1beta1 "k8s.io/api/authentication/v1beta1"
      	authorizationv1 "k8s.io/api/authorization/v1"
      	authorizationv1beta1 "k8s.io/api/authorization/v1beta1"
      	autoscalingv1 "k8s.io/api/autoscaling/v1"
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
      	networkingv1 "k8s.io/api/networking/v1"
      	networkingv1beta1 "k8s.io/api/networking/v1beta1"
      	nodev1 "k8s.io/api/node/v1"
      	nodev1alpha1 "k8s.io/api/node/v1alpha1"
      	nodev1beta1 "k8s.io/api/node/v1beta1"
      	policyv1 "k8s.io/api/policy/v1"
      	policyv1beta1 "k8s.io/api/policy/v1beta1"
      	rbacv1 "k8s.io/api/rbac/v1"
      	rbacv1alpha1 "k8s.io/api/rbac/v1alpha1"
      	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
      	schedulingv1 "k8s.io/api/scheduling/v1"
      	schedulingv1alpha1 "k8s.io/api/scheduling/v1alpha1"
      	schedulingv1beta1 "k8s.io/api/scheduling/v1beta1"
      	storagev1 "k8s.io/api/storage/v1"
      	storagev1alpha1 "k8s.io/api/storage/v1alpha1"
      	storagev1beta1 "k8s.io/api/storage/v1beta1"
      	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
      	runtime "k8s.io/apimachinery/pkg/runtime"
      	schema "k8s.io/apimachinery/pkg/runtime/schema"
      	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
      	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
      )
      
      var Scheme = runtime.NewScheme()
      var Codecs = serializer.NewCodecFactory(Scheme)
      var ParameterCodec = runtime.NewParameterCodec(Scheme)
      var localSchemeBuilder = runtime.SchemeBuilder{
      	admissionregistrationv1.AddToScheme,
      	admissionregistrationv1beta1.AddToScheme,
      	internalv1alpha1.AddToScheme,
      	appsv1.AddToScheme,
      	appsv1beta1.AddToScheme,
      	appsv1beta2.AddToScheme,
      	authenticationv1.AddToScheme,
      	authenticationv1beta1.AddToScheme,
      	authorizationv1.AddToScheme,
      	authorizationv1beta1.AddToScheme,
      	autoscalingv1.AddToScheme,
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
      	networkingv1.AddToScheme,
      	networkingv1beta1.AddToScheme,
      	nodev1.AddToScheme,
      	nodev1alpha1.AddToScheme,
      	nodev1beta1.AddToScheme,
      	policyv1.AddToScheme,
      	policyv1beta1.AddToScheme,
      	rbacv1.AddToScheme,
      	rbacv1beta1.AddToScheme,
      	rbacv1alpha1.AddToScheme,
      	schedulingv1alpha1.AddToScheme,
      	schedulingv1beta1.AddToScheme,
      	schedulingv1.AddToScheme,
      	storagev1beta1.AddToScheme,
      	storagev1.AddToScheme,
      	storagev1alpha1.AddToScheme,
      }
      
      var AddToScheme = localSchemeBuilder.AddToScheme
      
      func init() {
      	v1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
      	utilruntime.Must(AddToScheme(Scheme))
      }
      ```

[^2]: 对于kind的第三种特殊种类，`apimachinery`库的`metav1`包统一为这些特殊kind提供了一个注册函数`AddToGroupVersion`:
      ```go title="k8s.io/apimachinery/pkg/apis/meta/v1/register.go"
      func AddToGroupVersion(scheme *runtime.Scheme, groupVersion schema.GroupVersion) {
        
      	scheme.AddKnownTypeWithName(groupVersion.WithKind(WatchEventKind), &WatchEvent{})
        // ...
      	scheme.AddKnownTypes(groupVersion, 
            &ListOptions{},
	        &GetOptions{},
	        &DeleteOptions{},
	        &CreateOptions{},
	        &UpdateOptions{},
	        &PatchOptions{},
        )
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
      这个函数被包括在`k8s.io/api`库中每个API分组对应的注册函数中，例如`core/v1`的`addKnownTypes`：
      ```go title="k8s.io/api/core/v1/register.go"
      func addKnownTypes(scheme *runtime.Scheme) error {
          // ...
          metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
      }
      ```
      所以当每个分组的注册函数被执行的时候这些特殊kind也被注册进了相应的`scheme`对象中。
      
      另外，由于`client-go`会执行所有API分组的注册函数，所以这个特殊kind的注册函数也会被执行多次，因此我们可以推断出对于这些特殊kind存在重复注册的情况。
    




