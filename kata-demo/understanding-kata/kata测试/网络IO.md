 [TOC]


# 当前容器云网络模型及网络开销说明
k8s的最小调度单位为pod，pod网络通信的实现依赖于第三方插件；容器云使用calico纯三层虚拟网络方案，可以避免与其他二层方案相关的数据包封装的操作，中间没有任何的NAT，没有任何的overlay，几乎能达到主机性能。
kata增加了一层tcfilter（默认）打通虚拟机与容器之间的网络

**备注：**
- 网络性能主要有两个指标是带宽和延时。延迟决定最大的QPS(Query Per Second)，而带宽决定了可支撑的最大负荷。
- qperf和iperf/netperf一样可以评测两个节点之间的带宽和延时。可以在测试tcp/ip协议和RDMA传输。不过相比netperf和iperf，支持RDMA是qperf工具的独有特性。
- wrk是一款高性能的http请求压测工具，它使用了Epoll模型，使所有请求都是异步非阻塞模式的，因此对系统资源能够应用到极致，可以压满 cpu。


# 测试配置说明
- 限制容器request/limit 1C2G
- kata设置debug_console_enabled=true（虚拟机开销占用业务开销）
- kata设置debug_console_enabled=false（虚拟机开销不限制）

``` bash
[root@telecom-k8s-phy01 hff]# kubectl get pod -o wide
NAME                                READY   STATUS    RESTARTS   AGE    IP               NODE                NOMINATED NODE   READINESS GATES
qperf-server-kata-fb5b7d54d-42sn4   1/1     Running   0          3m6s   10.196.142.196   telecom-k8s-phy02   <none>           <none>
qperf-server-z5rwk                  1/1     Running   0          16h    10.196.192.101   telecom-k8s-phy01   <none>           <none>
[root@telecom-k8s-phy01 hff]# kubectl get svc -o wide
NAME                TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                               AGE     SELECTOR
qperf-server        ClusterIP   10.196.99.250   <none>        4000/UDP,4000/TCP,4001/UDP,4001/TCP   16h     k8s-app=qperf-server
qperf-server-kata   ClusterIP   10.196.126.41   <none>        4000/UDP,4000/TCP,4001/UDP,4001/TCP   3m59s   k8s-app=qperf-server-kata
```

# qperf

- 循环测试1bytes-64KiB的带宽和延迟
```bash
qperf <ip> -oo msg_size:1:256K:*64 -vu tcp_bw tcp_lat
```

## 测试结果:
|               |msg_size   |宿主机服务端|runc容器服务端|kata容器服务端|
|---------------|-----------|-----------|-------------|-------------|
| 从主机到服务  |1bytes      |tcp_lat:13 us<br>tcp_bw:1.13 MB/sec |tcp_lat:12.5 us<br>tcp_bw:1.14 MB/sec | tcp_lat:20.2 us<br>tcp_bw:1.05 MB/sec
|               |64bytes    |tcp_lat:12.9 us<br>tcp_bw:77 MB/sec |tcp_lat:12.6 us<br>tcp_bw:42.6 MB/sec | tcp_lat:19.3 us<br>tcp_bw:45.6 MB/sec
|               |4KiB       |tcp_lat:14.1 us<br>tcp_bw:1.71 GB/sec |tcp_lat:20.2 us<br>tcp_bw:335 MB/sec | tcp_lat:25.7 us<br>tcp_bw:355 MB/sec
|               |256KiB     |tcp_lat:68.8 us<br>tcp_bw:3.48 GB/sec |tcp_lat:92.1 us<br>tcp_bw:3.39 GB/sec | tcp_lat:167 us<br>tcp_bw:3.13 GB/sec
|跨节点主机到服务|1bytes     |tcp_lat:15.9 us<br>tcp_bw:1.19 MB/sec|tcp_lat:19.2 us<br>tcp_bw:1.17 MB/sec   | tcp_lat:28.9 us<br>tcp_bw:1.08 MB/sec
|               |64bytes    |tcp_lat:16.2 us<br>tcp_bw:54.2 MB/sec|tcp_lat:18.8 us<br>tcp_bw:69 MB/sec    | tcp_lat:29.7 us<br>tcp_bw:55 MB/sec
|               |4KiB       |tcp_lat:28.8 us<br>tcp_bw:1.13 GB/sec|tcp_lat:34.1 us<br>tcp_bw:1.16 GB/sec  | tcp_lat:45.7 us<br>tcp_bw:1.13 GB/sec
|               |256KiB     |tcp_lat:315 us<br>tcp_bw:1.15 GB/sec|tcp_lat:397 us<br>tcp_bw:1.17 GB/sec  | tcp_lat:338 us<br>tcp_bw:1.16 GB/sec





