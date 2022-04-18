[TOC]

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

# curl
echo "hfftest 111" |curl --data-binary @- http://10.19.0.13:9091/metrics/job/schedulerStatus/instance/tztest


curl -X POST -g 'http://127.0.0.1:9090/api/v1/admin/tsdb/delete_series?match[]=scheduler_effective_dynamic_schedule_count' 

 curl -X POST -g 'http://10.19.0.13:9090/api/v1/admin/tsdb/delete_series?match[]={job="schedulerStatus"}'



# pushgateway
