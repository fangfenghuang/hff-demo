[TOC]
# docker 兼容问题

- kata已去掉docker cli支持，请使用crictl命令

- DinD问题：参考[[流水线适配问题]]



# 使用限制：

## 不支持subPaths(emptyDir )使用

## kata不支持host网络

​		一些使用主机网络的k8s组件和业务无法使用kata容器，所以runc（containerd）必须保留作为默认运行时，而kata-container作为可选运行时给特定负载使用。

## 不支持网络命名空间共享

​		Docker 支持容器使用docker run --net=containers语法加入另一个容器命名空间的能力。这允许多个容器共享一个公共网络命名空间和放置在网络命名空间中的网络接口。Kata Containers 不支持网络命名空间共享。如果将 Kata 容器设置为共享runc容器的网络命名空间，则运行时会有效地接管分配给命名空间的所有网络接口并将它们绑定到 VM。因此，runc容器失去其网络连接。



# docker切containerd影响




# 测试默认配置文件
如果存在/etc/kata-containers/configuration.toml，测试ctr run和kata pod是否使用了改配置文件


/etc/kata-containers/configuration.toml：设置default_cpus=3

* 没有/opt/kata/share/defaults/kata-containers/configuration-qemu.toml文件，kata pod可以起来，但是ctr run不行
* 增加/etc/kata-containers/configuration.toml后，kata pod cpu数变成3，但是ctr run cpu为1

