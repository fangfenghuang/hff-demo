Cgroups 是 Control Groups 的缩写，由 Linux内核提供。用于限制、记录和隔离进程组使用的物理资源（CPU、内存、i/o）。

```
[root@localhost ~]# mount -t cgroup
cgroup on /sys/fs/cgroup/systemd type cgroup (rw,nosuid,nodev,noexec,relatime,xattr,release_agent=/usr/lib/systemd/systemd-cgroups-agent,name=systemd)
cgroup on /sys/fs/cgroup/memory type cgroup (rw,nosuid,nodev,noexec,relatime,memory)
cgroup on /sys/fs/cgroup/net_cls,net_prio type cgroup (rw,nosuid,nodev,noexec,relatime,net_prio,net_cls)
cgroup on /sys/fs/cgroup/freezer type cgroup (rw,nosuid,nodev,noexec,relatime,freezer)
cgroup on /sys/fs/cgroup/cpu,cpuacct type cgroup (rw,nosuid,nodev,noexec,relatime,cpuacct,cpu)
cgroup on /sys/fs/cgroup/cpuset type cgroup (rw,nosuid,nodev,noexec,relatime,cpuset)
cgroup on /sys/fs/cgroup/devices type cgroup (rw,nosuid,nodev,noexec,relatime,devices)
cgroup on /sys/fs/cgroup/perf_event type cgroup (rw,nosuid,nodev,noexec,relatime,perf_event)
cgroup on /sys/fs/cgroup/hugetlb type cgroup (rw,nosuid,nodev,noexec,relatime,hugetlb)
cgroup on /sys/fs/cgroup/blkio type cgroup (rw,nosuid,nodev,noexec,relatime,blkio)
cgroup on /sys/fs/cgroup/pids type cgroup (rw,nosuid,nodev,noexec,relatime,pids)
```

cpuset、cpu、 memory子系统

cfs_period和cfs_quota:限制进程在长度为cfs_period的一段时间内，只能被分配到总量为cfs_quota的CPU时间,CPU quota默认没有任何限制（即：-1）,CPU period默认100 ms（100000 us）
```
[root@localhost ~]# cat /sys/fs/cgroup/cpu/docker/cpu.cfs_period_us
100000
[root@localhost ~]# cat /sys/fs/cgroup/cpu/docker/cpu.cfs_quota_us
-1
```

注：
可压缩资源：cpu等(当可压缩资源不足时，Pod只会“饥饿”，但不会退出)
非可压缩资源： mem(当不可压缩资源不足时，Pod就会因为OOM（Out-Of-Memory）被内核杀掉。)



# K8s中的Cgroups
kubelet启动时，它会根据需要创建一个4层的cgroup树。
- Node(root) cgroup
- QoS cgroup
- Pod cgroup
- Container cgroup

Kubelet运行在Kubernetes集群的工作节点上，管理pod的生命周期，连接到CRI、CNI、CSI等运行时接口，通过cgroup按照配置向pod提供资源，并在使用超过时删除或OOMKilled pod。

kubelet中的所有cgroup操作都是由其内部的containerManager模块实现的，该模块通过cgroup(从下到上)对资源使用进行逐层限制

>container-> pod-> qos -> node(root)。

## 容器级别的Cgroup
- 如果容器超出其内存限制，则可能终止该容器。如果它是可重新启动的，kubelet将重新启动它，就像处理任何其他类型的运行时故障一样。
- 如果容器超出了它的内存请求，那么每当节点耗尽内存时，它的Pod很可能会被逐出。
- 容器可能被允许，也可能不被允许在较长时间内超过其CPU限制。但是，它不会因为CPU占用过多而被杀死。

## Pod级别的Cgroup
每个pod都有一定的开销资源，比如由sandbox容器或docker containerd-shim(在1.22中删除)使用的资源。另外，pod在指定内存类型的卷时也会占用内存资源。

- 计算所有pod容器的CPU请求/限制和内存限制，并将CPU请求/限制转换为cpuhares和cpuQuota。
>pod<UID>/cpu.shares = sum(pod.spec.containers.resources.requests[cpu])
pod<UID>/cpu.cfs_quota_us = sum(pod.spec.containers.resources.limits[cpu])
pod<UID>/memory.limit_in_bytes = sum(pod.spec.containers.resources.limits[memory])

- 只有pod cgroup的cpu。如果其中一个容器只指定请求而不指定限制，则将设置共享值，而不设置其限制值。

- 只设置cpu。当无容器指定请求值或限制值时共享，即pod在资源空闲时使用所有节点资源，在资源受限时不获取任何资源，符合低优先级原则

## QoS级别的Cgroup
- Guaranteed(保证型)::由请求设置的值等于由限制设置的值。
当Pod里的每一个Container都同时设置了requests和limits，并且requests和limits值相等的时候，这个Pod就属于Guaranteed类别
- Burstable(突增型):请求设置的值小于限制设置的值，但是!= 0。
而当Pod不满足Guaranteed的条件，但至少有一个Container设置了requests。那么这个Pod就会被划分到Burstable类别
- BestEffort(尽力而为型):请求和限制设置的值都是0。
优先级是
而如果一个Pod既没有设置requests，也没有设置limits，那么它的QoS类别就是BestEffort

>Guaranteed > Burstable > BestEffort

## 节点Node级别的Cgroup
- 业务流程使用的资源，即pod使用的资源。
- Kubernetes组件使用的资源，如kubelet和docker。
- 系统组件使用的资源，如logind、journald和其他进程。


# 资源限制
在调度的时候，kube-scheduler只会按照requests的值进行计算。而在真正设置Cgroups限制的时候，kubelet则会按照limits的值来进行设置。

- 当你指定了requests.cpu=250m之后，相当于将Cgroups的cpu.shares的值设置为(250/1000)*1024。而当你没有设置requests.cpu的时候，cpu.shares默认则是1024

- 如果你指定了limits.cpu=500m之后，则相当于将Cgroups的cpu.cfs_quota_us的值设置为(500/1000)*100ms，而cpu.cfs_period_us的值始终是100ms。这样，Kubernetes就为你设置了这个容器只能用到CPU的50%。

- 对于内存来说，当你指定了limits.memory=128Mi之后，相当于将Cgroups的memory.limit_in_bytes设置为128 * 1024 * 1024。


>Kubernetes这种对CPU和内存资源限额的设计，实际上参考了Borg论文中对“动态资源边界”的定义，既：容器化作业在提交时所设置的资源边界，并不一定是调度系统所必须严格遵守的，这是因为在实际场景中，大多数作业使用到的资源其实远小于它所请求的资源限额。

## Eviction
目前，Kubernetes为你设置的Eviction的默认阈值如下所示：

memory.available<100Mi
nodefs.available<10%
nodefs.inodesFree<5%
imagefs.available<15%
当然，上述各个触发条件在kubelet里都是可配置的。比如下面这个例子：

kubelet --eviction-hard=imagefs.available<10%,memory.available<500Mi,nodefs.available<5%,nodefs.inodesFree<5% --eviction-soft=imagefs.available<30%,nodefs.available<10% --eviction-soft-grace-period=imagefs.available=2m,nodefs.available=2m --eviction-max-pod-grace-period=600
在这个配置中，你可以看到Eviction在Kubernetes里其实分为Soft和Hard两种模式。

# 容器是一个“单进程”模型
