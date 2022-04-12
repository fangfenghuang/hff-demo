[TOC]

https://github.com/kata-containers/kata-containers/blob/main/docs/design/kata-2-0-metrics.md

containerd-shim-kata-v2 通过虚拟串口向 VM GuestOS（POD） 内的 kata-agent 请求监控数据，kata-agent 采集 GuestOS 内的容器监控数据并响应


# 指标列表
### Kata agent metrics
### Firecracker metrics
### Kata guest OS metrics
### Hypervisor metrics
### Kata monitor metrics
### Kata containerd shim v2 metrics


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





