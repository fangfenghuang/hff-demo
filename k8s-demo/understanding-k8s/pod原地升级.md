    所谓原地升级模式，就是在应用升级过程中避免将整个 Pod 对象删除、新建，而是基于原有的 Pod 对象升级其中某一个或多个容器的镜像版本

原地升级为发布效率带来了以下优化点：
节省了调度的耗时，Pod 的位置、资源都不发生变化；
节省了分配网络的耗时，Pod 还使用原有的 IP；
节省了分配、挂载远程盘的耗时，Pod 还使用原有的 PV（且都是已经在 Node 上挂载好的）；
节省了大部分拉取镜像的耗时，因为 Node 上已经存在了应用的旧镜像，当拉取新版本镜像时只需要下载很少的几层 layer。

# 分析
## Kubelet 针对 Pod 容器的版本管理
每个 Node 上的 Kubelet，会针对本机上所有 Pod.spec.containers 中的每个 container 计算一个 hash 值，并记录到实际创建的容器中。

如果我们修改了 Pod 中某个 container 的 image 字段，kubelet 会发现 container 的 hash 发生了变化、与机器上过去创建的容器 hash 不一致，而后 kubelet 就会把旧容器停掉，然后根据最新 Pod spec 中的 container 来创建新的容器。
## Pod 更新限制
对于一个已经创建出来的 Pod，在 Pod Spec 中只允许修改 containers/initContainers 中的 image 字段，以及 activeDeadlineSeconds 字段。对 Pod Spec 中所有其他字段的更新，都会被 kube-apiserver 拒绝。

## containerStatuses 上报
一个 Pod 中 spec 和 status 的 image 字段不一致，并不意味着宿主机上这个容器运行的镜像版本和期望的不一致。

## ReadinessGate 控制 Pod 是否 Ready
目前 kubelet 判定一个 Pod 是否 Ready 的两个前提条件：

1. Pod 中容器全部 Ready（其实对应了 ContainersReady condition 为 True）；
2. 如果 pod.spec.readinessGates 中定义了一个或多个 conditionType，那么需要这些 conditionType 在 pod.status.conditions 中都有对应的 status: "true" 的状态。


# 方案
。。。


# 参考资料
https://developer.aliyun.com/article/765421
https://jimmysong.io/kubernetes-handbook/practice/in-place-update.html
