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




# ctr依赖/configuration-qemu.toml路径问题
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

# kata deploy重启/节点重启会导致配置还原问题（必现）
```bash
cat /opt/kata/share/defaults/kata-containers/configuration-qemu.toml  | grep sandbox_cgroup_only
sandbox_cgroup_only=true

## 重启后
[root@localhost ~]# cat /opt/kata/share/defaults/kata-containers/configuration-qemu.toml  | grep sandbox_cgroup_only
sandbox_cgroup_only=false

```
# 资源开销问题
## 容器内fio测试(压内存)会导致host上对应qemu进程oom，或节点卡死，节点异常问题
问题1： pod sanbox change，exec退出，pod重启次数加1
问题2： 节点卡死，节点notrady，相同配置runc测试无该问题

怀疑： 未开启SandboxCgroupOnly，导致sanbox overhead无限制占用主机资源，导致节点异常，这样的话就说明kata的资源开销还是很大的。

开启SandboxCgroupOnly后测试卡死，但是pod/节点未异常,进程结束。。。
所以，问题是，为什么没有错误信息

## Overhead的CPU和内存占用应该纳入已分配配额？？？

## 容器使用内存不会自动release？？

# 存储性能问题
## 未enable DAX, fio测试结果较差？？

# 网络问题
## qemu 不能直接使用 veth-pair, 导致vm+container的网络拓扑比较复杂且容易有性能问题, kubevirt同样的问题


# 参考资料
https://github.com/kata-containers/kata-containers/blob/main/docs/Limitations.md
https://github.com/pulls?q=label%3Alimitation+org%3Akata-containers+is%3Aclosed
https://github.com/chestack/k8s-blog/blob/master/kata-container/kata-container.md

