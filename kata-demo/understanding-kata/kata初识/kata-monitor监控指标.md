[TOC]



https://github.com/kata-containers/kata-containers/blob/main/docs/design/kata-2-0-metrics.md

# 指标设计

Kata 实现 CRI 的 API，并支持 ContainerStats 和 ListContainerStats 接口以公开容器指标。用户可以使用这些界面来获取有关容器的基本指标。
但是与 runc 不同，Kata 是基于 VM 的运行时，并且具有不同的体系结构。

在 Kata 2.0 中，以下组件将能够提供有关系统的更多详细信息。

- containerd shim v2 (effectively kata-runtime)
- Hypervisor statistics
- Agent process
- Guest OS statistics

Kata 2.0 指标强烈依赖于 Prometheus。 Kata Containers 2.0 引入了一个名为 kata-monitor 的新 Kata 组件，该组件用于监视主机上的其他 Kata 组件。


containerd-shim-kata-v2 通过虚拟串口向 VM GuestOS（POD） 内的 kata-agent 请求监控数据，kata-agent 采集 GuestOS 内的容器监控数据并响应


# 指标列表
- Kata agent metrics
- Firecracker metrics
- Kata guest OS metrics
- Hypervisor metrics
- Kata monitor metrics
- Kata containerd shim v2 metrics

在kata vm中/proc/<pid>/io  stat status等的数据
kata_agent_io_stat代理进程 IO 统计
kata_agent_process_cpu_seconds_total 以秒为单位花费的总用户和系统 CPU 时间。
kata_agent_total_vm 

在kata vm中/proc/stat  diskstats meminfo等的数据
kata_guest_cpu_time
kata_guest_diskstat
kata_guest_load
kata_guest_meminfo
kata_guest_vm_stat

kata_hypervisor_io_stat 处理IO统计
kata_hypervisor_proc_stat  进程统计

kata_shim_pod_overhead_cpu CPU 资源的 Kata Pod 开销（百分比）
kata_shim_pod_overhead_memory_in_bytes Kata Pod 的内存资源开销（字节）




# promethues监控负载指标
- container_fs_writes_bytes_total 

- container_cpu_usage_seconds_total没有container字段
- 
```
sum(irate(container_cpu_usage_seconds_total{namespace=~"${allNamespace}",pod=~"^${loadNames}",container!=""}[3m]))by(pod)
```





# 指标的性能与开销
-  端到端（从 Prometheus 服务器到kata-monitor并kata-monitor写回响应）：20 毫秒（平均）
-  代理（从 shim 到agent的所有 RPC）：3 毫秒（平均）
-  Prometheus 默认scrape_interval为 1 分钟，但通常设置为 15 秒。较小scrape_interval会导致更多开销，因此用户应根据自己的监控需求进行设置。

	Prometheus 发出的一个指标获取请求的大小。当没有 gzip 压缩时，计算预期大小的公式是：  
9 + (144 - 9) *`number of kata sandboxes`
	Prometheus支持gzip压缩. 启用后，每个请求的响应大小会更小：  
2 + (10 - 2) *`number of kata sandboxes`

# endpoint
`kata-monitor`公开了以下endpoint·：
  *  `/metrics`              : 获取 Kata 沙箱指标。
  *  `/sandboxes`            : 列出主机上运行的所有 Kata 沙箱。
  *  `/agent-url`            : 获取 Kata 沙箱的代理 URL。
  *  `/debug/vars`           : Kata 运行时 shim 的内部数据。
  *  `/debug/pprof/`         : Kata 运行时 shim 的 Golang 分析数据：索引页。
  *  `/debug/pprof/cmdline` : Kata 运行时 shim 的 Golang 分析数据：`cmdline`endpoint。
  *  `/debug/pprof/profile` : Kata 运行时 shim 的 Golang 分析数据：`profile`endpoint（CPU 分析）。
  *  `/debug/pprof/symbol`   : Kata 运行时 shim 的 Golang 分析数据：`symbol`endpoint。
  *  `/debug/pprof/trace`    : Kata 运行时 shim 的 Golang 分析数据：`trace`endpoint。

`/agent-url`和所有`/debug/` * 都需要在查询字符串中指定`sandbox_id` 


# kata-monitor启动方式

1. kata节点运行kata-monitor守护进程
   curl 127.0.0.1:8090/sandboxes
   curl 127.0.0.1:8090/agent-url?sandboxes=df96b24bd49ec437c872c1a758edc084121d607ce1242ff5d2263a0e1b693343
2. daemonset(建议)


# `enable_pprof = true` 
configuration-qemu.toml


# daemonset部署
kubectl apply -f https://raw.githubusercontent.com/kata-containers/kata-containers/main/docs/how-to/data/kata-monitor-daemonset.yml
Once the daemonset is running, Prometheus should discover kata-monitor as a target. You can open http://<hostIP>:30909/service-discovery and find kubernetes-pods under the Service Discovery list

- 关于没有kata-monitor 镜像问题
https://github.com/kata-containers/kata-containers/issues/2421


- 镜像编译问题


