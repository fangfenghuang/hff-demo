[TOC]

# 参考：
[host-cgroups.md](https://github.com/kata-containers/kata-containers/blob/main/docs/design/host-cgroups.md)
[vcpu-handling.md](https://github.com/kata-containers/kata-containers/blob/main/docs/design/vcpu-handling.md)

# 结论
- cgroupsPath的路径根据Qos类别不同
- kata_overhead的cpu和mem没有限制，能用到宿主机所有资源，如果不开启sandbox_cgroup_only会导致主机资源被抢占
- kubepod cgroup限制=limit+overhead，容器业务负载就是在这个限制之内使用，然后区分sandbox_cgroup_only是否开启决定kata vm(sandbox)开销是否统计到kubepod中，如果是则业务申请cpu mem资源就必须考虑kata sandbox开销
- 虚拟机大小=kata vm中或者说pod中看到的cpu mem（lscpu lsmem）=default+limit
- kubectl describe node | grep test-kata看到的资源申请=limit+overhead（同runc）
- kubectl top pod看到的是业务真正的负载开销？？
- pod overhead会影响调度，overhead只作用于sandbox开销，不能作用与业务负载（见验证3）
- 为了限制资源抢占问题，建议开启sandbox_cgroup_only=true,overhead建议不设置（待定），以下说明均在sandbox_cgroup_only=true前提下进行说明
- 业务负载+sandbox开销最大能使用的资源限制=limit+overhead（overhead部分只能是额外开销用）
- default_vcpus/memory，影响虚拟机的大小/启动时间，这个值官方不建议修改
- 如果container未设置limit，则kubepod cgroup无限制（但是业务负载只能使用default设置的1G2G???暂未找到在哪做的限制）；如果设置了limit，则这个default则只影响虚拟机大小，但不影响kubepod cgroup限制
- ~~~一个Pod的最小规格是1C 256M，当低于 256M 时，会重置为 2G。，支持的最大内存规格是256GB。如果用户分配的内存规格超过256GB，可能会出现未定义的错误，安全容器暂不支持超过256GB的大内存场景。~~~（未找到官方出处，应该只是openEuler内部做的限制）
- 如果不设置request，则request的值和limit默认相等




# kata cgroup说明
Kata Containers 在两层 cgroup 上运行。
- workload guest
- VMM host
（Kata Containers 在两层 cgroup 上运行。第一层在虚拟机中放置工作负载，而第二层在主机上运行 VMM和 关联线程。）


从 Kubernetes 角度来讲，Cgroup 指的是 Pod Cgroup，由 Kubelet 创建，限制的是 Pod 的资源; 从 Container 角度来讲，Cgroup 指的是 Container Cgroup，由对应的 runtime 创建，限制的是 container 的资源。但是为了可以获取到更准确的容器资源，Kubelet 会根据 Container Cgroup 去调整 Pod Cgroup。 在传统的 runtime 中，两者没有太大的区别。

而 Kata Containers 引入 VM 的概念，所以针对这种情况有两种处理方式：
- 启用 SandboxCgroupOnly（默认），Kubelet 在调整 Pod Cgroup 的大小时，会将 sandbox 的开销统计进去
（就是说，用户申请的资源limit+overhead是业务负载和kata虚拟机开销的最大使用限制，如果不设置limit，默认是default_vcpus/mem??）
- 禁用 SandboxCgroupOnly，sandbox 的开销和 Pod Cgroup 分开计算，独立存在
（就是说，用户申请的资源limit+overhead是业务负载自己的最大使用限制，kata虚拟机开销单独一个cgroup而且没有限制，如果不设置limit，业务使用默认是default_vcpus/mem）


## 两种cgrouppath:
/kubepods/burstable/pod287707dd-0fda-4fab-874d-e1b00e87390a/c2bfeccb7580707e7559c00f7ee9d46e98745cc2a0d267fbb41c68da41df784f
()
/kubepods/pod24356d87-3993-4e9a-8d3f-55207af763f9/kata_842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6
()

## pod overhead

## sandbox_cgroup_only
- 如果启用，运行时将在一个专用cgroup中添加所有kata进程。每个沙盒只创建一个cgroup，主机中的container cgroup将不会创建
- kata runtime调用者可以自由地限制或收集整个Kata沙盒的cgroup统计信息。
- sandbox cgroup路径是具有PodSandbox annotation(???)的容器的父cgroup。
- 如果没有container type annotation(??)，则sandbox cgroup将受到限制
- 意味着 pod cgroup 的大小适当，并考虑了 pod 开销（设置podoverhead考虑调度）
- Kata shim 将在sandbox_cgroup_onlypod 的专用 cgroup 下为每个 pod 创建一个子 cgroup
- kata shim、qemu 线程将与沙箱在同一个 cgroup 中。在这种情况下，它们将与 vcpu 线程共享同一组 cpu 内核

**=true好处：**
- 便于Pod资源统计
- 更好的主机资源隔离


**=fasle优点及缺点：**
- 不限制会消耗大量主机资源
- 在专用开销cgroup下运行所有非vcpu线程可以提供有关实际的开销准确值，可以设置开销cgroup大小（手动修改，无接口）


## default_vcpus/memory


## kata资源开销
- 半虚拟化 I/O 后端
- VMM 实例
- Kata shim 进程
- kata monitor开销（如果开启）

## kata组件进程：
- sandbox qemu-system-x86_64进程cgroup
- containerd-shim-kata-v2进程
- 两个virtiofsd进程

kata pod起来后，其中，kata shim进程和qemu进程占用了大量的内存开销，其中，沙箱里的内核和 agent 是直接分享了应用内存但并不是应用的一部分，也就是说，这一部分是用户可见的开销。而 VMM 本身在宿主机上的开销以及 shim-v2 占用的内存，虽然用户应用不可见，但同样是 Kata 带来的开销，影响基础设施的资源效率。


这些内存开销里，还包括了可共享开销和不可共享开销。所有宿主这边的代码段内存（只读内存），可以共享 page cache，所以是共享的，而在沙箱里的代码（只读）内存，在使用 DAX 或模版的情况下，也可以共享同一份。所以，可以共享的内存开销在节点全局范围内是唯一一份，所以在有多个Pod的情况下，这部分的常数开销相对而言是不太重要的，相反的，不可共享开销，尤其是堆内存开销，会随着 Pod 的增长而增长，的换句话说，堆内存（匿名页）是最需要被重视的开销。

综上，第一位需要被遏制的开销是沙箱内的用户可见的不可共享开销，尤其是 Agent 的开销，而 VMM 和 shim 的匿名页开销次之。


## 不支持Cgroups V2



# 一些需要知道的配置
```bash
# Default number of vCPUs per SB/VM:
# unspecified or 0                --> will be set to 1（不设置或0则默认是1）
# < 0                             --> will be set to the actual number of physical cores(小于0则设为实际物理核)
# > 0 <= number of physical cores --> will be set to the specified number
# > number of physical cores      --> will be set to the actual number of physical cores（大于物理核默认为物理核）
default_vcpus = 1

# Default memory size in MiB for SB/VM.
# If unspecified then it will be set 2048 MiB.（不设置默认是2Gi）
default_memory = 2048

# Default maximum number of vCPUs per SB/VM:
# unspecified or == 0             --> will be set to the actual number of physical cores or to the maximum number（不设置或0默认为最大物理核数）
#                                     of vCPUs supported by KVM if that number is exceeded
# > 0 <= number of physical cores --> will be set to the specified number
# > number of physical cores      --> will be set to the actual number of physical cores or to the maximum number
#                                     of vCPUs supported by KVM if that number is exceeded
# WARNING: Depending of the architecture, the maximum number of vCPUs supported by KVM is used when
# the actual number of physical cores is greater than it.
# WARNING: Be aware that this value impacts the virtual machine's memory footprint and CPU
# the hotplug functionality. For example, `default_maxvcpus = 240` specifies that until 240 vCPUs
# can be added to a SB/VM, but the memory footprint will be big. Another example, with
# `default_maxvcpus = 8` the memory footprint will be small, but 8 will be the maximum number of
# vCPUs supported by the SB/VM. In general, we recommend that you do not edit this variable,
# unless you know what are you doing.
# NOTICE: on arm platform with gicv2 interrupt controller, set it to 8.
default_maxvcpus = 0

```

## 通过注释修改配置：
```bash
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
 runtime_type = "io.containerd.kata.v2"
 privileged_without_host_devices = false
 shim_debug = true
 pod_annotations = ["io.katacontainers.*"] #  <-- look here
```

>  io.katacontainers.config.hypervisor.default_vcpus = 5

## 打开--feature-gates PodOverhead

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

## virtio-mem内存热插拔（仅支持QEMU）
[https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-use-virtio-mem-with-kata.md](https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-use-virtio-mem-with-kata.md)

```
$ sudo sed -i -e 's/^#enable_virtio_mem.*$/enable_virtio_mem = true/g' /etc/kata-containers/configuration.toml
```

# 一些验证
## 验证1： 查看cgroup大小
cpu.cfs_period_us是CFS算法的一个调度周期，一般它的值是100000us，即100ms。
cpu.cfs_quota_us表示在CFS算法中，在一个调度周期里该控制组被允许的运行时间，比如这个值为50000时，就是50ms。用这个值去除以调度周期cpu.cfs_period_us，即50ms/100ms=0.5，得到的值表示该控制组被允许使用的CPU最大配额是0.5个CPU。在我的系统里，这个值是-1，为默认值，表示不限制。

default_vcpu=1
default_memory=2Gi
overhead: "podFixed":{"cpu":"250m","memory":"160Mi"}
pod Qos:
        resources:
          requests:
            memory: "500Mi"
            cpu: "0.5"
          limits:
            memory: "3000Mi"
            cpu: "4"

### sandbox_cgroup_only=true
"cgroupsPath": "/kubepods/burstable/pod287707dd-0fda-4fab-874d-e1b00e87390a/c2bfeccb7580707e7559c00f7ee9d46e98745cc2a0d267fbb41c68da41df784f",

kata_overhead-mem: ~8G（不限）
kata_overhead-cpu: -1（不限）
kubepods_mem: 3313500160 (~3.08G)  (limit+overhead)
kubepods_cpu cfs_period_us: 100000  cpu.cfs_quota_us 425000  (4.25核) （limit+overhead）

lscpu: 5 (default+limit)
lsmem：
    Memory block size:       128M
    Memory block size:       128M
    Total online memory:       5G (default+limit)

kubectl describe node | grep test-kata： 750m (9%)     4250m (54%)(limit+overhead)  660Mi (10%)      3160Mi (48%)(limit+overhead)   4m54s
kubectl top pod： 1002m        1Mi
kubectl top node： 1306m        16%    5513Mi          84%

[root@localhost ~]# systemd-cgtop | grep kata
/kata_overhead                                                                 -      -     4.9M        -        -
/system.slice/kata-monitor.service                                             1      -    27.4M        -        -
[root@localhost ~]# systemd-cgtop | grep pod287707dd-0fda-4fab-874d-e1b00e87390a
/kubepods/burstable/pod287707dd-0fda-4fab-874d-e1b00e87390a                    -      -   193.4M        -        -

### sandbox_cgroup_only=false（结果同=true）
"cgroupsPath": "/kubepods/burstable/pod3a449bdd-90d5-41fd-a57a-9cfd102f9e44/87b5990ce61e4ff193124fc60e4a68c7061cf9299688c6c521cae29103bc3ef5",

kata_overhead-mem: ~8G（不限）
kata_overhead-cpu: -1（不限）

kubepods_mem: 3313500160 (~3.08G)  (limit+overhead)
kubepods_cpu cfs_period_us: 100000  cpu.cfs_quota_us 425000  (4.25核) （limit+overhead）


[root@localhost ~]# systemd-cgtop | grep kata
/kata_overhead                                                                 -      -     4.9M        -        -
/system.slice/kata-monitor.service                                             1      -    24.6M        -        -
[root@localhost ~]# systemd-cgtop | grep pod3a449bdd-90d5-41fd-a57a-9cfd102f9e44
/kubepods/burstable/pod3a449bdd-90d5-41fd-a57a-9cfd102f9e44                    -      -   151.7M        -        -

kubectl describe node | grep test-kata： 750m (9%)     4250m (54%)(limit+overhead)  660Mi (10%)      3160Mi (48%)(limit+overhead)   4m54s
kubectl top pod： 1002m        1Mi
kubectl top node：  1330m        17%    5507Mi          84%


### 查看命令：
```bash
[root@localhost ~]# systemd-cgtop | grep kata
/kata_overhead                                                                                                                -      -     4.9M        -        -
/kubepods/pod24356d87-3993-4e9a-8d3f-55207af763f9/kata_842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6       7      -   132.9M        -        -
/system.slice/kata-monitor.service                                                                                            1      -    23.2M        -        -

[root@localhost ~]# cat /sys/fs/cgroup/memory/kata_overhead/memory.limit_in_bytes
9223372036854771712
[root@localhost ~]#  cat /sys/fs/cgroup/cpu/kata_overhead/cpu.cfs_period_us
100000
[root@localhost ~]# cat /sys/fs/cgroup/cpu/kata_overhead/cpu.cfs_quota_us
-1

[root@localhost ~]# cat /sys/fs/cgroup/memory/kubepods/pod24356d87-3993-4e9a-8d3f-55207af763f9/memory.limit_in_bytes
1241513984
[root@localhost ~]# cat /sys/fs/cgroup/cpu/kubepods/pod24356d87-3993-4e9a-8d3f-55207af763f9/cpu.cfs_period_us
100000
[root@localhost ~]# cat /sys/fs/cgroup/cpu/kubepods/pod24356d87-3993-4e9a-8d3f-55207af763f9/cpu.cfs_quota_us
125000

[root@localhost ~]#  cat /sys/fs/cgroup/memory/system.slice/kata-monitor.service/memory.limit_in_bytes
9223372036854771712

[root@localhost ~]# kubectl top pod test-kata
NAME        CPU(cores)   MEMORY(bytes)
test-kata   0m           2Mi
[root@localhost ~]# kubectl top node
NAME                    CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
localhost.localdomain   355m         4%     5469Mi          83%

[root@localhost ~]# kubectl describe node | grep kata
  default                     test-kata                                                1250m (16%)   1250m (16%)  1184Mi (18%)     1184Mi (18%)   16d
  kube-system                 kata-deploy-q2cnv                                        0 (0%)        0 (0%)       0 (0%)           0 (0%)         21d

```

## 验证2：如果不设置request limit，kubepod cgroup的限制是怎样的

"cgroupsPath": "/kubepods/besteffort/pod09a392bc-1080-4044-8015-39f392e3862f/367368f8f5c44bc69b277febfa5d533063f096c65cd9b6729bbc7847a582f51f",


kata_overhead-mem: ~8G（不限）
kata_overhead-cpu: -1（不限）

kubepods_mem: 9223372036854771712  ~8G（不限）
kubepods_cpu cfs_period_us cfs_quota_us: 100000 -1 不限
kubectl top pod：  1000m        0Mi
kubectl top node：  1316m        16%    5517Mi          84%

[root@localhost hff]# systemd-cgtop | grep kata
/kata_overhead                                                                 -      -     4.9M        -        -
/system.slice/kata-monitor.service                                             1      -    25.1M        -        -

[root@localhost hff]# systemd-cgtop | grep pod09a392bc-1080-4044-8015-39f392e3862f
/kubepods/besteffort/pod09a392bc-1080-4044-8015-39f392e3862f                   -      -    95.8M        -        -

kubectl describe node | grep test-kata： 250m (3%)（overhead）     0 (0%)      160Mi (2%)（overhead）       0 (0%)         6m30s


kubepod不限制，但是pod确实只能用到limit的值

## 验证3： 验证overhead是否影响业务负载
设置overhead 2C 2G后，limit 3C3G, 打满cpu=6,mem=3Gi
[root@localhost hff]# kubectl get pod
NAME                              READY   STATUS        RESTARTS   AGE
test-kata-5687cdcb66-28lgn        0/1     OutOfmemory   0          17s
test-kata-5687cdcb66-2lmd5        0/1     OutOfmemory   0          11s
test-kata-5687cdcb66-2pbhm        0/1     OutOfmemory   0          23s
test-kata-5687cdcb66-575sb        0/1     OutOfmemory   0          19s
test-kata-5687cdcb66-5gvt2        0/1     OutOfmemory   0          5s
test-kata-5687cdcb66-62fbd        0/1     OutOfmemory   0          33s

## 验证4：验证如果不设置request limit，cpu mem的使用上线
设置overhead 2C 2G后，limit不设置， 打满cpu=6,mem=1Gi

[root@localhost hff]# kubectl get pod  -w
NAME                              READY   STATUS        RESTARTS   AGE
test-kata-55867ffb58-26ndp        0/1     OutOfmemory   0          4s
test-kata-55867ffb58-2blnx        0/1     OutOfmemory   0          2s
test-kata-55867ffb58-2d858        0/1     OutOfmemory   0          1s
test-kata-55867ffb58-2qwzf        0/1     Pending       0          0s
test-kata-55867ffb58-4j4hl        0/1     OutOfmemory   0          1s
test-kata-55867ffb58-8zzbc        0/1     OutOfmemory   0          2s
test-kata-55867ffb58-crflb        0/1     OutOfmemory   0          4s

## 验证5： 验证cpu mem小于1C 256Mi的情况（no overhead）
        resources:
          requests:
            memory: "1Mi"
            cpu: "0.1"
          limits:
            memory: "100Mi"
            cpu: "550m"
        args:
        - -cpus
        - "2"
        - -mem-total
        - 100Mi
        - -mem-alloc-size
        - 10Mi
        - -mem-alloc-sleep
        - 1s

没有OOM，但是pod不断重启，并没有被重置成2G
[root@localhost hff]# kubectl get pod -w
NAME                              READY   STATUS    RESTARTS   AGE
kubefilebrowser-896974bcc-z6scd   1/1     Running   3          31d
test-kata-6878f4fd9f-76c5k        1/1     Running   2          61s
test-runc-79d8bdc4cb-v6j28        1/1     Running   0          75m
test-kata-6878f4fd9f-76c5k        0/1     Error     2          72s
test-kata-6878f4fd9f-76c5k        0/1     CrashLoopBackOff   2          75s

[root@localhost hff]# kubectl logs test-kata-6878f4fd9f-76c5k
I0527 07:16:06.932758       1 main.go:26] Allocating "200Mi" memory, in "10Mi" chunks, with a 1s sleep between allocations
I0527 07:16:06.933058       1 main.go:39] Spawning a thread to consume CPU
I0527 07:16:06.933082       1 main.go:39] Spawning a thread to consume CPU

