[TOC]



https://github.com/kata-containers/kata-containers/blob/main/docs/design/kata-2-0-metrics.md


kata-monitor 进程运行在宿主机上，负责从各 Kata Containers 容器/VM中获取 metrics，并返回给 Prometheus。

默认情况下 kata-monitor 不需要指定参数，它会监听在本地的 8090 端口，这也是在 Prometheus 配置文件中 target 指定的端口号。如果要修改这个端口号，则需要注意两处要保持一致。


## kata-monitor启动方式
1. kata节点运行kata-monitor守护进程
```bash
[root@localhost ~]# cat /etc/systemd/system/kata-monitor.service
[Unit]
Description=kata monitor

[Service]
ExecStart=/opt/kata/bin/kata-monitor -listen-address 0.0.0.0:8090
Restart=always
StartLimitInterval=0
RestartSec=10

[Install]
WantedBy=multi-user.target
```

2. daemonset(？？没有镜像，手动编译不过)（TODO）
```bash
$ kubectl apply -f https://raw.githubusercontent.com/kata-containers/kata-containers/main/docs/how-to/data/kata-monitor-daemonset.yml
```

Once the daemonset is running, Prometheus should discover kata-monitor as a target. You can open http://<hostIP>:30909/service-discovery and find kubernetes-pods under the Service Discovery list

- 关于没有kata-monitor 镜像问题
https://github.com/kata-containers/kata-containers/issues/2421


## promethues增加scrape_configs
```bash
- job_name: 'kata'
    static_configs:
    - targets: ['<kata节点IP>:8090']
```
## 导入 Grafana dashborad

```bash
[root@localhost ~]# curl -XPOST -i <grafana节点IP>:3000/api/dashboards/import \
>     -u admin:admin \
>     -H "Content-Type: application/json" \
> -d "{\"dashboard\":$(curl -sL https://raw.githubusercontent.com/kata-containers/kata-containers/main/docs/how-to/data/dashboard.json )}"

HTTP/1.1 100 Continue

HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 16 May 2022 05:16:33 GMT
Content-Length: 253

{"pluginId":"","title":"Kata containers","imported":true,"importedUri":"db/kata-containers","importedUrl":"/d/75pdqURGk/kata-containers","slug":"","dashboardId":0,"folderId":0,"importedRevision":1,"revision":1,"description":"","path":"","removed":false}
```



# 监控指标

## Kata Containers 目前采集了下面几种类型的 metrics：

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




## promethues监控负载指标
- container_fs_writes_bytes_total 

- container_cpu_usage_seconds_total没有container字段
- 
```
sum(irate(container_cpu_usage_seconds_total{namespace=~"${allNamespace}",pod=~"^${loadNames}",container!=""}[3m]))by(pod)
```





## 指标的性能与开销
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

   curl 127.0.0.1:8090/sandboxes
   curl 127.0.0.1:8090/agent-url?sandboxes=df96b24bd49ec437c872c1a758edc084121d607ce1242ff5d2263a0e1b693343