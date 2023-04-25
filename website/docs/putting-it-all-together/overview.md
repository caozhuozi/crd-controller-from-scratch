---
sidebar_position: 0
id: putting-it-all-together
---
当你读到这里时，你应该已经拥有了实现一个简单的自定义控制器的所有预备知识。

本章的内容是一个"气球🎈控制器🤖️"的实现，完整代码也可以通过github仓库[caozhuozi/balloon-controller](https://github.com/caozhuozi/balloon-controller)获取。
另外，正如在前言中[本书结构](TODO[cross-reference]:)中所说的那样，本章不会对代码再做过多的解释和说明，本章更像是一个代码到知识点的索引。
我们尽量会为代码中的每个主要部分分配一个跳转到前两章对应知识点的链接。

## 气球控制器

### *气球*自定义资源（*Balloon* CRD）
```yaml
apiVersion: "apiextensions.k8s.io/v1"
kind: "CustomResourceDefinition"
metadata:
  name: "balloons.book.dong.io"
spec:
  group: "book.dong.io"
  versions:
    - name: v1
      # Each version can be enabled/disabled by Served flag.
      served: true
      # One and only one version must be marked as the storage version.
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                releaseTime:
                  type: "string"
                  format: "date-time"
            status:
              type: object
              properties:
                status:
                  type: "string"
      subresources:
        # status enables the status subresource.
        status: { }

  scope: "Namespaced"
  names:
    plural: "balloons"
    singular: "balloon"
    kind: "Balloon"
```
就像在本书[前言](../intro#谁适合阅读本书)中所说的一样， 安装*CRD*来扩展Kubernetes API应该是读者需要预先掌握的知识，我们在这里不多做赘述。
为了尽量让我们的例子简单一些，*气球*自定义资源的`spec`字段下仅有一个字段`releaseTime`描述气球期望的释放时间，这也就是*气球*资源的*期望状态*。
同时，`status`字段下也仅有一个字段`status`用于记录气球资源实际是否已经被释放，这也代表*气球*资源的*实际状态*。 
在此CRD中，唯一需要注意的地方是以下字段的声明用于开启`status`子资源：[🤖️](../client-go/controller#kubernetes对象子资源status)
```yaml
subresources: 
  status: { } 
```
### 气球控制器
  
本章实现的气球控制器大概分为四个部分，代码目录结构如下所示：
```text
|-- api
|   |-- deepcopy.go
|   |-- register.go
|   `-- types.go
|-- client
|   `-- client.go
|-- informer.go
`-- main.go
```
* 包`api`里包括了：
  * *气球*资源对应的*kind*：`Ballon`类型（单体类型）以及`BalloonList`的定义（集合类型）；
  * 将与`Ballon`相关的类型（包括集合类型，**特殊及通用类型**）注册进`client-go`库中的全局`Schema`。
* 包`client`里是基于`client-go`的`RESTClient`封装的一个简单的气球资源类型客户端；
* `informer.go`文件中实现了一个基于`client-go`的`Infomer`机制的气球资源本地缓存；
* `main.go`文件是控制器的主体逻辑：每分钟访问一次气球资源集合的本地缓存，当气球期望的释放时间和当前时间匹配时，
  则"释放"该气球，并更新气球资源的实际状态。

  