## 注意
qperf测试service有问题：
1. 对于runc容器，需要修改qperf服务监听端口，否则跨节点无法测试
2. 对于kata容器，本节点和跨节点都无法测试


# wrk测试
```bash
# 2个线程，100个连接，持续时间30s 每秒r个请求
./wrk -t2 -c100 -d30 -R<r> <url> -L
```

**备注：**
宿主机及pod服务起apache（httpd）服务，pod通过nodeport提供测试curl
TPS：每秒处理的事务数（比如每秒处理的订单数）
QPS：每秒处理的请求数
-c, --connections <N>  跟服务器建立并保持的TCP连接数量  
-d, --duration    <T>  压测时间           
-t, --threads     <N>  使用多少个线程进行压测   
-R, --rate        <T>  工作速率（吞吐量）即每个线程每秒钟完成的请求数

```bash
[root@telecom-k8s-phy01 hff]# kubectl get svc
NAME                TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                               AGE
test-kata-httpd     NodePort    10.196.97.40    <none>        80:28394/TCP                          50s
test-runc-httpd     NodePort    10.196.26.12    <none>        80:21020/TCP                          15s

# 宿主机：
./wrk -t2 -c100 -d30 -R3000 http://10.96.0.2:40080/ -L

# runc:
./wrk -t2 -c100 -d30 -R3000 http://10.96.0.2:21020/ -L

# kata:
./wrk -t2 -c100 -d30 -R3000 http://10.96.0.2:28394/ -L
```

## 测试结果:
|-R（吞吐量）     |宿主机服务端|runc容器服务端|kata容器服务端|
|-----------|-----------|-------------|-------------|
|1000       |QPS: 997.37<br>TPS: 278.69KB<br>timeout: 189<br>Latency: 848.14us|QPS: 997.35<br>TPS: 263.10KB<br>timeout:\ <br>Latency: 829.26us|QPS: 997.38<br>TPS: 263.10KB<br>timeout:\ <br>Latency: 1.02ms|
|3000       |QPS: 2989.21<br>TPS: 835.37KB<br>timeout: 257<br>Latency: 330.64ms|QPS: 2989.19<br>TPS: 788.66K<br>timeout:\ <br>Latency: 845.25us|QPS: 2989.23<br>TPS: 788.67KB<br>timeout:\ <br>Latency: 1.11ms|
|5000       |QPS: 4980.69<br>TPS: 1.36MB<br>timeout: 7<br>Latency: 58.18ms|QPS: 4980.73<br>TPS: 1.28MB<br>timeout:\ <br>Latency: 0.88ms|QPS: 3460.70<br>TPS: 0.89MB<br>timeout: 377<br>Latency: 1.51ms|
|10000      |QPS: 9959.90<br>TPS: 2.72MB<br>timeout:\ <br>Latency: 0.93ms|QPS: 5593.87<br>TPS: 1.44MB<br>timeout: 554<br>Latency: 811.38ms|QPS: 5011.95<br>TPS: 1.29MB<br>timeout: 591<br>Latency: 2.87s|
|20000      |QPS: 19917.95<br>TPS: 5.44MB<br>timeout:\ <br>Latency: 1.03ms|QPS: 5868.37<br>TPS: 1.51MB<br>timeout: 519<br>Latency: 10.76s|QPS:4777.17 <br>TPS: 1.23MB<br>timeout:\ <br>Latency: 15.06s|
|50000      |QPS: 49794.13<br>TPS: 13.59MB<br>timeout: 14<br>Latency: 120.71ms|QPS: 5729.92<br>TPS: 1.48MB<br>timeout: 413<br>Latency: 16.84s|QPS: 4804.63<br>TPS: 1.24MB<br>timeout:\ <br>Latency: 18.38s|