## 验证6： 极小资源占用测试
limit550m/200Mi，压2C100Mi，pod会不断重启
        resources:
          requests:
            memory: "1Mi"
            cpu: "0.1"
          limits:
            memory: "200Mi"
            cpu: "550m"
        args:
        - -cpus
        - "2"
        - -mem-total
        - 100Mi
        - -mem-alloc-size
        - 10Mi
        - -mem-alloc-sleep
        - 1s
Events:
  Type     Reason          Age                   From                            Message
  ----     ------          ----                  ----                            -------
  Normal   Scheduled       <unknown>             default-scheduler               Successfully assigned default/test-kata-68b77bc9f8-zlj5p to localhost.localdomain
  Normal   Created         36m (x4 over 38m)     kubelet, localhost.localdomain  Created container cpu-stress-kata
  Normal   Started         36m (x4 over 38m)     kubelet, localhost.localdomain  Started container cpu-stress-kata
  Normal   SandboxChanged  35m (x4 over 37m)     kubelet, localhost.localdomain  Pod sandbox changed, it will be killed and re-created.
  Normal   Pulled          7m53s (x11 over 38m)  kubelet, localhost.localdomain  Container image "vish/stress" already present on machine
  Warning  BackOff         3m5s (x159 over 37m)  kubelet, localhost.localdomain  Back-off restarting failed container
