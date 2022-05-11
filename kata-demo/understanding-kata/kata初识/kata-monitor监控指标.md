[TOC]



https://github.com/kata-containers/kata-containers/blob/main/docs/design/kata-2-0-metrics.md


kata-monitor 进程运行在宿主机上，负责从各 Kata Containers 容器/VM中获取 metrics，并返回给 Prometheus。

默认情况下 kata-monitor 不需要指定参数，它会监听在本地的 8090 端口，这也是在 Prometheus 配置文件中 target 指定的端口号。如果要修改这个端口号，则需要注意两处要保持一致。


# Kata Containers 目前采集了下面几种类型的 metrics：

Kata agent metrics：agent 进程的 metrics
Kata guest OS metrics：VM 中的 guest metrics
Hypervisor metrics：hypervisor 进程的 metrics（如果 hypervisor 本身提供了 metrics 接口，比如 firecracker，也会采集到 Kata Containers 的 metrics）
Kata monitor metrics：kata-monitor 进程的 metrics
Kata containerd shim v2 metrics：shimv2 进程的 metrics




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



# kata指标实例
```bash
[root@localhost hff]# kata-runtime metrics 842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6 | grep kata_shim_pod_overhead_cpu
# HELP kata_shim_pod_overhead_cpu Kata Pod overhead for CPU resources(percent).
# TYPE kata_shim_pod_overhead_cpu gauge
kata_shim_pod_overhead_cpu 1.016532023719192(????)

[root@localhost hff]# kata-runtime metrics 842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6 | grep kata_shim_pod_overhead_memory_in_bytes
# HELP kata_shim_pod_overhead_memory_in_bytes Kata Pod overhead for memory resources(bytes).
# TYPE kata_shim_pod_overhead_memory_in_bytes gauge
kata_shim_pod_overhead_memory_in_bytes 1.17354496e+08（对应kubepods）


[root@localhost hff]# kata-runtime metrics 842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6 | grep kata_agent_total_vm
# HELP kata_agent_total_vm Agent process total VM size
# TYPE kata_agent_total_vm gauge
kata_agent_total_vm 2.0017152e+07


[root@localhost hff]# kata-runtime metrics 842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6 | grep kata_guest_load
# HELP kata_guest_load Guest system load.
# TYPE kata_guest_load gauge
kata_guest_load{item="load1"} 0
kata_guest_load{item="load15"} 0
kata_guest_load{item="load5"} 0

[root@localhost hff]# systemd-cgtop | grep pod24356d87-3993-4e9a-8d3f-55207af763f9
/kubepods/pod24356d87-3993-4e9a-8d3f-55207af763f9                                                                             -      -   115.3M        -        -
/kubepods/pod24356d87-3993-4e9a-8d3f-55207af763f9/kata_842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6       7      -   115.3M        -        -
/kata_overhead                                                                                                                -      -    47.9M        -        -

[root@localhost hff]# kubectl top pod
NAME                              CPU(cores)   MEMORY(bytes)
test-kata                         0m           0Mi

## 开启kata-monitor
[root@localhost ~]# curl 127.0.0.1:8090/sandboxes
842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6
[root@localhost ~]# kata-runtime metrics 842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6 | grep kata_shim_pod_overhead_cpu
# HELP kata_shim_pod_overhead_cpu Kata Pod overhead for CPU resources(percent).
# TYPE kata_shim_pod_overhead_cpu gauge
kata_shim_pod_overhead_cpu 0.3821455739873718(降下来了？？？)
[root@localhost ~]# kata-runtime metrics 842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6 | grep kata_shim_pod_overhead_memory_in_bytes
# HELP kata_shim_pod_overhead_memory_in_bytes Kata Pod overhead for memory resources(bytes).
# TYPE kata_shim_pod_overhead_memory_in_bytes gauge
kata_shim_pod_overhead_memory_in_bytes 1.19877632e+08
[root@localhost ~]# kata-runtime metrics 842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6 | grep kata_agent_total_vm
# HELP kata_agent_total_vm Agent process total VM size
# TYPE kata_agent_total_vm gauge
kata_agent_total_vm 2.0033536e+07
[root@localhost ~]#  kata-runtime metrics 842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6 | grep kata_guest_load
# HELP kata_guest_load Guest system load.
# TYPE kata_guest_load gauge
kata_guest_load{item="load1"} 0
kata_guest_load{item="load15"} 0
kata_guest_load{item="load5"} 0
[root@localhost ~]# systemd-cgtop | grep pod24356d87-3993-4e9a-8d3f-55207af763f9
/kubepods/pod24356d87-3993-4e9a-8d3f-55207af763f9                                                                             -      -   117.2M        -        -
/kubepods/pod24356d87-3993-4e9a-8d3f-55207af763f9/kata_842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6       7      -   117.2M        -        -


[root@localhost ~]# curl 127.0.0.1:8090/metrics | grep kata_shim_pod_overhead_cpu
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  407k    0  407k    0     0  31.5M      0 --:--:-- --:--:-- --:--:-- 33.1M
# HELP kata_shim_pod_overhead_cpu Kata Pod overhead for CPU resources(percent).
# TYPE kata_shim_pod_overhead_cpu gauge
kata_shim_pod_overhead_cpu{sandbox_id="842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6",cri_uid="24356d87-3993-4e9a-8d3f-55207af763f9",cri_name="test-kata",cri_namespace="default"} 0.7711305188284795
[root@localhost ~]# curl 127.0.0.1:8090/metrics | grep kata_shim_pod_overhead_memory_in_bytes
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0# HELP kata_shim_pod_overhead_memory_in_bytes Kata Pod overhead for memory resources(bytes).
# TYPE kata_shim_pod_overhead_memory_in_bytes gauge
kata_shim_pod_overhead_memory_in_bytes{sandbox_id="842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6",cri_uid="24356d87-3993-4e9a-8d3f-55207af763f9",cri_name="test-kata",cri_namespace="default"} 1.19918592e+08
100  407k    0  407k    0     0  28.6M      0 --:--:-- --:--:-- --:--:-- 30.5M
[root@localhost ~]# curl 127.0.0.1:8090/metrics | grep kata_agent_total_vm
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0# HELP kata_agent_total_vm Agent process total VM size
# TYPE kata_agent_total_vm gauge
kata_agent_total_vm{sandbox_id="842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6",cri_uid="24356d87-3993-4e9a-8d3f-55207af763f9",cri_name="test-kata",cri_namespace="default"} 2.0037632e+07
100  407k    0  407k    0     0  30.4M      0 --:--:-- --:--:-- --:--:-- 33.1M
[root@localhost ~]# curl 127.0.0.1:8090/metrics | grep kata_guest_load
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0# HELP kata_guest_load Guest system load.
# TYPE kata_guest_load gauge
kata_guest_load{item="load1",sandbox_id="842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6",cri_uid="24356d87-3993-4e9a-8d3f-55207af763f9",cri_name="test-kata",cri_namespace="default"} 0
kata_guest_load{item="load15",sandbox_id="842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6",cri_uid="24356d87-3993-4e9a-8d3f-55207af763f9",cri_name="test-kata",cri_namespace="default"} 0
kata_guest_load{item="load5",sandbox_id="842463ada9994438f6c663b082bbf5735c64309f77a0b57838b8a16f347433a6",cri_uid="24356d87-3993-4e9a-8d3f-55207af763f9",cri_name="test-kata",cri_namespace="default"} 0
100  407k    0  407k    0     0  31.3M      0 --:--:-- --:--:-- --:--:-- 33.1M
```


不起kata monitor可以看到指标，但是 curl 127.0.0.1:8090/sandboxes不能查看


# 使用 Prometheus + Grafana 监控 Kata Containers
http://liubin.org/kata-dev-book/src/kata-prom-grafana.html

编辑 Prometheus 配置文件
```bash
# 在 scrape_configs 部分的最后，加入下面的 target：
  - job_name: 'kata'
    static_configs:
    - targets: ['localhost:8090']

curl -XPOST <prometheus-url>/-/reload  (curl -XPOST http://10.240.229.101:32090/-/reload)
```

http://10.208.11.110:3000/
http://10.208.11.110:9090/




kata-monitor 启动后，就可以在 Prometheus targets 页面（ http://<your_server>:9090/targets ）看到我们的 target 的状态了（UP还是DOWN）。