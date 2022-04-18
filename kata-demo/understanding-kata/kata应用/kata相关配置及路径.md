 
[TOC]
# containerd配置

> /etc/containerd/config.toml
> /etc/systemd/system/kubelet.service.d/0-containerd.conf

## containerd数据路径
> root = "/app/docker/containerd"
> state = "/app/docker/run/containerd"


## containerd插件数据
> ctr plugin ls

 

# kata配置configuration.toml

默认的配置文件位于/opt/kata/share/defaults/kata-containers/configuration.toml，如果/etc/kata-containers/configuration.toml的配置文件存在，则会替代默认的配置文件。

```bash
[root@rqy-k8s-1 hff]# kata-runtime --kata-show-default-config-paths
/etc/kata-containers/configuration.toml
/opt/kata/share/defaults/kata-containers/configuration.toml
```

```bash
[hypervisor.qemu]
use_vsock ：使用vsocks与agent直接通信（前提支持vsocks），默认false 
[runtime]
enable_cpu_memory_hotplug ：使能cpu和内存热插拔，默认false
[agent.kata]
debug_console_enabled = true
[hypervisor.qemu]
sed -i -e 's/^kernel_params = "\(.*\)"/kernel_params = "\1 agent.debug_console"/g' "${kata_configuration_file}"
```

修改后新容器生效

## 注意配置文件问题
https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/containerd-kata.md


# crictl

crictl 默认连接到 unix:///var/run/dockershim.sock

```bash
[root@rqy-k8s-3 fio-iperf]# cat /etc/crictl.yaml
runtime-endpoint: unix:///run/containerd/containerd.sock
timeout: 0
debug: false
```




# 镜像配置：

> [plugins."io.containerd.grpc.v1.cri".registry]

**镜像存储路径：**
> /var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/
   原docker: 
   /app/docker/overlay2


# 存储配置
## 存储路径：
/run/kata-containers/shared/sandboxes/
/run/vc/vm/
/run/vc/sbs/

```bash
[root@localhost ~]# find / -name hfftest0413-etc
/run/kata-containers/shared/sandboxes/3b54b3b02fc7f6905d01aedfc4eb209cfb11fd9136006ed6e11e1e26c0f48562/mounts/c7c33d3c7666933c6f1c182bb49bf850c5ca99f08b4595b0f37e6f817bb52768/rootfs/etc/hfftest0413-etc
/run/kata-containers/shared/sandboxes/3b54b3b02fc7f6905d01aedfc4eb209cfb11fd9136006ed6e11e1e26c0f48562/shared/c7c33d3c7666933c6f1c182bb49bf850c5ca99f08b4595b0f37e6f817bb52768/rootfs/etc/hfftest0413-etc
/run/containerd/io.containerd.runtime.v2.task/k8s.io/c7c33d3c7666933c6f1c182bb49bf850c5ca99f08b4595b0f37e6f817bb52768/rootfs/etc/hfftest0413-etc
/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/199015/fs/etc/hfftest0413-etc
```

# 与存储相关的参数

## virtio_fs_cache = "auto"
- none
Metadata, data, and pathname lookup are not cached in guest. They are
always fetched from host and any changes are immediately pushed to host.
- auto
Metadata and pathname lookup cache expires after a configured amount of
 time (default is 1 second). Data is cached while the file is open (close to open consistency).
- always
Metadata, data, and pathname lookup are cached in guest and never expire.

## virtio_fs_cache_size = 0


## 共享内存目录file_mem_backend



# 修改sanbox(虚拟机)配置

https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-set-sandbox-config-kata.md

## containerd配置：
```bash
​ [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
​  runtime_type = "io.containerd.kata.v2"
​  pod_annotations = ["io.katacontainers.*"]
​  container_annotations = ["io.katacontainers.*"]
```
## 受限注释：

一些注释是*受限*的，这意味着配置文件指定了可接受的值。目前，出于安全原因，仅管理程序注释受到限制，目的是控制 Kata Containers 运行时将代表您启动哪些二进制文件。

> configuration.toml：
  enable_annotations = []


# runtime 名字怎么写？
在 K8s 和 containerd 中，我们会看到很多用于设置 runtime 的地方，比如 RuntimeClass 、Pod 的 runtimeClassName 定义，以及 ctr run --runtime io.containerd.run.kata.v2 和 crictl runp -r kata ，里面都有参数指定运行时的名字。

Pod 的 runtimeClassName 属性会查找同名的 RuntimeClass 资源
根据 该资源的 handler ，在 containerd 的配置文件查找相应的运行时（ [plugins.cri.containerd.runtimes.${HANDLER_NAME}] ）。
一般情况下 containerd 配置会像这样：

>[plugins.cri.containerd.runtimes.kata]
  runtime_type = "io.containerd.kata.v2"
- ctr 命令使用的是 containerd 配置文件中的 runtime_type 属性（ containerd 用）。
- crictl 和 K8s（实际也是 CRI 接口） 使用的是 containerd 配置中的 HANDLER_NAME（ CRI 用）。

默认情况下，containerd 会根据 runtime_type 按规则对应到具体的运行时的可执行文件名。比如 Kata Containers(io.containerd.kata.v2) 运行时最终会转换为 containerd-shim-kata-v2 命令，该命令默认安装在 /usr/local/bin/containerd-shim-kata-v2。

```
[root@localhost ~]# cat /usr/local/bin/containerd-shim-kata-v2
#!/usr/bin/env bash
KATA_CONF_FILE=/opt/kata/share/defaults/kata-containers/configuration-qemu.toml /opt/kata/bin/containerd-shim-kata-v2 "$@"
```

