 [TOC]


# 当前容器云网络模型及网络开销说明
k8s的最小调度单位为pod，pod网络通信的实现依赖于第三方插件；容器云使用calico纯三层虚拟网络方案，可以避免与其他二层方案相关的数据包封装的操作，中间没有任何的NAT，没有任何的overlay，几乎能达到主机性能。
kata增加了一层tcfilter（默认）打通虚拟机与容器之间的网络


**测试配置说明:**
- 资源限制request/limit 1C2G
- kata设置debug_console_enabled=true（虚拟机开销占用业务开销）
- kata设置debug_console_enabled=false（虚拟机开销不限制）

# qperf
qperf和iperf/netperf一样可以评测两个节点之间的带宽和延时。可以在测试tcp/ip协议和RDMA传输。不过相比netperf和iperf，支持RDMA是qperf工具的独有特性。
- 循环测试1bytes-64KiB的带宽和延迟
```bash
qperf <ip> -oo msg_size:1:256K:*64 -vu tcp_bw tcp_lat
```

**注意:**
qperf测试service有问题：
1. 对于runc容器，需要修改qperf服务监听端口，否则跨节点无法测试
2. 对于kata容器，本节点和跨节点都无法测试

## 测试结果:
|               |msg_size   |宿主机服务端|runc容器服务端|kata容器服务端(true)|
|---------------|-----------|-----------|-------------|-------------|
| 从主机到服务   |1bytes     |tcp_lat:13 us<br>tcp_bw:1.13 MB/sec |tcp_lat:12.5 us<br>tcp_bw:1.14 MB/sec | tcp_lat:20.2 us<br>tcp_bw:1.05 MB/sec
|               |64bytes    |tcp_lat:12.9 us<br>tcp_bw:77 MB/sec |tcp_lat:12.6 us<br>tcp_bw:42.6 MB/sec | tcp_lat:19.3 us<br>tcp_bw:45.6 MB/sec
|               |4KiB       |tcp_lat:14.1 us<br>tcp_bw:1.71 GB/sec |tcp_lat:20.2 us<br>tcp_bw:335 MB/sec | tcp_lat:25.7 us<br>tcp_bw:355 MB/sec
|               |256KiB     |tcp_lat:68.8 us<br>tcp_bw:3.48 GB/sec |tcp_lat:92.1 us<br>tcp_bw:3.39 GB/sec| tcp_lat:167 us<br>tcp_bw:3.13 GB/sec
|跨节点主机到服务|1bytes      |tcp_lat:15.9 us<br>tcp_bw:1.19 MB/sec|tcp_lat:19.2 us<br>tcp_bw:1.17 MB/sec | tcp_lat:28.9 us<br>tcp_bw:1.08 MB/sec
|               |64bytes    |tcp_lat:16.2 us<br>tcp_bw:54.2 MB/sec|tcp_lat:18.8 us<br>tcp_bw:69 MB/sec    | tcp_lat:29.7 us<br>tcp_bw:55 MB/sec
|               |4KiB       |tcp_lat:28.8 us<br>tcp_bw:1.13 GB/sec|tcp_lat:34.1 us<br>tcp_bw:1.16 GB/sec  | tcp_lat:45.7 us<br>tcp_bw:1.13 GB/sec
|               |256KiB     |tcp_lat:315 us<br>tcp_bw:1.15 GB/sec|tcp_lat:397 us<br>tcp_bw:1.17 GB/sec    | tcp_lat:338 us<br>tcp_bw:1.16 GB/sec


# wrk测试
wrk是一款高性能的http请求压测工具，它使用了Epoll模型，使所有请求都是异步非阻塞模式的，因此对系统资源能够应用到极致，可以压满 cpu。
```bash
# 64个线程，20000个连接，压测时间3m
./wrk -t64 -c20000 -d3m <url>
```


TPS：每秒处理的事务数（比如每秒处理的订单数）
QPS：每秒处理的请求数
-c, --connections <N>  跟服务器建立并保持的TCP连接数量  
-d, --duration    <T>  压测时间           
-t, --threads     <N>  使用多少个线程进行压测   
-R, --rate        <T>  工作速率（吞吐量）即每个线程每秒钟完成的请求数


