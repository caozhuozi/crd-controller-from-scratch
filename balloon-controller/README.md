> 该部分内容提供了 **balloon-controller** 的具体实现，我们将自定义控制器打包成容器镜像并上传到镜像仓库。通过在 kubernetes 集群中部署相关的 manifest 让自定义控制器真正的运行起来。
> 
> 本部分的代码与[《从零实现Kubernetes自定义控制器》](https://caozhuozi.github.io/crd-controller-from-scratch/)系列内容以及 [balloon-controller](https://github.com/caozhuozi/balloon-controller) 仓库中的内容基本相同。
> 添加了些 log 相关的内容便于调试。
> 
> 我们假定你已经对容器镜像构建、kubernetes 权限管理等内容有基本的了解。

# Prerequisite

将 **balloon-controller** 部署到 kubernetes 集群，或是上手进行编码工作，你需要准备下面的环境：

- 可用的 kubernetes 集群，用以部署我们的自定义控制器
- Golang 开发环境，用以编写控制器的代码
- Docker，用以打包和上传自定义控制器的镜像

# Install

自定义控制器的镜像可以使用我们已经打包好的，也可以选择自己[打包镜像](#Build Image)。

### Step1: 部署权限配置文件

**balloon-controller** 在 `default` 命名空间使用默认的服务账户，我们需要为其配置权限使其能够对我们自定义的 `Balloon` 资源进行一些基本的 `watch/get/list/update` 的操作：

```shell
# cd crd-controller-from-scratch/balloon-controller/manifest/
# kubectl apply -f role-binding-sa.yaml 
role.rbac.authorization.k8s.io/balloon-controller-role created
rolebinding.rbac.authorization.k8s.io/balloon-controller-rolebinding created
```

### Step2: 部署 CRD 和 Balloon 资源

```shell
# kubectl apply -f balloon-crd.yaml 
customresourcedefinition.apiextensions.k8s.io/balloons.book.dong.io created
# kubectl apply -f balloon.yaml 
balloon.book.dong.io/my-balloon created
# kubectl get balloons
NAME         STATUS
my-balloon   
```

气球目前还没有到达释放时间，因此它的 `STATUS` 为空。为了更好的看到效果，你需要根据实际情况修改 manifest 中气球的释放时间:

```yaml
apiVersion: "book.dong.io/v1"
kind: Balloon
metadata:
  name: my-balloon
  namespace: default
spec:
  # 根据实际时间修改
  releaseTime: "2023-09-07T17:24:00+08:00"
```

### Step3: 部署 balloon-controller

目前集群中存在我们自定义的气球资源，但是集群中没有控制器来对其执行特定的动作。我们部署气球控制器使其可以控制气球的“释放”:

```shell
# kubectl apply -f balloon-controller-deploy.yaml 
deployment.apps/balloon-controller created
```

等待 Pod Ready:

```shell
# kubectl get pod
NAME                                  READY   STATUS    RESTARTS       AGE
balloon-controller-7884b6489b-z8nv7   1/1     Running   0              97s
```

在到了气球的释放时间后，查看气球状态已变成“释放”：

```shell
# kubectl logs balloon-controller-7884b6489b-pr8dd
2023/09/08 06:52:43 balloon controller started successfully!
2023/09/08 06:52:43 The balloon: my-balloon is created.
2023/09/08 06:52:43 The balloon: my-balloon is updated.
2023/09/08 06:53:03 releasing balloon: my-balloon
2023/09/08 06:53:03 The balloon: my-balloon is updated.
2023/09/08 06:53:03 balloon: my-balloon is released, status is Released
# kubectl get balloons
NAME         STATUS
my-balloon   Released
```

# Build Image

由于本系列教程的主要目的是实现一个极为简单的自定义控制器，控制器代码中对于一些可能遇到的需要处理的错误，或是更为复杂和高级的内容并没有做相应的实现。

这些内容有兴趣的读者可以自行实现，为你的气球添加更多更有趣的玩法。

由于在 Dockerfile 中使用 `go mod` 拉取一些代码库会出现问题，因此我们使用 vendor 形式在容器中编译。在打包之前，使用 `go mod vendor` 更新引用包到本地，然后构建镜像时一同将所有代码依赖放到容器中进行编译。

接下来，你就可以构建自己的 **balloon-controller** 镜像并按照仓库的具体要求打上 tag，再推送到自己的仓库中。

```shell
# cd crd-controller-from-scratch/balloon-controller/
# docker build -t {$YourImageName}:{$Version} .
# docker login {$YourRepository}
# docker push {$YourImageName}:{$Version}
```

最后，将 manifest 中的镜像替换成你的控制器，部署你的气球，等待控制器工作。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  ...
    spec:
      containers:
      - name: balloon-controller
        image: {$YourImageName}:{$Version}
        imagePullPolicy: IfNotPresent
```