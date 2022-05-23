 [TOC]

 资源的配置应分为两部分：对轻量级虚拟机的资源配置，即Host资源配置；对虚拟机内容器的配置，即Guest容器资源配置。


# 注意事项（总结）
- 一个Pod的最小规格是1C 256M，当低于 256M 时，会重置为 2G。，支持的最大内存规格是256GB。如果用户分配的内存规格超过256GB，可能会出现未定义的错误，安全容器暂不支持超过256GB的大内存场景。
- Kata 配置文件中默认的 VM 大小为 1C 2G（不设置limit不是无限制）
- Kata VM 中额外的资源是通过 hotplug 的方式实现，资源目前特指 CPU 和 Memory 两种；终 VM 的资源大小为 limit + default，其中 limit 为 Pod 声明的 limit，不包含 overhead 在内，default 为Kata 配置文件中的基础 VM 大小。（VM的大小limit+default如果超出主机资源，pod会创建失败（CPU可能会在使用过程中pod异常））
- 如果pod没有CPU限制，则增加default_vcpus可以提高性能。然而，对于许多线程，增加CPU限制并没有帮助，并且会让事情变得更糟。因此，根据经验：要提高性能，请使用 CPU 限制或default_vcpus 注释，但不能同时使用两者。 
- SandboxCgroupOnly默认为false，此时kata容器虚拟机开销可能会占用过多主机资源。（实测fio测试可能会导致k8s节点异常）
- 如果不设置request，则request的值和limit默认相等
- Kata 对于资源并不是完全占用，不同的 Kata VM 之间会存在资源抢占现象。在此方面，Kata Containers 和传统容器的设计理念相同。



# 一些需要知道的配置

## configuration.toml

enable_cpu_memory_hotplug（默认false）

default_vcpus=1（默认）

default_maxvcpus

default_memory=2048（默认）

## 通过注释修改配置：
```bash
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
 runtime_type = "io.containerd.kata.v2"
 privileged_without_host_devices = false
 shim_debug = true
 pod_annotations = ["io.katacontainers.*"] #  <-- look here
```

>  io.katacontainers.config.hypervisor.default_vcpus = 5

## SandboxCgroupOnly

默认禁用：此时kata容器非业务负载开销在另外的kata_overhead cgroup中，并且无限制，这样可能会导致资源无管控，占用过多主机资源，好处是业务不需评估非业务负载所需资源

如果开启，则非业务负载和业务负载在一个pod cgroup中，需要评估所有负载及开销，好处是隔离性更好，资源管控更好。

## PodOverhead
PodOverhead作用于 Kata Containers 的额外开销，而不是业务负载

### 打开--feature-gates PodOverhead

> /etc/kubernetes/manifests/kube-apiserver.yaml

1.17默认关闭，1.18默认打开

> Environment="KUBELET_FEATURE=--feature-gates=RotateKubeletServerCertificate=true,VolumeSnapshotDataSource=true,ExpandCSIVolumes=true,VolumePVCDataSource=true,ServiceTopology=true,EndpointSlice=true,PodOverhead=true"

### 启用RuntimeClass准入控制

> --feature-gates=VolumeSnapshotDataSource=true,ExpandCSIVolumes=true,VolumePVCDataSource=true,TTLAfterFinished=true,ServiceTopology=true,EndpointSlice=true,PodOverhead=true

### 设置podOverHead

```yaml
kind: RuntimeClass
apiVersion: node.k8s.io/v1beta1
metadata:
  name: kata-containers
handler: kata
overhead:
  podFixed:
  memory: "100Mi"
  cpu: "100m"
```



# vcpu

[https://github.com/kata-containers/kata-containers/blob/main/docs/design/vcpu-handling.md](https://github.com/kata-containers/kata-containers/blob/main/docs/design/vcpu-handling.md)

## CPU热插拔

kata-runtime中复用了**—cpus**选项实现了CPU热插拔的功能，通过统计Pod中所有容器的**—cpus**选项的和，然后确定需要热插多少个CPU到轻量级虚机中。



## cpu绑核

# 内存


## virtio-mem内存热插拔（仅支持QEMU）
[https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-use-virtio-mem-with-kata.md](https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-use-virtio-mem-with-kata.md)

```
$ sudo sed -i -e 's/^#enable_virtio_mem.*$/enable_virtio_mem = true/g' /etc/kata-containers/configuration.toml
```



使用以下命令将容器内存限制设置为 2g，并将 VM 的内存大小设置为其 default_memory + 2g。
```bash
$ sudo crictl update --memory $((2*1024*1024*1024)) $cid
```

内存资源当前只支持热插，不支持内存热拔。


# 资源限制

## cpu/mem资源限制

### 不设置limit:
默认使用default设置的cpu/mem限制1C2G

### overhead

### 设置limit

容器业务（pod）最大使用上限：limit

最终VM的资源大小为：limit+default（lscpu、free -h）

VM的的最大使用量（开启SandboxCgroupOnly）：overhead+limit(memory.limit_in_bytes)(describe node) 

如果不设置request，则request的值和limit默认相等

