 [TOC]



系统管理程序 hypervisor将启动一个虚拟机，该虚拟机包括最小的 虚拟机内核和虚拟机镜像。
	
# 配置
```bash
[hypervisor.qemu]
path = "/opt/kata/bin/qemu-system-x86_64"
kernel = "/opt/kata/share/kata-containers/vmlinux.container"
image = "/opt/kata/share/kata-containers/kata-containers.img"
machine_type = "q35"
[Kernel]
 Path = "/opt/kata/share/kata-containers/vmlinux-5.15.23-89"
 Parameters = "systemd.unit=kata-containers.target systemd.mask=systemd-networkd.service systemd.mask=systemd-networkd.socket scsi_mod.scan=none agent.debug_console agent.debug_console_vport=1026"
 [Image]
 Path = "/opt/kata/share/kata-containers/kata-clearlinux-latest.image"
 [Initrd]
 Path = ""
```
 

# Guest kernel

用于启动VM。Kata-container高度优化了内核启动时间和极小的内存占用，只用于一个容器的运行。

# Guest image

支持基于initrd和rootfs(image)的最小镜像，默认包中提供一个镜像和一个initrd，他们都是通过osbuilder生成的。

## image type(rootfs)
默认的 root filesystem image (有时称为mini O/S)是一个基于 Clear Linux 的高度优化的容器引导系统。它提供了一个极小的环境，并有一个高度优化的引导路径。


## initrd type
压缩的 cpio(1) 文件，从 rootfs （被载入内存并作为 linux 启动过程的一部分被使用）创建。



Kata 运行时配置文件中的initrd和image选项之一**必须**设置，但**不能同时**设置。选项之间的主要区别在于initrd(10MB+) 的大小明显小于 rootfs image(100MB+)。

通过initrd=和image=配置决定使用哪一个类型

# 最小镜像极度简化

[https://github.com/kata-containers/kata-containers/issues/2010](https://github.com/kata-containers/kata-containers/issues/2010)



# 内核构建build-kernel.sh

[https://github.com/kata-containers/kata-containers/tree/main/tools/packaging/kernel](https://github.com/kata-containers/kata-containers/tree/main/tools/packaging/kernel)

例子：
```bash
$ ./build-kernel.sh -v 5.10.25 -g nvidia -f -d setup
· -v 5.10.25：指定来宾内核版本。
· -g nvidia: 构建一个支持 Nvidia GPU 的来宾内核。
· -f:.config即使内核目录已经存在也强制生成文件。
· -d: 启用 bash 调试模式。
```

添加补丁：${GOPATH}/src/github.com/kata-containers/kata-containers/tools/packaging/kernel/patches/

内核配置：${GOPATH}/src/github.com/kata-containers/kata-containers/tools/packaging/kernel/configs/

# osbuilder

[https://github.com/kata-containers/kata-containers/tree/main/tools/osbuilder](https://github.com/kata-containers/kata-containers/tree/main/tools/osbuilder)

# 修改内核参数
```bash
[Kernel]
 Path = "/opt/kata/share/kata-containers/vmlinux-5.15.23-89"
 Parameters = "systemd.unit=kata-containers.target systemd.mask=systemd-networkd.service systemd.mask=systemd-networkd.socket scsi_mod.scan=none agent.debug_console agent.debug_console_vport=1026"
 ```

# 加载内核模块

[https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-load-kernel-modules-with-kata.md](https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-load-kernel-modules-with-kata.md)

## 使用 Kata 配置文件configuration.toml（全局）

> 
kernel_modules =[ “ e1000e InterruptThrottleRate=3000,3000,3000 EEE=1 ” , “ i915 ” ]

## 使用注释

```bash
annotations:
  io.katacontainers.config.agent.kernel_modules: "e1000e EEE=1; i915"spec:
```

# 使用 Kata 设置sysctl
[https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-use-sysctls-with-kata.md](https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-use-sysctls-with-kata.md)

sysctl 使用 pod 的 securityContext 在 pod 上设置。securityContext 适用于同一 pod 中的所有容器。
```yaml
apiVersion: v1kind: Podmetadata:
name: sysctl-examplespec:
securityContext:
  sysctls:
   - name: kernel.shm_rmid_forced
     value: "0"
   - name: net.ipv4.route.min_pmtu
     value: "552"
   - name: kernel.msgmax
     value: "65536"
 ...
 ```

所有安全sysctls默认被开启

使用不安全的 sysctls，集群管理员需要允许这些：

```
$ kubelet --allowed-unsafe-sysctls 'kernel.msg*,net.ipv4.route.min_pmtu' ...
```
