[TOC]

# 组件
- Prometheus Server: 用于收集和存储时间序列数据。
- Client Library: 客户端库，为需要监控的服务生成相应的 metrics 并暴露给 Prometheus server。当Prometheus server 来 pull 时，直接返回实时状态的 metrics。对于机器层面的 metrices，需要使用 node exporter。
- Push Gateway: 主要用于短期的 jobs。
- Exporters: 用于暴露已有的第三方服务的 metrics 给 Prometheus。
- Alertmanager: 从 Prometheus server 端接收到 alerts 后，会进行去除重复数据，分组，并路由到对收的接受方式，发出报警。常见的接收方式有：电子邮件，pagerduty，OpsGenie, webhook 等。



# curl
echo "hfftest 111" |curl --data-binary @- http://10.19.0.13:9091/metrics/job/schedulerStatus/instance/tztest

curl -X POST -g 'http://127.0.0.1:9090/api/v1/admin/tsdb/delete_series?match[]=scheduler_effective_dynamic_schedule_count' 

 curl -X POST -g 'http://10.19.0.13:9090/api/v1/admin/tsdb/delete_series?match[]={job="schedulerStatus"}'

## 热加载
```shell
curl -XPOST <prometheus-url>/-/reload
```

# pushgateway


# Grafana


# 指标
## 指标类型
Counter（计数器）对数据只增不减
Gauage（仪表盘）可增可减
Histogram（直方图）,Summary（摘要）提供更多的统计信

## 计算表达式：
Prometheus为不同的数据提供了非常多的计算函数，其中有个小技巧就是遇到counter数据类型，在做任何操作之前，先套上一个rate()或者increase()函数

>100-avg(irate(node_cpu_seconds_total{mode='idle'}[5m])) by (node_name)*100

## 指标命名规范
一个指标名称：
	• 必须符合有效字符的数据模型。
	• 应该具有与指标所属域相关的（单个词汇）应用程序前缀。前缀有时被客户端库称为命名空间。对于特定于应用程序的指标，前缀通常是应用程序名称本身。然而，有时候指标更通用，比如客户端库导出的标准化指标。例如：
		○ prometheus_notifications_total （针对Prometheus 服务器）
		○ process_cpu_seconds_total （由客户端库导出）
		○ http_request_duration_seconds （用于所有HTTP请求）
	• 必须有一个单一的单位（即，不要把秒与毫秒，或秒与字节混用）。
	• 应该使用基本单位（如秒、字节、米——而不是毫秒、兆字节、公里）。参见下面的基本单位列表。
	• 应以复数形式用后缀来描述单位。请注意，累计计数以total作为后缀，附加在单位之后。
		○ http_request_duration_seconds
		○ node_memory_usage_bytes
		○ http_requests_total （用于无单位的累计计数）
		○ process_cpu_seconds_total （用于有单位的累计计数）
		○ foobar_build_info （用于提供关于正在运行的二进制文件的元数据的伪指标）
	• 应该在所有的标签维度中表示相同的监控逻辑。
		○ 请求持久时长
		○ 传输的数据字节数
		○ 瞬时资源使用百分比
