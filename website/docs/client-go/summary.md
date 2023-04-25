---
sidebar_position: 4
---

# 本章小结

本章我们从`client-go`中构成`Clientset`的基础组件`RESTClient`的用法说起，介绍了如何使用它与Kubernetes API直接交互进而管理资源。
之后我们正式介绍了什么是Kubernetes控制器——它本质上是一种具有特殊行为的资源客户端。
接着我们又以优化Kubernetes控制器的实现为线索，引入了`client-go`中的`Informer`组件。

至此，我们已经掌握了实现一个自定义资源控制器的所有组件。

