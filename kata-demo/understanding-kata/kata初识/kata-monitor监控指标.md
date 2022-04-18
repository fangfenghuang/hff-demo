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


# 指标的性能与开销
-  端到端（从 Prometheus 服务器到kata-monitor并kata-monitor写回响应）：20 毫秒（平均）
-  代理（从 shim 到agent的所有 RPC）：3 毫秒（平均）
-  Prometheus 默认scrape_interval为 1 分钟，但通常设置为 15 秒。较小scrape_interval会导致更多开销，因此用户应根据自己的监控需求进行设置。

	Prometheus 发出的一个指标获取请求的大小。当没有 gzip 压缩时，计算预期大小的公式是：  
9 + (144 - 9) *`number of kata sandboxes`
	Prometheus支持gzip压缩. 启用后，每个请求的响应大小会更小：  
2 + (10 - 2) *`number of kata sandboxes`


# kata-monitor启动方式

每个节点只运行一个kata-monitor进程





