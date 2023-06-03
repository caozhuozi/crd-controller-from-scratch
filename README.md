# CRD controller from scratch

![kubernetes-balloon logo](./website/static/img/logo.svg)

A book aiming to teach you to implement a Kubernetes custom controller from scratch.

The book is still being written and an English version is not yet available. For Chinese/中文 readers,
the pre-release version is hosted at [github pages](https://caozhuozi.github.io/crd-controller-from-scratch/).


## Plan towards v0.6

I'm currently working towards version 0.6, where I plan to add more details to the existing content to make it easier to read.
The main points that need to be improved are listed below for tracking.

- [ ] we should add mind map for the chapter structures.
- [X] the importance of the concept "GV" should be emphasized.
- [ ] we should mention the representation of GV in the "apiVersion" field.
- [ ] the "internal version" of a resource should be emphasized and expanded.
- [X] we should explain more in "本可不必如此" section.
- [ ] resource "collection" should be emphasized in "Controller" section.
- [ ] it's better to tell the origin of the informer name.
- [ ] we should add a complete example about how to initialize an informer.
- [ ] we must add text details in the "serialization" section.
- [ ] whole picture about the codec flow? more explanation about the codec and the JSON serializer?
- [ ] we should add a separate section for "SchemeBuilder".
- [ ] informer resync( and also resyncPeriod parameter) should be explained.
- [ ] should we introduce the `resourceVersion` when explaining watch and informer?
- [ ] the history of informer is quite interesting. its helpful if we could add it. 
- [ ] "informer is all you need" section needs a major revision. "Reflector" actually doesn't contain 
      a storage but need a storage. the origin of the name "Reflector" may be also helpful to understand Reflector.
- [ ] "RESTClient" section needs a major revision including:
       1) the relation between clientset, CoreV1Client
          and "PodInterface"(this is helpful to explain "ListWatcher" in "informer is all you need" section).
       2) 2 ways to initialize a RESTClient (inside and outside cluster) initialization should be emphasized.
- [ ] ~~we should also mention the difference between stdlib json serializer and the apimachinery one: case-sensitive.~~
- [ ] we should give a separate section for "Request"&"Result" types.
- [ ] we need to explain how to interact with Kubernetes API using raw http package or curl?
      1) need to move http package example from subscript to main body(maybe "restclient" is a good place to put it)
      2) also need to add a curl example
- [ ] we should emphasize that `client-gen` is also not required in "introduction" chapter.
- [ ] the history of "group" and its corresponding concept "apiVersion"(slash) and "core group".
- [ ] we should add the reference to "controller pattern".
- [ ] should we use release page when we refer to kubernetes version?
- [ ] should we use the exact PR name instead of number when we refer to an PR?
- [ ] the history of API group
- [ ] the history of "API object" ?



## Copyright

The work is licensed under a [CC BY-NC-ND 4.0](https://creativecommons.org/licenses/by-nc-nd/4.0/)  international license.







