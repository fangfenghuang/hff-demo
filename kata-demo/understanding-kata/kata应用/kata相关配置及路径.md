 
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

## 1. [hypervisor.qemu]
```
path = "/usr/bin/qemu-system-x86_64" 指定 qemu 的路径
kernel = "/usr/share/kata-containers/vmlinuz.container" 指定启动内核路径
initrd = "/usr/share/kata-containers/kata-containers-initrd.img" 指定 initrd
image = "/usr/share/kata-containers/kata-containers-centos.img" 指定系统盘，initrd 和 image 不可以同时配置，否则会出错
kernel_params = "" 配置-append 参数，定制虚拟机内核启动参数
firmware = "" 指定固件，影响 qemu 的-bios 参数
machine_accelerators="" virtcontainers/qemu.go 的 getQemuMachine() 进行处理，加到 machine.Options 中，最终是加到-machine 参数中
default_vcpus = 1 默认 vcpu 个数
default_maxvcpus = 0  默认最大 vcpu 个数，设置为 0 时实际上时 240
default_bridges = 1 默认 PCI 桥个数
default_memory = 2048  VM 默认内存大小
memory_slots = 10 内存插槽个数
disable_block_device_use = false 是否禁用块设备
block_device_driver = "virtio-scsi" 块设备驱动，可以是 virtio-scsi、virtio-blk 或 nvdimm
#block_device_cache_set = true
#block_device_cache_direct = true
#block_device_cache_noflush = true  块设备是否设置 cache
enable_iothreads = false  //enable_iothreas 当前仅针对 virtio-scsi 块设备生效
#enable_mem_prealloc = true
#enable_hugepages = true
#file_mem_backend = ""  //这几个配置统一用于 kata-runtime qemu 插件启动虚拟机时的内存配置
#enable_swap=true  //是否允许虚拟机内存 swap，以支持更大的虚拟机密度
#enable_debug=false  //影响 guest kernel 内核启动，在 enable_debug 后，guest kernel 启动项会加上 systemd.show_status true systemd.log_level debug
#disable_nesting_checks = true 虚拟机标志是否 nestedRun，即虚拟化嵌套
msize=8192 9p fs msize 选项
#use_vsock = true 是否使用 vsock
#hotplug_vfio_on_root_bus = true 对于 vfio 设备，会挂在到 root bus 上，否则挂载到 PCI bridge 上, 默认是 false
#disable_vhost_net = true 禁用 vhost_net
#entropy_source= "/dev/urandom" 指定随机数发生器，默认为/dev/urandom,kata 启动虚拟机时会给虚拟机附加一个随机数发生器
#guest_hook_path = "/usr/share/oci/hooks" guest 钩子函数执行路径,用于 OCI
#enable_template = true 默认为 false，enable 后新的虚拟机从模板通过虚拟机克隆方式启动，所有 VM 共享相同的初始化 kernel、initramfs 和 agent 内存
#enable_debug = true 默认为 false，enable 后，shim 将消息发往 system log
#enable_tracing = true 默认为 false，用于跟踪
#diable_new_netns=true 默认为 false，enable 后，runtime 不会再为 shim 和 hypervisor 进程创建一个网络 namespace，在 enable_netmon、网络模式采用 bridged 或者 macvtap 后不能 enable 该选项


```

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

