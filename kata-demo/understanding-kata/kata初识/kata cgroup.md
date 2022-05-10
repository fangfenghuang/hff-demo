[TOC]

参考：[host-cgroups](https://github.com/kata-containers/kata-containers/blob/main/docs/design/host-cgroups.md)

Kata Containers 在两层 cgroup 上运行。
- workload guest
- VMM host

对于 Kubernetes，pod cgroup(sandbox沙箱)由Kubelet创建，而container cgroup由运行时处理。


# kata cgroup

# 资源开销
- 半虚拟化 I/O 后端
- VMM 实例
- Kata shim 进程
- sandbox资源开销
- kata monitor开销（如果开启）


区分是否并入pod cgroup的情况，如果不并入则单独一个sandbox cgroup，否则开销将记入到pod cgroup中影响业务性能


kata pod起来后，其中，kata shim进程和qemu进程占用了大量的内存开销，其中，沙箱里的内核和 agent 是直接分享了应用内存但并不是应用的一部分，也就是说，这一部分是用户可见的开销。而 VMM 本身在宿主机上的开销以及 shim-v2 占用的内存，虽然用户应用不可见，但同样是 Kata 带来的开销，影响基础设施的资源效率。


这些内存开销里，还包括了可共享开销和不可共享开销。所有宿主这边的代码段内存（只读内存），可以共享 page cache，所以是共享的，而在沙箱里的代码（只读）内存，在使用 DAX 或模版的情况下，也可以共享同一份。所以，可以共享的内存开销在节点全局范围内是唯一一份，所以在有多个Pod的情况下，这部分的常数开销相对而言是不太重要的，相反的，不可共享开销，尤其是堆内存开销，会随着 Pod 的增长而增长，的换句话说，堆内存（匿名页）是最需要被重视的开销。

综上，第一位需要被遏制的开销是沙箱内的用户可见的不可共享开销，尤其是 Agent 的开销，而 VMM 和 shim 的匿名页开销次之。



## sandbox_cgroup_only=true（开销记入pod cgroup）
- 意味着 pod cgroup 的大小适当，并考虑了 pod 开销（设置podoverhead考虑调度）
- Kata shim 将在sandbox_cgroup_onlypod 的专用 cgroup 下为每个 pod 创建一个子 cgroup


**好处：**
- 便于Pod资源统计
- 更好的主机资源隔离


例如：
```bash
# pod cgroup（sandbox开销记入pod cgroup）
/sys/fs/cgroup/memory/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca/
```

```bash
# pod cgroup
[root@localhost hff]# crictl inspect 79bccf2621042 | grep cgroupsPath
"cgroupsPath": "/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/79bccf2621042bd6f202187f9f24d544a04d5b841c83c2a929ffe1b17fdee081",
```

```bash
# sandbox qemu-system-x86_64进程cgroup
[root@localhost ~]# cat /proc/32597/cgroup
11:pids:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
10:blkio:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
9:hugetlb:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
8:perf_event:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
7:devices:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
6:cpuset:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
5:cpuacct,cpu:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
4:freezer:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
3:net_prio,net_cls:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
2:memory:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
1:name=systemd:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
## 当前进程已经被加到kata pod 专用cgroup中中

# containerd-shim-kata-v2进程
[root@localhost ~]# cat /proc/32583/cgroup
11:pids:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
10:blkio:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
9:hugetlb:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
8:perf_event:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
7:devices:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
6:cpuset:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
5:cpuacct,cpu:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
4:freezer:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
3:net_prio,net_cls:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
2:memory:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
1:name=systemd:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca

# 两个virtiofsd进程
[root@localhost ~]# cat /proc/32591/cgroup
11:pids:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
10:blkio:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
9:hugetlb:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
8:perf_event:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
7:devices:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
6:cpuset:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
5:cpuacct,cpu:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
4:freezer:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
3:net_prio,net_cls:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
2:memory:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
1:name=systemd:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
[root@localhost ~]# cat /proc/32604/cgroup
11:pids:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
10:blkio:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
9:hugetlb:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
8:perf_event:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
7:devices:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
6:cpuset:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
5:cpuacct,cpu:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
4:freezer:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
3:net_prio,net_cls:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
2:memory:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
1:name=systemd:/kubepods/pod64cf4d8f-9420-438b-9376-6d6e0cc56a16/kata_e7ff2d95cf2e933d5d876c54bbb746886d0c2099a04bd8941be17669907bb0ca
```


## sandbox_cgroup_only=false(默认)(开销单独一个kata_overhead开销cgroup，不限制占用主机资源)
Kata Containers 创建了一个不受约束的开销cgroup（kata_overhead），除了虚拟 CPU 线程之外的任何进程都移动到该cgroup下一个pod对应一个子cgroup。
```bash
[root@localhost ~]# find / -name kata_overhead
/sys/fs/cgroup/pids/kata_overhead
/sys/fs/cgroup/blkio/kata_overhead
/sys/fs/cgroup/hugetlb/kata_overhead
/sys/fs/cgroup/perf_event/kata_overhead
/sys/fs/cgroup/devices/kata_overhead
/sys/fs/cgroup/cpuset/kata_overhead
/sys/fs/cgroup/cpu,cpuacct/kata_overhead
/sys/fs/cgroup/freezer/kata_overhead
/sys/fs/cgroup/net_cls,net_prio/kata_overhead
/sys/fs/cgroup/memory/kata_overhead
/sys/fs/cgroup/systemd/kata_overhead
```

**优点及缺点：**
- 不限制会消耗大量主机资源
- 在专用开销cgroup下运行所有非vcpu线程可以提供有关实际的开销准确值，可以设置开销cgroup大小（手动修改，无接口）



例如：
```bash
# pod cgroup
/sys/fs/cgroup/memory/kubepods/pod82d8457f-63e8-4b2c-9f38-d157a5873bed/kata_ac5be1ec5262ac9cbad8e0cb6fdda98f53128946eb63545ec246ee6a565b8b29
# 开销cgroup
/sys/fs/cgroup/memory/kata_overhead/ac5be1ec5262ac9cbad8e0cb6fdda98f53128946eb63545ec246ee6a565b8b29
```


```bash
# pod cgroup
[root@localhost ~]# crictl inspect 28a1e515a7168 | grep cgroupsPath
"cgroupsPath": "/kubepods/pod82d8457f-63e8-4b2c-9f38-d157a5873bed/28a1e515a7168b097f218135eba44852c0619f497e943155bc54a5812322a791",
```

```bash
# sandbox qemu-system-x86_64进程cgroup
[root@localhost hff]#  cat /proc/1973/cgroup
11:pids:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
10:blkio:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
9:hugetlb:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
8:perf_event:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
7:devices:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
6:cpuset:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
5:cpuacct,cpu:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
4:freezer:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
3:net_prio,net_cls:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
2:memory:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
1:name=systemd:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
## 当前进程已经被加到systemd:/kata_overhead/中

# containerd-shim-kata-v2进程
[root@localhost hff]# cat /proc/1950/cgroup
11:pids:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
10:blkio:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
9:hugetlb:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
8:perf_event:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
7:devices:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
6:cpuset:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
5:cpuacct,cpu:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
4:freezer:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
3:net_prio,net_cls:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
2:memory:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
1:name=systemd:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
## 当前进程已经被加到systemd:/kata_overhead/中

# 两个virtiofsd进程
[root@localhost hff]# cat /proc/1967/cgroup
11:pids:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
10:blkio:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
9:hugetlb:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
8:perf_event:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
7:devices:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
6:cpuset:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
5:cpuacct,cpu:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
4:freezer:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
3:net_prio,net_cls:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
2:memory:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
1:name=systemd:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
## 当前进程已经被加到systemd:/kata_overhead/中
[root@localhost hff]# cat /proc/1978/cgroup
11:pids:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
10:blkio:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
9:hugetlb:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
8:perf_event:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
7:devices:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
6:cpuset:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
5:cpuacct,cpu:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
4:freezer:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
3:net_prio,net_cls:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
2:memory:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
1:name=systemd:/kata_overhead/9d381fb97f3375549aea9174941bf86b3e4d46e7cc41fa6f491b0c5f8c9124de
## 当前进程已经被加到systemd:/kata_overhead/中
```

# pod overhead
只影响调度？？？？



# 不支持Cgroups V2


# runc pod
```bash
[root@localhost ~]# cat /proc/17955/cgroup
11:pids:/system.slice/containerd.service
10:blkio:/system.slice/containerd.service
9:hugetlb:/
8:perf_event:/
7:devices:/system.slice/containerd.service
6:cpuset:/
5:cpuacct,cpu:/system.slice/containerd.service
4:freezer:/
3:net_prio,net_cls:/
2:memory:/system.slice/containerd.service
1:name=systemd:/system.slice/containerd.service
```