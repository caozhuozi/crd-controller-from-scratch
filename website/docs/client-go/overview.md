---
sidebar_position: 0
id: client-go
---
# Overview

在前一章，我们主要介绍了kind的概念以及`apimachinery`的序列化原理。
在本章中，我们将正式开始介绍`client-go`库中与实现控制器有关的重要组件。
正如我们之前就提及的那样，理论上我们可以不使用`client-go`直接与Kubernetes API交互。
但这并不是本书所聚焦的目标。
我们仍然想复用`client-go`库中官方开发者提供的那些"优雅"的组件。