QPS: <br>TPS: <br>timeout: <br>Latency: |


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
### 宿主机
```bash
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R1000 http://10.96.0.2:40080/ -L
Running 30s test @ http://10.96.0.2:40080/
  2 threads and 100 connections
  Thread calibration: mean lat.: 1452.784ms, rate sampling interval: 9437ms
  Thread calibration: mean lat.: 1474.524ms, rate sampling interval: 9445ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   848.14us  414.76us   4.63ms   64.28%
    Req/Sec   499.50      0.50   500.00    100.00%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%  832.00us
 75.000%    1.15ms
 90.000%    1.40ms
 99.000%    1.81ms
 99.900%    2.03ms
 99.990%    4.43ms
 99.999%    4.63ms
100.000%    4.63ms

。。。
#[Mean    =        0.848, StdDeviation   =        0.415]
#[Max     =        4.628, Total count    =        19740]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  29923 requests in 30.00s, 8.17MB read
  Socket errors: connect 0, read 0, write 0, timeout 189
Requests/sec:    997.37
Transfer/sec:    278.69KB
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R3000 http://10.96.0.2:40080/ -L
Running 30s test @ http://10.96.0.2:40080/
  2 threads and 100 connections
  Thread calibration: mean lat.: 1458.655ms, rate sampling interval: 5877ms
  Thread calibration: mean lat.: 1472.184ms, rate sampling interval: 5890ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   330.64ms  750.20ms   3.30s    86.39%
    Req/Sec     1.54k    70.36     1.64k    66.67%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    1.01ms
 75.000%    1.66ms
 90.000%    1.56s
 99.000%    3.06s
 99.900%    3.27s
 99.990%    3.30s
 99.999%    3.30s
100.000%    3.30s
。。。
#[Mean    =      330.642, StdDeviation   =      750.200]
#[Max     =     3299.328, Total count    =        60648]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  89680 requests in 30.00s, 24.47MB read
  Socket errors: connect 0, read 0, write 0, timeout 257
Requests/sec:   2989.21
Transfer/sec:    835.37KB

[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R5000 http://10.96.0.2:40080/ -L
Running 30s test @ http://10.96.0.2:40080/
  2 threads and 100 connections
  Thread calibration: mean lat.: 396.697ms, rate sampling interval: 2820ms
  Thread calibration: mean lat.: 400.296ms, rate sampling interval: 2828ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    58.18ms  754.88ms  15.00s    99.28%
    Req/Sec     2.50k    75.97     2.75k    92.31%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    0.91ms
 75.000%    1.18ms
 90.000%    1.46ms
 99.000%    1.89ms
 99.900%   13.03s
 99.990%   14.80s
 99.999%   14.98s
100.000%   15.01s

。。。
#[Mean    =       58.181, StdDeviation   =      754.876]
#[Max     =    14999.552, Total count    =        98749]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  149426 requests in 30.00s, 40.78MB read
  Socket errors: connect 0, read 0, write 1, timeout 7
Requests/sec:   4980.69
Transfer/sec:      1.36MB


[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R10000 http://10.96.0.2:40080/ -L
Running 30s test @ http://10.96.0.2:40080/
  2 threads and 100 connections
  Thread calibration: mean lat.: 184.304ms, rate sampling interval: 1388ms
  Thread calibration: mean lat.: 187.370ms, rate sampling interval: 1396ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     0.93ms  435.61us   7.03ms   66.44%
    Req/Sec     5.00k     6.11     5.02k    82.14%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    0.93ms
 75.000%    1.24ms
 90.000%    1.45ms
 99.000%    1.98ms
 99.900%    2.17ms
 99.990%    4.66ms
 99.999%    6.99ms
100.000%    7.03ms

。。。
#[Mean    =        0.929, StdDeviation   =        0.436]
#[Max     =        7.032, Total count    =       197500]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  298812 requests in 30.00s, 81.55MB read
Requests/sec:   9959.90
Transfer/sec:      2.72MB
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R20000 http://10.96.0.2:40080/ -L
Running 30s test @ http://10.96.0.2:40080/
  2 threads and 100 connections
  Thread calibration: mean lat.: 100.146ms, rate sampling interval: 685ms
  Thread calibration: mean lat.: 100.207ms, rate sampling interval: 685ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.03ms  463.70us   7.54ms   62.64%
    Req/Sec    10.01k    23.37    10.06k    69.64%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    1.04ms
 75.000%    1.39ms
 90.000%    1.63ms
 99.000%    1.99ms
 99.900%    2.32ms
 99.990%    5.36ms
 99.999%    7.14ms
100.000%    7.54ms

。。。
#[Mean    =        1.028, StdDeviation   =        0.464]
#[Max     =        7.536, Total count    =       395000]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  597555 requests in 30.00s, 163.09MB read
Requests/sec:  19917.95
Transfer/sec:      5.44MB

```