## 64线程20000连接数3min测试结果:
|             | QPS       |   TPS    | Latency(avg) |    error    |
|-------------|-----------|----------|--------------|-------------|
|宿主机服务端  | 235282.32 | 64.22MB  | 4.63ms   | read 1365980, write 137402, timeout 11
|runc容器服务端| 104251.00 | 26.85MB  | 5.30ms   | read 1802279, write 158301, timeout 165
|kata容器服务端| 99475.75  | 25.62MB  | 5.55ms   | read 1859794, write 157643, timeout 158

![](../images/wrk-t64-c20000.png)


# 测试数据
## qperf
### 宿主机
```bash
[root@telecom-k8s-phy01 kbuser]# qperf 10.96.0.1  -oo msg_size:1:256K:*64 -vu tcp_bw tcp_lat
tcp_bw:
    bw        =  1.13 MB/sec
    msg_size  =     1 bytes
tcp_bw:
    bw        =  77 MB/sec
    msg_size  =  64 bytes
tcp_bw:
    bw        =  1.71 GB/sec
    msg_size  =     4 KiB (4,096)
tcp_bw:
    bw        =  3.48 GB/sec
    msg_size  =   256 KiB (262,144)
tcp_lat:
    latency   =  13 us
    msg_size  =   1 bytes
tcp_lat:
    latency   =  12.9 us
    msg_size  =    64 bytes
tcp_lat:
    latency   =  14.1 us
    msg_size  =     4 KiB (4,096)
tcp_lat:
    latency   =  68.8 us
    msg_size  =   256 KiB (262,144)
[root@telecom-k8s-phy02 kbuser]# qperf 10.96.0.1  -oo msg_size:1:256K:*64 -vu tcp_bw tcp_lat
tcp_bw:
    bw        =  1.19 MB/sec
    msg_size  =     1 bytes
tcp_bw:
    bw        =  54.2 MB/sec
    msg_size  =    64 bytes
tcp_bw:
    bw        =  1.13 GB/sec
    msg_size  =     4 KiB (4,096)
tcp_bw:
    bw        =  1.15 GB/sec
    msg_size  =   256 KiB (262,144)
tcp_lat:
    latency   =  15.9 us
    msg_size  =     1 bytes
tcp_lat:
    latency   =  16.2 us
    msg_size  =    64 bytes
tcp_lat:
    latency   =  28.8 us
    msg_size  =     4 KiB (4,096)
tcp_lat:
    latency   =  315 us
    msg_size  =  256 KiB (262,144)
```
### runc
```bash
[root@telecom-k8s-phy01 kbuser]# qperf 10.196.192.101 -lp 4000 -ip 4001 -oo msg_size:1:256K:*64 -vu tcp_bw tcp_lat
tcp_bw:
    bw        =   1.14 MB/sec
    msg_size  =      1 bytes
    port      =  4,001
tcp_bw:
    bw        =   42.6 MB/sec
    msg_size  =     64 bytes
    port      =  4,001
tcp_bw:
    bw        =    335 MB/sec
    msg_size  =      4 KiB (4,096)
    port      =  4,001
tcp_bw:
    bw        =   3.39 GB/sec
    msg_size  =    256 KiB (262,144)
    port      =  4,001
tcp_lat:
    latency   =   12.5 us
    msg_size  =      1 bytes
    port      =  4,001
tcp_lat:
    latency   =   12.6 us
    msg_size  =     64 bytes
    port      =  4,001
tcp_lat:
    latency   =   20.2 us
    msg_size  =      4 KiB (4,096)
    port      =  4,001
tcp_lat:
    latency   =   92.1 us
    msg_size  =    256 KiB (262,144)
    port      =  4,001


[root@telecom-k8s-phy02 kbuser]# qperf 10.196.192.101 -lp 4000 -ip 4001 -oo msg_size:1:256K:*64 -vu tcp_bw tcp_lat
tcp_bw:
    bw        =   1.17 MB/sec
    msg_size  =      1 bytes
    port      =  4,001
tcp_bw:
    bw        =     69 MB/sec
    msg_size  =     64 bytes
    port      =  4,001
tcp_bw:
    bw        =   1.16 GB/sec
    msg_size  =      4 KiB (4,096)
    port      =  4,001
tcp_bw:
    bw        =   1.17 GB/sec
    msg_size  =    256 KiB (262,144)
    port      =  4,001
tcp_lat:
    latency   =   19.2 us
    msg_size  =      1 bytes
    port      =  4,001
tcp_lat:
    latency   =   18.8 us
    msg_size  =     64 bytes
    port      =  4,001
tcp_lat:
    latency   =   34.1 us
    msg_size  =      4 KiB (4,096)
    port      =  4,001
tcp_lat:
    latency   =    397 us
    msg_size  =    256 KiB (262,144)
    port      =  4,001

```

