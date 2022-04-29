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


## 不支持cpuset-cpus
https://github.com/kata-containers/runtime/issues/1079
docker run -d --cpuset-cpus= “ 0-1 ”   ubuntu sleep 30000



# docker切containerd影响




# 测试默认配置文件
如果存在/etc/kata-containers/configuration.toml，测试ctr run和kata pod是否使用了改配置文件


/etc/kata-containers/configuration.toml：设置default_cpus=3

* 没有/opt/kata/share/defaults/kata-containers/configuration-qemu.toml文件，kata pod可以起来，但是ctr run不行
* 增加/etc/kata-containers/configuration.toml后，kata pod cpu数变成3，但是ctr run cpu为1


## ctr依赖/configuration-qemu.toml路径问题
```bash
/etc/kata-containers/configuration.toml已存在，为测试删除了默认配置文件，但是containerd配置保留

[root@localhost ~]# cat  /etc/containerd/config.toml 
      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
        runtime_type = "io.containerd.kata.v2"
        privileged_without_host_devices = true
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata.options]
           ConfigPath = "/etc/kata-containers/configuration.toml"

[root@localhost ~]# kubectl get runtimeclasses.node.k8s.io kata-containers -o yaml | grep handler
handler: kata

[root@localhost ~]# ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hfftest dmesg 
ctr: Cannot find usable config file (file /opt/kata/share/defaults/kata-containers/configuration-qemu.toml does not exist): not found


[root@rqy-k8s-1 kbuser]# ll /usr/local/bin/containerd-shim-kata-v2
lrwxrwxrwx. 1 root root 43 Mar 18 11:31 /usr/local/bin/containerd-shim-kata-v2 -> /usr/local/bin/containerd-shim-kata-qemu-v2
```


## kata-deploy会修改containerd配置
即便部署前已修改，配置会替换：
>[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
        runtime_type = "io.containerd.kata.v2"
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata.options]
           ConfigPath = "/opt/kata/share/defaults/kata-containers/configuration.toml"

https://github.com/containerd/containerd/issues/3073
https://github.com/containerd/containerd/issues/5006

## 规避方法
（无效）
设置KATA_CONF_FILE环境变量??
```
[root@localhost ~]# export KATA_CONF_FILE=/etc/kata-containers/configuration.toml
[root@localhost ~]# env | grep KATA_CONF_FILE
KATA_CONF_FILE=/etc/kata-containers/configuration.toml
[root@localhost ~]# ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hfftest dmesg
ctr: Cannot find usable config file (file /opt/kata/share/defaults/kata-containers/configuration-qemu.toml does not exist): not found
```
（无效）
```
[root@localhost ~]# ctr -n k8s.io run --env KATA_CONF_FILE=/etc/kata-containers/configuration.toml --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hfftest dmesg
ctr: Cannot find usable config file (file /opt/kata/share/defaults/kata-containers/configuration-qemu.toml does not exist): not found
```

（无效）
修改containerd配置：
```
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
  runtime_type = "io.containerd.kata.v2"
  privileged_without_host_devices = true
  pod_annotations = ["io.katacontainers.*"]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata.options]
    ConfigPath = "/etc/kata-containers/configuration.toml"
```

（有效）
```
[root@localhost ~]# cat /usr/local/bin/containerd-shim-kata-qemu-v2
#!/usr/bin/env bash
KATA_CONF_FILE=/opt/kata/share/defaults/kata-containers/configuration-qemu.toml /opt/kata/bin/containerd-shim-kata-v2 "$@"
改为：
[root@localhost ~]# cat /usr/local/bin/containerd-shim-kata-qemu-v2
#!/usr/bin/env bash
KATA_CONF_FILE=/etc/kata-containers/configuration.toml /opt/kata/bin/containerd-shim-kata-v2 "$@"
```
（有效）
```bash
ln -s /opt/kata/bin/containerd-shim-kata-v2 /usr/local/bin/containerd-shim-kata-v2 
```

# kata pod 大文件io性能测试导致pod重启，或节点卡死，节点异常问题
问题1： pod sanbox change，exec退出，pod重启次数加1
问题2： 节点卡死，节点notrady，相同配置runc测试无该问题

怀疑： 未开启SandboxCgroupOnly，导致sanbox overhead无限制占用主机资源，导致节点异常，这样的话就说明kata的资源开销还是很大的。

开启SandboxCgroupOnly后测试卡死，但是pod/节点未异常,进程结束。。。
所以，问题是，为什么没有错误信息




# 参考资料
https://github.com/kata-containers/kata-containers/blob/main/docs/Limitations.md

https://github.com/pulls?q=label%3Alimitation+org%3Akata-containers+is%3Aclosed

