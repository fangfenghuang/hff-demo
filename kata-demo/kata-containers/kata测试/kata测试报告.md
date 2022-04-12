[TOC]
# 从哪些维度去考量kata containers?

1. 容器启动速度

2. cpu、内存性能损耗

3. 资源开销

4. 网络性能：见[[网络IO]]

5. 文件IO性能：见[[文件IO]]

6. 已知问题验证

7. 边界测试

8. 异常测试

## 调优方向：

使用高性能固态硬盘，或者上spdk，或者使用device mapper

打开iothread

根据不同应用设置内核调优参数

使用轻量化hypervisor（ccloud-hypervisor）

开启vm template，优化启动时间

## 测试环境

### 单机：

**10.208.11.110（单机）**

> CPU:Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz(4核8线程)
>
> Host Kernel: Linux localhost.localdomain 3.10.0-1160.59.1.el7.x86_64 #1 SMP Wed Feb 23 16:47:03 UTC 2022 x86_64 x86_64 x86_64 GNU/Linux
>
> Guest Kernel: Linux clr-64b293ce5be44f6d9f521c20c5a36249 5.15.23 #2 SMP Mon Mar 7 22:16:36 UTC 2022 x86_64 GNU/Linux

### **集群：**

**10.91.0.1-3（集群）**

> CPU: Intel(R) Xeon(R) Silver 4216 CPU @ 2.10GHz（16核64线程）
>
> Host Kernel: Linux rqy-k8s-1 3.10.0-1160.59.1.el7.x86_64 #1 SMP Wed Feb 23 16:47:03 UTC 2022 x86_64 x86_64 x86_64 GNU/Linux
>
> Guest Kernel: Linux clr-f9e79d9d0eb74bbc9f59c716e2fd9795 5.15.23 #2 SMP Mon Mar 7 22:16:36 UTC 2022 x86_64 GNU/Linux



**containerd版本1.4.6**

**kata版本2.4.0-rc0**



# 容器启动速度

kata container启动时间构成分析：

1. qemu启动+virtiofsd(vhost-user-fs server)启动+kvm创建vm资源

2. quest os内核态内核bootup

3. quest os用户态systemd+agent启动

4. vm中agent创建container+conatiner启动

## 测试1：测试启动销毁总耗时

**(容器启动+执行命令+容器销毁)**

> time ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hfftest uname -a

## 测试2： vm启动时间

> ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hfftest dmesg 

> $ dmesg | grep Startup
>
> $ systemd-analyze

## 测试数据

单机

|      | 启动销毁总耗时 | vm启动时间      |
| ---- | -------------- | --------------- |
| runc | 0.262s         | /               |
| kata | 2.744s         | 226ms（pod/vm） |