### kata
```bash
[root@telecom-k8s-phy01 kbuser]# qperf 10.196.192.79 -lp 4000 -ip 4001 -oo msg_size:1:256K:*64 -vu tcp_bw tcp_lattcp_bw:
    bw        =   1.05 MB/sec
    msg_size  =      1 bytes
    port      =  4,001
tcp_bw:
    bw        =   45.6 MB/sec
    msg_size  =     64 bytes
    port      =  4,001
tcp_bw:
    bw        =    355 MB/sec
    msg_size  =      4 KiB (4,096)
    port      =  4,001
tcp_bw:
    bw        =   3.13 GB/sec
    msg_size  =    256 KiB (262,144)
    port      =  4,001
tcp_lat:
    latency   =   20.2 us
    msg_size  =      1 bytes
    port      =  4,001
tcp_lat:
    latency   =   19.3 us
    msg_size  =     64 bytes
    port      =  4,001
tcp_lat:
    latency   =   25.7 us
    msg_size  =      4 KiB (4,096)
    port      =  4,001
tcp_lat:
    latency   =    167 us
    msg_size  =    256 KiB (262,144)
    port      =  4,001
[root@telecom-k8s-phy02 kbuser]# qperf 10.196.192.79 -lp 4000 -ip 4001 -oo msg_size:1:256K:*64 -vu tcp_bw tcp_lattcp_bw:
    bw        =   1.08 MB/sec
    msg_size  =      1 bytes
    port      =  4,001
tcp_bw:
    bw        =     55 MB/sec
    msg_size  =     64 bytes
    port      =  4,001
tcp_bw:
    bw        =   1.13 GB/sec
    msg_size  =      4 KiB (4,096)
    port      =  4,001
tcp_bw:
    bw        =   1.16 GB/sec
    msg_size  =    256 KiB (262,144)
    port      =  4,001
tcp_lat:
    latency   =   28.9 us
    msg_size  =      1 bytes
    port      =  4,001
tcp_lat:
    latency   =   29.7 us
    msg_size  =     64 bytes
    port      =  4,001
tcp_lat:
    latency   =   45.7 us
    msg_size  =      4 KiB (4,096)
    port      =  4,001
tcp_lat:
    latency   =    338 us
    msg_size  =    256 KiB (262,144)
    port      =  4,001
```

## wrk
```bash
# 宿主机
[root@telecom-k8s-phy03 wrk-master]# ./wrk -t64 -c20000 -d3m http://10.96.0.2:40080/
Running 3m test @ http://10.96.0.2:40080/
  64 threads and 20000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.63ms   27.66ms   1.99s    97.07%
    Req/Sec     3.76k     2.33k   51.04k    81.79%
  42374542 requests in 3.00m, 11.29GB read
  Socket errors: connect 0, read 1365980, write 137402, timeout 113
Requests/sec: 235282.32
Transfer/sec:     64.22MB


# runc
[root@telecom-k8s-phy03 wrk-master]# ./wrk -t64 -c20000 -d3m http://10.96.0.2:21020/
Running 3m test @ http://10.96.0.2:21020/
  64 threads and 20000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.30ms   23.22ms   1.81s    96.07%
    Req/Sec     1.64k     1.32k   25.39k    69.10%
  18774201 requests in 3.00m, 4.72GB read
  Socket errors: connect 0, read 1802279, write 158301, timeout 165
Requests/sec: 104251.00
Transfer/sec:     26.85MB


# kata
[root@telecom-k8s-phy03 wrk-master]# ./wrk -t64 -c20000 -d3m http://10.96.0.2:28394/
Running 3m test @ http://10.96.0.2:28394/
  64 threads and 20000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.55ms   24.44ms   1.81s    95.84%
    Req/Sec     1.57k     1.32k   26.18k    72.43%
  17913983 requests in 3.00m, 4.51GB read
  Socket errors: connect 0, read 1859794, write 157643, timeout 158
Requests/sec:  99475.75
Transfer/sec:     25.62MB



```