# CRD controller from scratch

![kubernetes-balloon logo](./website/static/img/logo.svg)

A book aiming to teach you to implement a Kubernetes custom controller from scratch.

The book is still being written and an English version is not yet available. For Chinese/中文 readers,
the pre-release version is hosted at [github pages](https://caozhuozi.github.io/crd-controller-from-scratch/).


## Plan towards v0.6

I'm currently working towards version 0.6, where I plan to add more details to the existing content to make it easier to read.
The main points that need to be improved are listed below for tracking.

- [ ] ~~we should add mind map for the chapter structures.~~
- [X] the importance of the concept "GV" should be emphasized.
- [X] we should mention the representation of GV in the "apiVersion" field.
- [X] the "internal version" of a resource should be emphasized and expanded.
- [X] ~~we should explain more in "本可不必如此" section.~~ <- I already have a major revision on "Kubernetes API" section.
- [ ] ~~resource "collection" should be emphasized in "Controller" section.~~
- [X] it's better to tell the origin of the informer name.
- [ ] ~~we should add a complete example about how to initialize an informer.~~ <- decide to move this to 6.1
- [X] we must add text details in the "serialization" section.
- [X] whole picture about the codec flow? more explanation about the codec and the JSON serializer?
- [X] we should add a separate section for "SchemeBuilder".
- [X] informer resync( and also resyncPeriod parameter) should be explained.
- [ ] ~~should we introduce the `resourceVersion` when explaining watch and informer?~~
- [X] the history of informer is quite interesting. its helpful if we could add it. 
- [X] "informer is all you need" section needs a major revision. "Reflector" actually doesn't contain 
      a storage but need a storage. the origin of the name "Reflector" may be also helpful to understand Reflector.
- [X] "RESTClient" section needs a major revision including:
       1) the relation between clientset, CoreV1Client
          and "PodInterface"(this is helpful to explain "ListWatcher" in "informer is all you need" section).
       2) 2 ways to initialize a RESTClient (inside and outside cluster) initialization should be emphasized.
- [ ] ~~we should also mention the difference between stdlib json serializer and the apimachinery one: case-sensitive.~~
- [X] we should give a separate section for "Request"&"Result" types.
- [ ] ~~we need to explain how to interact with Kubernetes API using raw http package or curl?~~
      1) need to move http package example from subscript to main body(maybe "restclient" is a good place to put it)
      2) also need to add a curl example
- [X] we should emphasize that `client-gen` is also not required in "introduction" chapter.
- [X] the history of "group" and its corresponding concept "apiVersion"(slash) and "core group".
- [ ] ~~we should add the reference to "controller pattern".~~
- [ ] ~~should we use the exact PR name instead of number when we refer to an PR?~~
- [ ] the history of API group <- need to explain more. maybe move to the following version of v0.6.0.
- [ ] ~~the history of "API object" ?~~
- [ ] ~~Do we need to mention other cache types or even custom cache when initializing `refector`?~~
- [X] a major revision of "Kubernetes API" section
- [ ] ~~kind & json schema~~
- [X] To make "virtual resource" more clear, we need to explain deeply about what is "SelfSubjectAccessReview" and what it is for.
- [X] the group name is not limited to a single word
- [X] the group name common dot pattern
- [X] Do we need to put subsection "Kubernetes对象的期望状态与实际状态" into "从Kubernetes API谈起" ?
- [X] reference official document in "subresource" section
- [X] **Spec and Status for Kubernetes Objects**
- [X] hard decision: we need to do a major revision on "serializer" section. 
      the section should be mainly divided into 3 parts: 1) resource version conversion 2) before v1.2 and 3) after v1.2 <- I'm currently working on this.
- [X] the history of `apimachinery` and `client-go`
- [X] update "introduction": 
      limit about why we ignore CRD and those  manifests
      style: history

### Remaining TODO List 
The remaining good points are grouped together for easy tracking.
- [ ] Do we also need to add the history of "spec" and "status" fields?
- [ ] Do we need to find the history of "hub and spoke" conversion model?
- [ ] should we use release page when we refer to kubernetes version?


## Copyright

The work is licensed under a [CC BY-NC-ND 4.0](https://creativecommons.org/licenses/by-nc-nd/4.0/)  international license.







