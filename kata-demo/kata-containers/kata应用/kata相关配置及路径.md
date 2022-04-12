 
[TOC]
# containerd配置

> /etc/containerd/config.toml
   /etc/systemd/system/kubelet.service.d/0-containerd.conf

## containerd数据路径
>root = "/app/docker/containerd"
state = "/app/docker/run/containerd"


## containerd插件数据
>ctr plugin ls

 

# kata配置configuration.toml

默认的配置文件位于/opt/kata/share/defaults/kata-containers/configuration.toml，如果/etc/kata-containers/configuration.toml的配置文件存在，则会替代默认的配置文件。

>[root@rqy-k8s-1 hff]# kata-runtime --kata-show-default-config-paths
/etc/kata-containers/configuration.toml
/opt/kata/share/defaults/kata-containers/configuration.toml



> [hypervisor.qemu]
> use_vsock ：使用vsocks与agent直接通信（前提支持vsocks），默认false 
> [runtime]
> enable_cpu_memory_hotplug ：使能cpu和内存热插拔，默认false
> [agent.kata]
> debug_console_enabled = true
> [hypervisor.qemu]
> sed -i -e 's/^kernel_params = "\(.*\)"/kernel_params = "\1 agent.debug_console"/g' "${kata_configuration_file}"

  	修改后新容器生效


## ctr依赖/configuration-qemu.toml路径问题
/etc/kata-containers/configuration.toml已存在，为测试删除了默认配置文件，但是containerd配置保留

>[root@localhost ~]# cat  /etc/containerd/config.toml 
      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
        runtime_type = "io.containerd.kata.v2"
        privileged_without_host_devices = true
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata.options]
           ConfigPath = "/etc/kata-containers/configuration.toml"

>[root@localhost ~]# kubectl get runtimeclasses.node.k8s.io kata-containers -o yaml | grep handler
handler: kata

>[root@localhost ~]# ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hfftest dmesg 
ctr: Cannot find usable config file (file /opt/kata/share/defaults/kata-containers/configuration-qemu.toml does not exist): not found

## kata-deploy会修改containerd配置问题
即便部署前已修改：
>[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
        runtime_type = "io.containerd.kata.v2"
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata.options]
           ConfigPath = "/opt/kata/share/defaults/kata-containers/configuration.toml" ##这里还是默认配置，所以ctr还是使用的旧配置？


https://github.com/kata-containers/runtime/issues/1091
https://github.com/kata-containers/kata-containers/blob/main/src/runtime/README.md

## 规避方法
（无效）
>设置KATA_CONF_FILE环境变量
[root@localhost ~]# export KATA_CONF_FILE=/etc/kata-containers/configuration.toml
[root@localhost ~]# env | grep KATA_CONF_FILE
KATA_CONF_FILE=/etc/kata-containers/configuration.toml
[root@localhost ~]# ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hfftest dmesg
ctr: Cannot find usable config file (file /opt/kata/share/defaults/kata-containers/configuration-qemu.toml does not exist): not found

（无效）
>[root@localhost ~]# ctr -n k8s.io run --env KATA_CONF_FILE=/etc/kata-containers/configuration.toml --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hfftest dmesg
ctr: Cannot find usable config file (file /opt/kata/share/defaults/kata-containers/configuration-qemu.toml does not exist): not found


（无效）
>修改containerd配置：
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
  runtime_type = "io.containerd.kata.v2"
  privileged_without_host_devices = true
  pod_annotations = ["io.katacontainers.*"]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata.options]
    ConfigPath = "/etc/kata-containers/configuration.toml"


（有效）
>[root@localhost ~]# cat /usr/local/bin/containerd-shim-kata-qemu-v2
>#!/usr/bin/env bash
KATA_CONF_FILE=/opt/kata/share/defaults/kata-containers/configuration-qemu.toml /opt/kata/bin/containerd-shim-kata-v2 "$@"
改为：
[root@localhost ~]# cat /usr/local/bin/containerd-shim-kata-qemu-v2
#!/usr/bin/env bash
KATA_CONF_FILE=/etc/kata-containers/configuration.toml /opt/kata/bin/containerd-shim-kata-v2 "$@"


kata部署完成后，新增一个/usr/local/bin/containerd-shim-kata-containers-v2，然后软连接 /usr/local/bin/containerd-shim-kata-v2指向该文件


# crictl

crictl 默认连接到 unix:///var/run/dockershim.sock

 

> [root@rqy-k8s-3 fio-iperf]# cat /etc/crictl.yaml
> runtime-endpoint: unix:///run/containerd/containerd.sock
> timeout: 0
> debug: false





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

## virtio_fs_cache = "auto"
- none
Metadata, data, and pathname lookup are not cached in guest. They are
always fetched from host and any changes are immediately pushed to host.
- auto
Metadata and pathname lookup cache expires after a configured amount of
 time (default is 1 second). Data is cached while the file is open (close to open consistency).
- always
Metadata, data, and pathname lookup are cached in guest and never expire.

## 共享内存目录file_mem_backend



# 修改sanbox(虚拟机)配置

https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-set-sandbox-config-kata.md

## containerd配置：

​     [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]

​      runtime_type = "io.containerd.kata.v2"

​      pod_annotations = ["io.katacontainers.*"]

​      container_annotations = ["io.katacontainers.*"]

## 受限注释：

一些注释是*受限*的，这意味着配置文件指定了可接受的值。目前，出于安全原因，仅管理程序注释受到限制，目的是控制 Kata Containers 运行时将代表您启动哪些二进制文件。

> configuration.toml：
>
> enable_annotations = []