### runc
```bash
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R1000 http://10.96.0.2:21020/ -L
Running 30s test @ http://10.96.0.2:21020/
  2 threads and 100 connections
  Thread calibration: mean lat.: 0.798ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 0.806ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   829.26us  400.32us   6.26ms   64.88%
    Req/Sec   535.24    120.86     1.00k    53.13%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%  812.00us
 75.000%    1.13ms
 90.000%    1.33ms
 99.000%    1.82ms
 99.900%    2.04ms
 99.990%    3.71ms
 99.999%    6.27ms
100.000%    6.27ms

。。。
#[Mean    =        0.829, StdDeviation   =        0.400]
#[Max     =        6.264, Total count    =        19740]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  29922 requests in 30.00s, 7.71MB read
Requests/sec:    997.35
Transfer/sec:    263.10KB

[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R3000 http://10.96.0.2:21020/ -L
Running 30s test @ http://10.96.0.2:21020/
  2 threads and 100 connections
  Thread calibration: mean lat.: 0.837ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 0.856ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   845.25us  401.18us   5.81ms   64.23%
    Req/Sec     1.59k   147.83     2.11k    57.86%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%  843.00us
 75.000%    1.15ms
 90.000%    1.34ms
 99.000%    1.82ms
 99.900%    2.12ms
 99.990%    3.94ms
 99.999%    4.83ms
100.000%    5.81ms

。。。
#[Mean    =        0.845, StdDeviation   =        0.401]
#[Max     =        5.808, Total count    =        59242]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  89676 requests in 30.00s, 23.11MB read
Requests/sec:   2989.19
Transfer/sec:    788.66KB
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R5000 http://10.96.0.2:21020/ -L
Running 30s test @ http://10.96.0.2:21020/
  2 threads and 100 connections
  Thread calibration: mean lat.: 0.886ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 0.897ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     0.88ms  447.03us   7.49ms   64.01%
    Req/Sec     2.61k   271.88     3.89k    83.97%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    0.85ms
 75.000%    1.17ms
 90.000%    1.45ms
 99.000%    1.93ms
 99.900%    2.13ms
 99.990%    5.59ms
 99.999%    7.16ms
100.000%    7.49ms

。。。
#[Mean    =        0.875, StdDeviation   =        0.447]
#[Max     =        7.492, Total count    =        98748]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  149426 requests in 30.00s, 38.50MB read
Requests/sec:   4980.73
Transfer/sec:      1.28MB

[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R10000 http://10.96.0.2:21020/ -L
Running 30s test @ http://10.96.0.2:21020/
  2 threads and 100 connections
  Thread calibration: mean lat.: 1029.016ms, rate sampling interval: 4849ms
  Thread calibration: mean lat.: 1304.352ms, rate sampling interval: 5758ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   811.38ms    1.12s    4.64s    82.11%
    Req/Sec     2.84k   262.88     3.08k    85.71%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%   89.47ms
 75.000%    1.41s
 90.000%    2.62s
 99.000%    4.04s
 99.900%    4.52s
 99.990%    4.61s
 99.999%    4.64s
100.000%    4.64s

。。。
#[Mean    =      811.381, StdDeviation   =     1118.025]
#[Max     =     4636.672, Total count    =       110042]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  167825 requests in 30.00s, 43.24MB read
  Socket errors: connect 0, read 0, write 0, timeout 554
Requests/sec:   5593.87
Transfer/sec:      1.44MB
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R20000 http://10.96.0.2:21020/ -L
Running 30s test @ http://10.96.0.2:21020/
  2 threads and 100 connections
  Thread calibration: mean lat.: 3219.057ms, rate sampling interval: 11321ms
  Thread calibration: mean lat.: 3366.027ms, rate sampling interval: 12115ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    10.76s     2.51s   16.71s    62.26%
    Req/Sec     2.93k   228.00     3.16k    50.00%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%   10.96s
 75.000%   12.81s
 90.000%   13.91s
 99.000%   15.51s
 99.900%   16.61s
 99.990%   16.68s
 99.999%   16.72s
100.000%   16.72s

。。。
#[Mean    =    10756.578, StdDeviation   =     2507.203]
#[Max     =    16711.680, Total count    =       118064]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  176297 requests in 30.04s, 45.43MB read
  Socket errors: connect 0, read 0, write 0, timeout 519
Requests/sec:   5868.37
Transfer/sec:      1.51MB
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R50000 http://10.96.0.2:40080/ -L
Running 30s test @ http://10.96.0.2:40080/
  2 threads and 100 connections
  Thread calibration: mean lat.: 37.812ms, rate sampling interval: 251ms
  Thread calibration: mean lat.: 38.001ms, rate sampling interval: 252ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   120.71ms    1.09s   15.02s    98.54%
    Req/Sec    25.27k     2.09k   42.12k    96.79%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    0.99ms
 75.000%    1.40ms
 90.000%    1.80ms
 99.000%    5.59s
 99.900%   14.09s
 99.990%   14.94s
 99.999%   15.02s
100.000%   15.03s

。。。
#[Mean    =      120.708, StdDeviation   =     1088.328]
#[Max     =    15024.128, Total count    =       996000]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  1493868 requests in 30.00s, 407.72MB read
  Socket errors: connect 0, read 0, write 2, timeout 14
Requests/sec:  49794.13
Transfer/sec:     13.59MB
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R50000 http://10.96.0.2:21020/ -L
Running 30s test @ http://10.96.0.2:21020/
  2 threads and 100 connections
  Thread calibration: mean lat.: 4660.332ms, rate sampling interval: 16351ms
  Thread calibration: mean lat.: 4280.951ms, rate sampling interval: 15745ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    16.84s     4.50s   25.51s    57.82%
    Req/Sec     2.90k   504.00     3.40k    50.00%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%   16.86s
 75.000%   20.79s
 90.000%   22.97s
 99.000%   24.48s
 99.900%   25.08s
 99.990%   25.44s
 99.999%   25.53s
100.000%   25.53s

。。。
#[Mean    =    16844.032, StdDeviation   =     4499.305]
#[Max     =    25509.888, Total count    =       115305]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  172141 requests in 30.04s, 44.35MB read
  Socket errors: connect 0, read 0, write 0, timeout 413
Requests/sec:   5729.92
Transfer/sec:      1.48MB
```