> [root@localhost ~]# time ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hello uname -a
>
> Linux clr-612c65aedeac414ea353a3af34380e25 5.15.23 #2 SMP Mon Mar 7 22:16:36 UTC 2022 x86_64 GNU/Linux
> real   0m2.744s
> user   0m0.016s
> sys   0m0.015s
>
>
> [root@localhost ~]# time ctr -n k8s.io run  -t --rm docker.io/library/busybox:latest hello uname -a
>
> Linux localhost.localdomain 3.10.0-1160.59.1.el7.x86_64 #1 SMP Wed Feb 23 16:47:03 UTC 2022 x86_64 GNU/Linux
> real   0m0.262s
> user   0m0.014s
> sys   0m0.015s
>  
>
> [root@localhost kata]# ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest dmesg
>
> / # dmesg | grep Startup
>
> [   0.228203] systemd[1]: Startup finished in 141ms (kernel) + 84ms (userspace) = 226ms.
>
>
> [root@localhost ~]# kata-runtime exec dmesg
>
> bash: grep: command not found
>
> bash: grep: command not found
>
> bash: tty: command not found
>
> bash: expr: command not found
>
> bash: [: : integer expression expected
>
> bash-5.1# dmesg
>
> [   0.228203] systemd[1]: Startup finished in 141ms (kernel) + 84ms (userspace) = 226ms.



# cpu、内存性能损耗

## CPU：

>  ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/dotnetdr/sysbench:0.5 hfftest sysbench --test=cpu run

## 内存：

> ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/dotnetdr/sysbench:0.5 hfftest sysbench --test=memory --memory-block-size=4k --memory-total-size=4G run


## 测试数据

单机

|      | cpu总耗时 | mem总耗时 |
| ---- | --------- | --------- |
| runc | 7.3730s   | 0.3924s   |
| kata | 7.4050s   | 0.3960s   |
| 主机 | 10.0004s  | 0.3175s   |


> ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/dotnetdr/sysbench:0.5 hfftest sysbench --test=cpu run
>
>
> [root@localhost kata]# ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/dotnetdr/sysbench:0.5 hfftest sysbench --test=cpu run
>
> sysbench 0.5:  multi-threaded system evaluation benchmark
>
> Running the test with following options:
>
> Number of threads: 1
>
> Random number generator seed is 0 and will be ignored
>
> Primer numbers limit: 10000
>
> Threads started!
>
> General statistics:
>
>   total time:              7.4050s
>
>   total number of events:        10000
>
>   total time taken by event execution: 7.4020s
>
>   response time:
>
> ​     min:                  0.70ms
>
> ​     avg:                  0.74ms
>
> ​     max:                  3.23ms
>
> ​     approx.  95 percentile:        0.81ms
>
>
> Threads fairness:
>
>   events (avg/stddev):      10000.0000/0.00
>
>   execution time (avg/stddev):  7.4020/0.00
>
> [root@localhost kata]# ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/dotnetdr/sysbench:0.5 hfftest sysbench --test=memory --memory-block-size=4k --memory-total-size=4G run
>
> sysbench 0.5:  multi-threaded system evaluation benchmark
>
>  
> Running the test with following options:
>
> Number of threads: 1
>
> Random number generator seed is 0 and will be ignored
>
> Threads started!
>
> Operations performed: 1048576 (2647668.86 ops/sec)
>
> 4096.00 MB transferred (10342.46 MB/sec)
>
> General statistics:
>
>   total time:              0.3960s
>
>   total number of events:        1048576
>
>   total time taken by event execution: 0.2954s
>
>   response time:
>
> ​     min:                  0.00ms
>
> ​     avg:                  0.00ms
>
> ​     max:                  0.04ms
>
> ​     approx.  95 percentile:        0.00ms
>
>  
> Threads fairness:
>
>   events (avg/stddev):      1048576.0000/0.00
>
>   execution time (avg/stddev):  0.2954/0.00
>
> [root@localhost kata]# ctr -n k8s.io run -t --rm docker.io/dotnetdr/sysbench:0.5 hfftest sysbench --test=cpu run
>
> sysbench 0.5:  multi-threaded system evaluation benchmark
>
> Running the test with following options:
>
> Number of threads: 1
>
> Random number generator seed is 0 and will be ignored
>
> Primer numbers limit: 10000
>
> Threads started!
>
> General statistics:
>
>   total time:              7.3730s
>
>   total number of events:        10000
>
>   total time taken by event execution: 7.3705s
>
>   response time:
>
> ​     min:                  0.70ms
>
> ​     avg:                  0.74ms
>
> ​     max:                  1.05ms
>
> ​     approx.  95 percentile:        0.81ms
>
> Threads fairness:
>
>   events (avg/stddev):      10000.0000/0.00
>
>   execution time (avg/stddev):  7.3705/0.00

**性能损耗低原因：**

高版本的kernal kvm模块、qemu代码做了优化，虚拟化开销<1.5%

# 资源开销

## 空负载(1C2G)
/proc/vm_pid/status
百兆级别与兆级别

|                                    | kata       | runc      |
| ---------------------------------- | ---------- | --------- |
| VmSize进程当前使用的虚拟内存的大小 | 2533708 kB | 113364 kB |
| VmRSS实际的物理内存的使用量        | 126252 kB  | 5596 kB   |
| VmHWM程序得到分配到物理内存的峰值  | 150728 kB  | 5828 kB   |
| VmData进程数据段的大小             | 249760 kB  | 103508 kB |
| %CPU                               | 0.0%       | 0.0%      |
| %MEM                               | 1.6%       | 0.1%      |

> runc:
>
> [root@localhost hff]# crictl ps
>
> CONTAINER ID     IMAGE        CREATED       STATE        NAME                  ATTEMPT       POD ID
>
> 4bfd4949e00f6    14701355bb465    35 seconds ago    Running       stress                 0          08f791371e094
>
> [root@localhost hff]# kubectl top pod
>
> NAME            CPU(cores)  MEMORY(bytes)
>
> stress-6c64c9c667-znwgr  0m      0Mi
>
>  PID USER    PR  NI   VIRT   RES   SHR S  %CPU %MEM   TIME+ COMMAND
>
> 20987 root    20  0  113364  4856  4000 S  0.0  0.1  0:00.02 containerd-shim
>
> [root@localhost hff]# cat /proc/20987/status
>
> Name:  containerd-shim
>
> Umask:  0022
>
> State:  S (sleeping)
>
> Tgid:  20987
>
> Ngid:  0
>
> Pid:   20987
>
> PPid:  1
>
> TracerPid:    0
>
> Uid:   0    0    0    0
>
> Gid:   0    0    0    0
>
> FDSize: 64
>
> Groups:
>
> VmPeak:  113364 kB
>
> VmSize:  113364 kB
>
> VmLck:     0 kB
>
> VmPin:     0 kB
>
> VmHWM:    5828 kB
>
> VmRSS:    5596 kB
>
> RssAnon:       1580 kB
>
> RssFile:       4016 kB
>
> RssShmem:        0 kB
>
> VmData:  103508 kB
>
> VmStk:    132 kB
>
> VmExe:    4652 kB
>
> VmLib:     0 kB
>
> VmPTE:     52 kB
>
> VmSwap:     0 kB
>
> Threads:     15
>
> SigQ:  0/30632
>
> SigPnd: 0000000000000000
>
> ShdPnd: 0000000000000000
>
> SigBlk: fffffffe3bfa2800
>
> SigIgn: 0000000000000000
>
> SigCgt: fffffffe7fc1feff
>
> CapInh: 0000000000000000
>
> CapPrm: 0000001fffffffff
>
> CapEff: 0000001fffffffff
>
> CapBnd: 0000001fffffffff
>
> CapAmb: 0000000000000000
>
> NoNewPrivs:   0
>
> Seccomp:     0
>
> Speculation_Store_Bypass:    thread vulnerable
>
> Cpus_allowed:  ff
>
> Cpus_allowed_list:    0-7
>
> Mems_allowed:  00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000001
>
> Mems_allowed_list:    0
>
> voluntary_ctxt_switches:     9
>
> nonvoluntary_ctxt_switches:   5
>
> kata:
>
> [root@localhost hff]# kubectl top pod
>
> NAME            CPU(cores)  MEMORY(bytes)
>
> stress-64b4544d6b-zbkhh  1m      0Mi
>
>  PID USER    PR  NI   VIRT   RES   SHR S  %CPU %MEM   TIME+ COMMAND
>
>  1370 root    20  0 2533708 126252 112092 S  0.0  1.6  0:00.50 qemu-system-x86
>
>
> [root@localhost hff]# cat /proc/1370/status
>
> Name:  qemu-system-x86
>
> Umask:  0027
>
> State:  S (sleeping)
>
> Tgid:  1370
>
> Ngid:  0
>
> Pid:   1370
>
> PPid:  1
>
> TracerPid:    0
>
> Uid:   0    0    0    0
>
> Gid:   0    0    0    0
>
> FDSize: 128
>
> Groups:
>
> VmPeak:  2537668 kB
>
> VmSize:  2533708 kB
>
> VmLck:     0 kB
>
> VmPin:     0 kB
>
> VmHWM:   150728 kB
>
> VmRSS:   126252 kB
>
> RssAnon:      14160 kB
>
> RssFile:      33668 kB
>
> RssShmem:      78424 kB
>
> VmData:  249760 kB
>
> VmStk:    132 kB
>
> VmExe:    8460 kB
>
> VmLib:     0 kB
>
> VmPTE:    548 kB
>
> VmSwap:     0 kB
>
> Threads:     4
>
> SigQ:  0/30632
>
> SigPnd: 0000000000000000
>
> ShdPnd: 0000000000000000
>
> SigBlk: 0000000010002240
>
> SigIgn: 0000000000381000
>
> SigCgt: 0000000180004243
>
> CapInh: 0000000000000000
>
> CapPrm: 0000001fffffffff
>
> CapEff: 0000001fffffffff
>
> CapBnd: 0000001fffffffff
>
> CapAmb: 0000000000000000
>
> NoNewPrivs:   0
>
> Seccomp:     0
>
> Speculation_Store_Bypass:    thread vulnerable
>
> Cpus_allowed:  ff
>
> Cpus_allowed_list:    0-7
>
> Mems_allowed:  00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000001
>
> Mems_allowed_list:    0
>
> voluntary_ctxt_switches:     502
>
> nonvoluntary_ctxt_switches:   5

## stress压测

```
containers:
      - image: vish/stress
        name: stress
        resources:
          limits:
            cpu: "2"
            memory: "2Gi"
          requests:
            cpu: "2"
            memory: "2Gi"
        args:
        - -cpus
        - "2"
        - -mem-total
        - "1Gi"
        - -mem-alloc-size
        - "100Mi"
        - -mem-alloc-sleep
        - "1s" 
```

2C2Gi(CPU打满内存使用一半)

|                                    | kata          | runc          |
| ---------------------------------- | ------------- | ------------- |
| VmSize进程当前使用的虚拟内存的大小 | 4778352 kB    | 1066304 kB    |
| VmRSS实际的物理内存的使用量        | 1218900 kB    | 1058608 kB    |
| VmHWM程序得到分配到物理内存的峰值  | 1219136 kB    | 1058740 kB    |
| VmData进程数据段的大小             | 397228 kB     | 1062520 kB    |
| %CPU                               | 200%(1998m)   | ( 200%)2000m  |
| %MEM                               | 15.4%(1035Mi) | (13.3%)1033Mi |

> kata:
>
> [root@localhost hff]# kubectl top pod
>
> NAME            CPU(cores)  MEMORY(bytes)
>
> stress-558c4747d7-dlg8s  1998m     1035Mi
>
>
>  PID USER    PR  NI   VIRT   RES   SHR S  %CPU %MEM   TIME+ COMMAND
>
>  4228 root    20  0 4778352  1.2g  1.1g S 207.1 15.4  7:15.84 qemu-system-x86
>
>
> [root@localhost hff]# cat /proc/4228/status
>
> Name:  qemu-system-x86
>
> Umask:  0027
>
> State:  S (sleeping)
>
> Tgid:  4228
>
> Ngid:  0
>
> Pid:   4228
>
> PPid:  1
>
> TracerPid:    0
>
> Uid:   0    0    0    0
>
> Gid:   0    0    0    0
>
> FDSize: 128
>
> Groups:
>
> VmPeak:  4778352 kB
>
> VmSize:  4778352 kB
>
> VmLck:     0 kB
>
> VmPin:     0 kB
>
> VmHWM:  1219136 kB
>
> VmRSS:  1219028 kB
>
> RssAnon:      14320 kB
>
> RssFile:      32560 kB
>
> RssShmem:     1172148 kB
>
> VmData:  397228 kB
>
> VmStk:    132 kB
>
> VmExe:    8460 kB
>
> VmLib:     0 kB
>
> VmPTE:    2736 kB
>
> VmSwap:     0 kB
>
> Threads:     6
>
> SigQ:  0/30632
>
> SigPnd: 0000000000000000
>
> ShdPnd: 0000000000000000
>
> SigBlk: 0000000010002240
>
> SigIgn: 0000000000381000
>
> SigCgt: 0000000180004243
>
> CapInh: 0000000000000000
>
> CapPrm: 0000001fffffffff
>
> CapEff: 0000001fffffffff
>
> CapBnd: 0000001fffffffff
>
> CapAmb: 0000000000000000
>
> NoNewPrivs:   0
>
> Seccomp:     0
>
> Speculation_Store_Bypass:    thread vulnerable
>
> Cpus_allowed:  ff
>
> Cpus_allowed_list:    0-7
>
> Mems_allowed:  00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000001
>
> Mems_allowed_list:    0
>
> voluntary_ctxt_switches:     455
>
> nonvoluntary_ctxt_switches:   32
>
>
> runc:
>
>
> [root@localhost hff]# kubectl top pod
>
> NAME            CPU(cores)  MEMORY(bytes)
>
> stress-7b5575b67f-jv7ls  2000m     1033Mi
>
> [root@localhost hff]# ps -ef | grep stress
>
> root   10905 10846 99 09:36 ?     00:16:39 /stress -logtostderr -cpus 2 -mem-total 1Gi -mem-alloc-size 100Mi -mem-alloc-sleep 1s
>
> root   31214  5109  0 09:45 pts/0   00:00:00 grep --color=auto stress
>
> [root@localhost hff]# top -p 10905
>
>  
>
>  PID USER    PR  NI   VIRT   RES   SHR S  %CPU %MEM   TIME+ COMMAND
>
> 10905 root    20  0 1066304  1.0g  1604 S 200.0 13.3  14:12.67 stress
>
> 
> [root@localhost hff]# cat /proc/10905/status
>
> Name:  stress
>
> Umask:  0022
>
> State:  S (sleeping)
>
> Tgid:  10905
>
> Ngid:  0
>
> Pid:   10905
>
> PPid:  10846
>
> TracerPid:    0
>
> Uid:   0    0    0    0
>
> Gid:   0    0    0    0
>
> FDSize: 64
>
> Groups:
>
> VmPeak:  1066304 kB
>
> VmSize:  1066304 kB
>
> VmLck:     0 kB
>
> VmPin:     0 kB
>
> VmHWM:  1058740 kB
>
> VmRSS:  1058608 kB
>
> RssAnon:     1057000 kB
>
> RssFile:       1608 kB
>
> RssShmem:        0 kB
>
> VmData:  1062520 kB
>
> VmStk:    132 kB
>
> VmExe:    1624 kB
>
> VmLib:     0 kB
>
> VmPTE:    2092 kB
>
> VmSwap:     0 kB
>
> Threads:     11
>
> SigQ:  0/30632
>
> SigPnd: 0000000000000000
>
> ShdPnd: 0000000000000000
>
> SigBlk: 0000000000000000
>
> SigIgn: 0000000000000000
>
> SigCgt: fffffffe7fc1feff
>
> CapInh: 00000000a80425fb
>
> CapPrm: 00000000a80425fb
>
> CapEff: 00000000a80425fb
>
> CapBnd: 00000000a80425fb
>
> CapAmb: 0000000000000000
>
> NoNewPrivs:   0
>
> Seccomp:     0
>
> Speculation_Store_Bypass:    thread vulnerable
>
> Cpus_allowed:  ff
>
> Cpus_allowed_list:    0-7
>
> Mems_allowed:  00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000001
>
> Mems_allowed_list:    0
>
> voluntary_ctxt_switches:     134
>
> nonvoluntary_ctxt_switches:   703



## nginx压测（4C4G）（ab压测）

> ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/httpd:latest hfftest sh
>
> ab -kn 100000 -c 100 [http://10.241.102.146](http://pod_ip)/

|      | kata           |
| ---- | -------------- |
| kata | 9704.11 req/s  |
| runc | 13740.20 req/s |

> This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
>
> Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
>
> Licensed to The Apache Software Foundation, http://www.apache.org/
>
>
> Benchmarking 10.241.102.146 (be patient)
>
> Completed 10000 requests
>
> Completed 20000 requests
>
> Completed 30000 requests
>
> Completed 40000 requests
>
> Completed 50000 requests
>
> Completed 60000 requests
>
> Completed 70000 requests
>
> Completed 80000 requests
>
> Completed 90000 requests
>
> Completed 100000 requests
>
> Finished 100000 requests
>
>
> Server Software:     nginx/1.21.5
>
> Server Hostname:     10.241.102.146
>
> Server Port:       80
>
> Document Path:      /
>
> Document Length:     615 bytes
>
>
> Concurrency Level:    30
>
> Time taken for tests:  10.305 seconds
>
> Complete requests:    100000
>
> Failed requests:     0
>
> Total transferred:    84800000 bytes
>
> HTML transferred:    61500000 bytes
>
> Requests per second:   9704.11 [#/sec] (mean)
>
> Time per request:    3.091 [ms] (mean)
>
> Time per request:    0.103 [ms] (mean, across all concurrent requests)
>
> Transfer rate:      8036.22 [Kbytes/sec] received
>
>  
>
> Connection Times (ms)
>
> ​       min  mean[+/-sd] median  max
>
> Connect:     0   1  0.4    1    11
>
> Processing:   0   2  1.3    2    27
>
> Waiting:     0   2  1.2    2    27
>
> Total:      0   3  1.3    3    33
>
>  
>
> Percentage of the requests served within a certain time (ms)
>
>  50%    3
>
>  66%    3
>
>  75%    4
>
>  80%    4
>
>  90%    4
>
>  95%    5
>
>  98%    7
>
>  99%    8
>
>  100%   33 (longest request)
>
>  
>
>  
>
> [root@localhost test]# kubectl logs apache-5f5665dfd8-vvtzk
>
> This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
>
> Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
>
> Licensed to The Apache Software Foundation, http://www.apache.org/
>
>  
>
> Benchmarking 10.241.102.146 (be patient)
>
> Completed 10000 requests
>
> Completed 20000 requests
>
> Completed 30000 requests
>
> Completed 40000 requests
>
> Completed 50000 requests
>
> Completed 60000 requests
>
> Completed 70000 requests
>
> Completed 80000 requests
>
> Completed 90000 requests
>
> Completed 100000 requests
>
> Finished 100000 requests
>
>  
>
>  
>
> Server Software:     nginx/1.21.5
>
> Server Hostname:     10.241.102.146
>
> Server Port:       80
>
>  
>
> Document Path:      /
>
> Document Length:     615 bytes
>
>  
>
> Concurrency Level:    30
>
> Time taken for tests:  7.278 seconds
>
> Complete requests:    100000
>
> Failed requests:     0
>
> Total transferred:    84800000 bytes
>
> HTML transferred:    61500000 bytes
>
> Requests per second:   13740.20 [#/sec] (mean)
>
> Time per request:    2.183 [ms] (mean)
>
> Time per request:    0.073 [ms] (mean, across all concurrent requests)
>
> Transfer rate:      11378.60 [Kbytes/sec] received
>
>  
>
> Connection Times (ms)
>
> ​       min  mean[+/-sd] median  max
>
> Connect:     0   1  0.2    1    4
>
> Processing:   0   1  0.3    1    6
>
> Waiting:     0   1  0.3    1    6
>
> Total:      1   2  0.3    2    7
>
>  
>
> Percentage of the requests served within a certain time (ms)
>
>  50%    2
>
>  66%    2
>
>  75%    2
>
>  80%    2
>
>  90%    3
>
>  95%    3
>
>  98%    3
>
>  99%    3
>
>  100%    7 (longest request)

## kata overhead

TODO



## default调整对虚拟机开销的影响

default_vcpus调整成5，开多个副本，cpu实际使用率并没有很高

> [root@localhost kata]# kubectl get pod
>
> NAME              READY  STATUS   RESTARTS  AGE
>
> test-kata-1-f9bd8f6f7-59chj  1/1   Running  0      10s
>
> test-kata-1-f9bd8f6f7-64hg5  1/1   Running  0      10s
>
> test-kata-1-f9bd8f6f7-7hxtp  1/1   Running  0      10s
>
> test-kata-1-f9bd8f6f7-jddwb  1/1   Running  0      115s
>
> test-kata-1-f9bd8f6f7-pgq4t  1/1   Running  0      10s
>
> test-kata-56b85bd45f-dr5qk   1/1   Running  0      3m25s
>
> test-kata-56b85bd45f-f8vg8   1/1   Running  0      6m9s
>
> test-kata-56b85bd45f-fss2v   1/1   Running  0      3m25s
>
> test-kata-56b85bd45f-mv9vp   1/1   Running  0      9m29s
>
> test-kata-56b85bd45f-x2tx2   1/1   Running  0      3m25s
>
> top - 09:17:12 up 9 days, 14:55,  1 user,  load average: 1.19, 0.79, 0.65
>
> Tasks: 453 total,  1 running, 452 sleeping,  0 stopped,  0 zombie
>
> %Cpu(s):  1.4 us,  1.0 sy,  0.0 ni, 95.5 id,  2.0 wa,  0.0 hi,  0.0 si,  0.0 st
>
> KiB Mem :  7935356 total,  959456 free,  4247356 used,  2728544 buff/cache
>
> KiB Swap:     0 total,     0 free,     0 used.  2054316 avail Mem

# 其他已知问题验证

## DinD

containerd runc容器支持挂在主机路径的docker.sock，但是kata容器不支持；

containerd runc容器挂载docker.sock有风险，docker重启等问题

kata容器如果需要支持docker命令需要通过DinD边车方案，实现方式参考[[流水线适配问题]]

 

## 临时存储限制ephemeral-storage（支持）

```

 runtimeClassName: kata-containers
   containers:
   - image: nginx
　　  name: nginx
　　　resources:
　　　 limits:
　　　　　ephemeral-storage: 2Gi
　　　　requests:
　　   　ephemeral-storage: 2Gi
```

Pod启动后，进入容器，执行 dd if=/dev/zero of=/test bs=4096 count=1024000 ，尝试创建一个4Gi的文件：

> [root@rqy-k8s-2 kata]# kubectl exec -it nginx-78cb94bbd5-mw6pc bash
>
> root@nginx-78cb94bbd5-mw6pc:/#
>
> root@nginx-78cb94bbd5-mw6pc:/#
>
> root@nginx-78cb94bbd5-mw6pc:/# dd if=/dev/zero of=/test bs=4096 count=1024000
>
> 1024000+0 records in
>
> 1024000+0 records out
>
> 4194304000 bytes (4.2 GB, 3.9 GiB) copied, 3.24974 s, 1.3 GB/s
>
> root@nginx-78cb94bbd5-mw6pc:/# command terminated with exit code 137
>
>  
>
> [root@rqy-k8s-2 kbuser]# kubectl get pod -w
>
> NAME           READY  STATUS   RESTARTS  AGE
>
> gitea-64f76f6567-fn72l  1/1   Running  0      14h
>
> nginx-78cb94bbd5-mw6pc  1/1   Running  0      81s
>
> nginx-78cb94bbd5-mw6pc  0/1   Evicted  0      116s
>
> nginx-78cb94bbd5-nngc4  0/1   Pending  0      1s
>
> nginx-78cb94bbd5-nngc4  0/1   Pending  0      1s
>
> nginx-78cb94bbd5-nngc4  0/1   ContainerCreating  0      1s
>
> nginx-78cb94bbd5-nngc4  1/1   Running       0      22s

## kata不支持subPaths(emptyDir )

但不影响挂载configmap使用

## kata不支持host网络

pod可以创建，但是不支持host网络特性，未使用主机端口

## privilige与runc容器不相同

### enable_cpu_memory_hotplug=false不支持--priviliged：

>  Warning  Failed   13s (x2 over 15s)  kubelet, rqy-k8s-3  Error: failed to create containerd task: Conflicting device updates for /dev/dm-1: unknown
>
>  Warning  BackOff   12s (x2 over 13s)  kubelet, rqy-k8s-3  Back-off restarting failed container

```
[plugins.cri.containerd.runtimes.kata]
  runtime_type = " io.containerd.kata.v2 " 
  privileged_without_host_devices = true
```



## 内核版本

### 5.15.23（2.4.0-rc0）

> [root@rqy-k8s-1 kbuser]# kubectl exec -it netperf-server-kata -- cat /proc/version
>
> Linux version 5.15.23 (root@655cd71e6195) (gcc (Ubuntu 9.4.0-1ubuntu1~20.04) 9.4.0, GNU ld (GNU Binutils for Ubuntu) 2.34) #2 SMP Mon Mar 7 22:16:36 UTC 2022
>
> [root@rqy-k8s-1 kbuser]# kubectl exec -it netperf-server-default -- cat /proc/version
>
> Linux version 3.10.0-1160.59.1.el7.x86_64 (mockbuild@kbuilder.bsys.centos.org) (gcc version 4.8.5 20150623 (Red Hat 4.8.5-44) (GCC) ) #1 SMP Wed Feb 23 16:47:03 UTC 2022
>
> [root@rqy-k8s-1 kbuser]#

### 5.10.25( 2.4.0-alpha2)

### 4.19.86（kata1.10.8）

# 边界测试



## cpu/mem资源限制

### 不设置limit:

默认使用default设置的cpu/mem限制

### 设置limit

容器业务最大使用上限：limit

最终VM的资源大小为：limit+default（lscpu、free -h）

VM的最大使用量：overhead+limit(memory.limit_in_bytes)(describe node) 

如果不设置request，则request 的值和 limit 默认相等

## cpu/mem超分

### cpu超分风险不大

> [root@localhost kata]# kubectl top pod
>
> NAME             CPU(cores)  MEMORY(bytes)
>
> test-1-5687dbdf79-ddxdn   4977m     1Mi
>
> test-2-5889c67db7-nm4cb   4917m     1Mi
>
> test-3-5d6bd7fb66-p9gv2   4948m     1Mi
>
> test-4-b47b67ff-srsn2    4938m     1Mi
>
> test-5-7db59fb7cb-9fxpw   4929m     1Mi
>
> test-6-566cc7dcff-wzsfv   4942m     1Mi
>
> test-7-5976f4569f-dpxd2   4904m     1Mi
>
> test-kata-cfbbd954b-jt49k  2007m     1Mi

 

### 内存超分有风险

> [root@localhost kata]# crictl stats 915587983fee9
>
> CONTAINER      CPU %        MEM         DISK         INODES
>
> 915587983fee9    95.24        3.254GB       0B          13

内存超分可能会导致业务被杀死，系统卡死

# 异常测试

## 重启节点上的kata-deploy

kata pod没有被重启