[root@localhost hff]#
[root@localhost hff]# kubectl logs test-kata-68b77bc9f8-zlj5p
I0527 08:13:52.926368       1 main.go:26] Allocating "100Mi" memory, in "10Mi" chunks, with a 1s sleep between allocations
I0527 08:13:52.926544       1 main.go:39] Spawning a thread to consume CPU
I0527 08:13:52.926570       1 main.go:39] Spawning a thread to consume CPU

## 验证7： 查看kata vm中的cgroup
不设置limit，overhead
bash-5.1# more /sys/fs/cgroup/memory/
cgroup.clone_children               memory.max_usage_in_bytes
cgroup.event_control                memory.memsw.failcnt
cgroup.procs                        memory.memsw.limit_in_bytes
cgroup.sane_behavior                memory.memsw.max_usage_in_bytes
kubepods/                           memory.memsw.usage_in_bytes
memory.failcnt                      memory.move_charge_at_immigrate
memory.force_empty                  memory.oom_control
memory.kmem.failcnt                 memory.pressure_level
memory.kmem.limit_in_bytes          memory.soft_limit_in_bytes
memory.kmem.max_usage_in_bytes      memory.stat
memory.kmem.slabinfo                memory.swappiness
memory.kmem.tcp.failcnt             memory.usage_in_bytes
memory.kmem.tcp.limit_in_bytes      memory.use_hierarchy
memory.kmem.tcp.max_usage_in_bytes  notify_on_release
memory.kmem.tcp.usage_in_bytes      release_agent
memory.kmem.usage_in_bytes          system.slice/
memory.limit_in_bytes               tasks
bash-5.1# more /sys/fs/cgroup/memory/kubepods/
besteffort/                         memory.max_usage_in_bytes
cgroup.clone_children               memory.memsw.failcnt
cgroup.event_control                memory.memsw.limit_in_bytes
cgroup.procs                        memory.memsw.max_usage_in_bytes
memory.failcnt                      memory.memsw.usage_in_bytes
memory.force_empty                  memory.move_charge_at_immigrate
memory.kmem.failcnt                 memory.oom_control
memory.kmem.limit_in_bytes          memory.pressure_level
memory.kmem.max_usage_in_bytes      memory.soft_limit_in_bytes
memory.kmem.slabinfo                memory.stat
memory.kmem.tcp.failcnt             memory.swappiness
memory.kmem.tcp.limit_in_bytes      memory.usage_in_bytes
memory.kmem.tcp.max_usage_in_bytes  memory.use_hierarchy
memory.kmem.tcp.usage_in_bytes      notify_on_release
memory.kmem.usage_in_bytes          tasks
memory.limit_in_bytes
bash-5.1# more /sys/fs/cgroup/cpu/cpu.cfs_quota_us
-1
bash-5.1# more /sys/fs/cgroup/cpu/kubepods/besteffort/cpu.cfs_quota_us
-1
bash-5.1# more /sys/fs/cgroup/cpu/kubepods/besteffort/pode0eb5953-c643-4429-8812-b8e36d75e3d9/f3cffcf893000b0c5f6105f1adb94e3be7fa62ebecd801a29b9a77a472d62c86/cfs_quota_us
-1
bash-5.1# more /sys/fs/cgroup/memory/memory.limit_in_bytes
9223372036854771712
bash-5.1#  more /sys/fs/cgroup/memory/kubepods/memory.limit_in_bytes
9223372036854771712
bash-5.1#  more /sys/fs/cgroup/memory/kubepods/besteffort/memory.limit_in_bytes
9223372036854771712
bash-5.1#  more /sys/fs/cgroup/memory/kubepods/besteffort/pode0eb5953-c643-4429-8812-b8e36d75e3d9/.limit_in_bytes
9223372036854771712


# yaml
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-kata
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-kata
  template:
    metadata:
      labels:
        app: test-kata
    spec:
      runtimeClassName: kata-containers
      containers:
      - image: vish/stress
        imagePullPolicy: IfNotPresent
        name: cpu-stress-kata
        resources:
          requests:
            memory: "1Mi"
            cpu: "0.1"
          limits:
            memory: "161Mi"
            cpu: "350m"
        args:
        - -cpus
        - "2"
        - -mem-total
        - 100Mi
        - -mem-alloc-size
        - 10Mi
        - -mem-alloc-sleep
        - 1s
```
