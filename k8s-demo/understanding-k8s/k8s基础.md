# apiVersion
```
kubectl api-versions
```

- alpha
名称中带有alpha的API版本是进入Kubernetes的新功能的早期候选版本。这些可能包含错误，并且不保证将来可以使用。

- beta
API版本名称中的beta表示测试已经超过了alpha级别，并且该功能最终将包含在Kubernetes中。 虽然它的工作方式可能会改变，并且对象的定义方式可能会完全改变，但该特征本身很可能以某种形式将其变为Kubernetes。

- stable
稳定的apiVersion这些名称中不包含alpha或beta。 它们可以安全使用。

- v1
这是Kubernetes API的第一个稳定版本。 它包含许多核心对象。

- apps/v1
apps是Kubernetes中最常见的API组，其中包含许多核心对象和v1。 它包括与在Kubernetes上运行应用程序相关的功能，如Deployments，RollingUpdates和ReplicaSets。

- autoscaling/v1
此API版本允许根据不同的资源使用指标自动调整容器。此稳定版本仅支持CPU扩展，但未来的alpha和beta版本将允许您根据内存使用情况和自定义指标进行扩展。

- batch/v1
batchAPI组包含与批处理和类似作业的任务相关的对象（而不是像应用程序一样的任务，如无限期地运行Web服务器）。 这个apiVersion是这些API对象的第一个稳定版本。

- batch/v1beta1
Kubernetes中批处理对象的新功能测试版，特别是包括允许您在特定时间或周期运行作业的CronJobs。

- certificates.k8s.io/v1beta1
此API版本添加了验证网络证书的功能，以便在群集中进行安全通信。 您可以在官方文档上阅读更多内容。

- extensions/v1beta1
此版本的API包含许多新的常用Kubernetes功能。 部署，DaemonSets，ReplicaSet和Ingresses都在此版本中收到了重大更改。

- policy/v1beta1
此apiVersion增加了设置pod中断预算和pod安全性新规则的功能

- rbac.authorization.k8s.io/v1
此apiVersion包含Kubernetes基于角色的访问控制的额外功能。这有助于您保护群集


# resourceVersion
这个版本号是一个 K8s 的内部机制，用户不应该假设它是一个数字或者通过比较两个版本号大小来确定资源对象的新旧，唯一能做的就是通过比较版本号相等来确定对象是否是同一个版本（即是否发生了变化）。而 resourceVersion 一个重要的用处，就是来做 update 请求的版本控制。
kube-apiserver 会校验用户 update 请求提交对象中的 resourceVersion 一定要和当前 K8s 中这个对象最新的 resourceVersion 一致，才能接受本次 update。否则，K8s 会拒绝请求，并告诉用户发生了版本冲突（Conflict）。

# Patch 机制
- json patch
在 json patch 中我们要指定操作类型，比如 add 新增还是 replace 替换，另外在修改 containers 列表时要通过元素序号来指定容器。

- merge patch（默认）
merge patch 无法单独更新一个列表中的某个元素，因此不管我们是要在 containers 里新增容器、还是修改已有容器的 image、env 等字段，都要用整个 containers 列表来提交 patch：

- strategic merge patch
在我们 patch 更新 containers 不再需要指定下标序号了，而是指定 name 来修改，K8s 会把 name 作为 key 来计算 merge。
目前 strategic 策略只能用于原生 K8s 资源以及 Aggregated API 方式的自定义资源，对于 CRD 定义的资源对象，是无法使用的。