### kata
```bash
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R1000 http://10.96.0.2:28394/ -L
Running 30s test @ http://10.96.0.2:28394/
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.033ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.063ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.02ms  433.27us   6.67ms   64.02%
    Req/Sec   522.85    115.86     1.00k    55.54%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    0.99ms
 75.000%    1.33ms
 90.000%    1.61ms
 99.000%    2.00ms
 99.900%    2.29ms
 99.990%    4.30ms
 99.999%    6.67ms
100.000%    6.67ms

。。。
#[Mean    =        1.020, StdDeviation   =        0.433]
#[Max     =        6.668, Total count    =        19740]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  29922 requests in 30.00s, 7.71MB read
Requests/sec:    997.38
Transfer/sec:    263.10KB

[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R3000 http://10.96.0.2:28394/ -L
Running 30s test @ http://10.96.0.2:28394/
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.229ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.232ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.11ms  842.48us  19.78ms   97.18%
    Req/Sec     1.58k   198.38     3.89k    91.10%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    1.05ms
 75.000%    1.38ms
 90.000%    1.63ms
 99.000%    2.32ms
 99.900%   13.86ms
 99.990%   18.21ms
 99.999%   19.39ms
100.000%   19.79ms

。。。
#[Mean    =        1.108, StdDeviation   =        0.842]
#[Max     =       19.776, Total count    =        59250]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  89679 requests in 30.00s, 23.11MB read
Requests/sec:   2989.23
Transfer/sec:    788.67KB
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R5000 http://10.96.0.2:28394/ -L
Running 30s test @ http://10.96.0.2:28394/
  2 threads and 100 connections
  Thread calibration: mean lat.: 6.374ms, rate sampling interval: 41ms
  Thread calibration: mean lat.: 6.224ms, rate sampling interval: 40ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.51ms  610.88us  17.07ms   71.69%
    Req/Sec     1.44k    59.78     2.08k    73.61%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    1.45ms
 75.000%    1.87ms
 90.000%    2.26ms
 99.000%    2.92ms
 99.900%    4.70ms
 99.990%   15.11ms
 99.999%   16.86ms
100.000%   17.09ms

。。。
#[Mean    =        1.512, StdDeviation   =        0.611]
#[Max     =       17.072, Total count    =        56344]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  103828 requests in 30.00s, 26.75MB read
  Socket errors: connect 0, read 0, write 0, timeout 377
Requests/sec:   3460.70
Transfer/sec:      0.89MB

[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R10000 http://10.96.0.2:28394/ -L
Running 30s test @ http://10.96.0.2:28394/
  2 threads and 100 connections
  Thread calibration: mean lat.: 874.239ms, rate sampling interval: 4110ms
  Thread calibration: mean lat.: 673.307ms, rate sampling interval: 3287ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.87s     2.53s   17.22s    86.14%
    Req/Sec     2.41k   688.87     3.35k    50.00%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    2.59s
 75.000%    3.56s
 90.000%    4.94s
 99.000%   16.57s
 99.900%   17.01s
 99.990%   17.20s
 99.999%   17.24s
100.000%   17.24s

。。。
#[Mean    =     2872.981, StdDeviation   =     2529.939]
#[Max     =    17219.584, Total count    =       100322]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  150487 requests in 30.03s, 38.78MB read
  Socket errors: connect 0, read 0, write 2, timeout 591
Requests/sec:   5011.95
Transfer/sec:      1.29MB
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R20000 http://10.96.0.2:28394/ -L
Running 30s test @ http://10.96.0.2:28394/
  2 threads and 100 connections
  Thread calibration: mean lat.: 3802.356ms, rate sampling interval: 14344ms
  Thread calibration: mean lat.: 3810.352ms, rate sampling interval: 14385ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    15.06s     4.21s   23.59s    59.80%
    Req/Sec     2.53k     3.00     2.53k    50.00%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%   15.01s
 75.000%   18.61s
 90.000%   20.89s
 99.000%   22.69s
 99.900%   23.28s
 99.990%   23.46s
 99.999%   23.61s
100.000%   23.61s

。。。
#[Mean    =    15057.365, StdDeviation   =     4210.496]
#[Max     =    23592.960, Total count    =        97620]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  143315 requests in 30.00s, 36.93MB read
Requests/sec:   4777.17
Transfer/sec:      1.23MB
[root@telecom-k8s-phy03 wrk2-master]# ./wrk -t2 -c100 -d30 -R50000 http://10.96.0.2:28394/ -L
Running 30s test @ http://10.96.0.2:28394/
  2 threads and 100 connections
  Thread calibration: mean lat.: 4272.538ms, rate sampling interval: 16121ms
  Thread calibration: mean lat.: 4316.274ms, rate sampling interval: 16171ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    18.38s     4.98s   27.36s    59.32%
    Req/Sec     2.43k     5.00     2.44k    50.00%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%   18.46s
 75.000%   22.66s
 90.000%   25.10s
 99.000%   26.82s
 99.900%   27.15s
 99.990%   27.30s
 99.999%   27.38s
100.000%   27.38s

。。。
#[Mean    =    18381.631, StdDeviation   =     4983.160]
#[Max     =    27361.280, Total count    =        95898]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  144526 requests in 30.08s, 37.24MB read
Requests/sec:   4804.63
Transfer/sec:      1.24MB


```
