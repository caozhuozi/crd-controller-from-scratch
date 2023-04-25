---
sidebar_position: 0
id: apimachinery
---
了解`apimachinery` 库对于编写自定义控制器至关重要。

我们知道，在一个Kubernetes集群中，`kube-apiserver`作为控制平面（control plane）的核心组件对外提供了集群中各类*资源*[^1]增、删、改、查（以及watch）的HTTP RESTful接口。
本质上来说，`kube-apiserver`仍然是一个服务端，它在接受、处理和返回请求的过程中也必然涉及到**序列化**和**反序列化**[^2]。
同样地，客户端（`client-go`）在发送请求和接受服务端返回的过程中也会涉及到序列化和反序列化[^3]。

[comment]: # (TODO: 是否应该用官方的描述？：https://kubernetes.io/docs/concepts/overview/kubernetes-api/)
[comment]: # (TODO[figure]:  这里需要补一章kubectl和api-server的请求交互图：参考：https://devopscube.com/kubernetes-objects-resources/)
而有关对*Kubernetes API对象*[^1]*编/解码*[^4]的实现则在`apimachinery`这个基础库中，
因此，`k8s.io/client-go`和`k8s.io/apiserver`都依赖于`k8s.io/apimachinery`。
正如`apimachinery`库的官方介绍一样：
> This library is a shared dependency for servers and clients to work with Kubernetes API infrastructure without direct type dependencies.
> Its first consumers are k8s.io/kubernetes, k8s.io/client-go, and k8s.io/apiserver.

回到自定义控制器，虽然你现在可能并不知道Kubernetes控制器具体是什么，但是资源的控制器显然是需要对资源（无论是否为自定义资源）进行管理（增删改查等），`kube-apiserver`作为集群**唯一**[^5]对外的接口，
我们不可避免地需要与它直接交互。 这个交互的过程自然涉及到对HTTP请求体/返回体的编/解码。这也是我们需要深入了解`apimachinery`基础库的原因。


:::tip 小结
我们从序列化的角度引入了`apimachinery`基础库，实际上序列化只是`apimachinery`库的一部分内容。
其实`apimachinery`库可以被认为是Kubernetes API的基石——它系统性地定义了*Kubernetes API对象*[^1]在Go语言中的代表形式。
:::


[^1]: 此处的*资源*指的是*Kubernetes资源*（它也可以被成为*Kubernetes API对象*），并不是指用于声明容器中CPU/memory的`resources`，而是Kubernetes API上下文中的一个专业术语。在之后的[Kubernetes API 基础](./kubernetes-api#)小节中，我们会介绍它具体指的是什么。
[^2]: `kube-apiserver`除了要将网络传输的数据（HTTP 请求、返回）序列化/反序列化以外，将*Kubernetes资源*写入磁盘中（etcd）同样也涉及到序列化操作。
[^3]: 一个Kubernetes API请求的生命周期远远不止序列化操作，但这些内容不在本书的讨论范围之内，请参考扩展阅读[Kubernetes API请求的生命周期](../intro#扩展阅读)。
[^4]: 在本书中，*编/解码*与*序列化/反序列化*意思等同并且被交换使用。
[^5]: Kubernetes集群架构设计不在本书的讨论范围，请参考扩展阅读[Kubernetes资源管理设计文档](intro#扩展阅读).


