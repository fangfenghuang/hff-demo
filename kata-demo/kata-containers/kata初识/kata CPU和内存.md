 [TOC]

# 配置

## configuration.toml

enable_cpu_memory_hotplug（默认false）

default_vcpus=1（默认）

default_maxvcpus

default_memory=2048（默认）

## 通过注释：
```
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
 runtime_type = "io.containerd.kata.v2"
 privileged_without_host_devices = false
 shim_debug = true
 pod_annotations = ["io.katacontainers.*"] #  <-- look here
```

>  io.katacontainers.config.hypervisor.default_vcpus = 5

## SandboxCgroupOnly

默认关闭：此时kata容器非业务负载的花销在另外的kata_overhead cgroup中，并且无限制，这样可能会导致资源无管控，占用过多主机资源，好处是业务不需评估非业务负载所需资源

如果开启，则非业务负载和业务负载在一个pod cgroup中，需要评估所有负载及开销，好处是隔离性更好，资源管控更好。

## PodOverhead

### 打开--feature-gates PodOverhead

> /etc/kubernetes/manifests/kube-apiserver.yaml

1.17默认关闭，1.18默认打开

> Environment="KUBELET_FEATURE=--feature-gates=RotateKubeletServerCertificate=true,VolumeSnapshotDataSource=true,ExpandCSIVolumes=true,VolumePVCDataSource=true,ServiceTopology=true,EndpointSlice=true,PodOverhead=true"

### 启用RuntimeClass准入控制

> --feature-gates=VolumeSnapshotDataSource=true,ExpandCSIVolumes=true,VolumePVCDataSource=true,TTLAfterFinished=true,ServiceTopology=true,EndpointSlice=true,PodOverhead=true

### 设置podOverHead

```
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

# 资源需求和限制

如果 pod 没有 CPU 限制，则增加 default_vcpus 可以提高性能。然而，对于许多线程，增加 CPU 限制并没有帮助，并且会让事情变得更糟。因此，根据经验：**要提高性能，请使用 CPU 限制或default_vcpus 注释，但不能同时使用两者**。 



# vcpu

[https://github.com/kata-containers/kata-containers/blob/main/docs/design/vcpu-handling.md](https://github.com/kata-containers/kata-containers/blob/main/docs/design/vcpu-handling.md)

## CPU热插拔

kata-runtime中复用了**—cpus**选项实现了CPU热插拔的功能，通过统计Pod中所有容器的**—cpus**选项的和，然后确定需要热插多少个CPU到轻量级虚机中。



## cpu绑核

# 内存


## virtio-mem内存热插拔

[https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-use-virtio-mem-with-kata.md](https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-use-virtio-mem-with-kata.md)

使用以下命令将容器内存限制设置为 2g，并将 VM 的内存大小设置为其 default_memory + 2g。

$ sudo crictl update --memory $((2*1024*1024*1024)) $cid

内存资源当前只支持热插，不支持内存热拔。