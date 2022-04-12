[TOC]

## cache与shm

> root@netperf-server-runc:/# df -h /dev/shm
> df: Warning: cannot read table of mounted file systems: No such file or directory
> Filesystem      Size  Used Avail Use% Mounted on
>
> 64M     0   64M   0% /dev/shm
>
> root@netperf-server-kata:/# df -h /dev/shm
> df: Warning: cannot read table of mounted file systems: No such file or directory
> Filesystem      Size  Used Avail Use% Mounted on
>
> 992M     0  992M   0% /dev/shm



## dd

>   ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hfftest sh
>
>   dd if=/dev/zero of=/hff/host.txt bs=4096 count=1024000 oflag=direct

rqy-k8s-3: oflag=direct

|      | count=1024000 | count=10240000 | hostpath | nfs      |
| ---- | ------------- | -------------- | -------- | -------- |
| host | 128~129MB/s   | 111 MB/s       | /        | /        |
| runc | 126~129MB/s   | 114.9MB/s      | 114MB/s  | 16.4MB/s |
| kata | 130~135 MB/s  | 128.0MB/s      | 134MB/s  | 89.3MB/s |

> [root@rqy-k8s-3 kbuser]# dd if=/dev/zero of=/hff/host.txt bs=4096 count=10240000 oflag=direct
> 10240000+0 records in
> 10240000+0 records out
> 41943040000 bytes (42 GB) copied, 376.446 s, 111 MB/s
>
> 
>
> [root@rqy-k8s-3 kbuser]# ctr -n k8s.io run  -t --rm docker.io/library/busybox:latest hfftest sh
> / # mkdir /hff
> / # dd if=/dev/zero of=/hff/host.txt bs=4096 count=10240000 oflag=direct
> 10240000+0 records in
> 10240000+0 records out
> 41943040000 bytes (39.1GB) copied, 347.990191 seconds, 114.9MB/s
>
> 
>
> [root@rqy-k8s-3 kbuser]# ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/busybox:latest hfftest sh
> / # mkdir /hff
> / # dd if=/dev/zero of=/hff/host.txt bs=4096 count=10240000 oflag=direct
> 10240000+0 records in
> 10240000+0 records out
> 41943040000 bytes (39.1GB) copied, 312.486063 seconds, 128.0MB/s

## fio

### 随机读写：

> echo 3 > /proc/sys/vm/drop_caches
>
> ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/hff-iperf3-fio:v0.1  hfftest sh
>
>  随机读： 
>
> fio --filename=/tmp/test -direct=1 -iodepth 1 -thread -rw=randread  -ioengine=libaio -bs=512k  -numjobs=8 --size=1G -group_reporting -name=randread
>
> 顺序读： 
>
> fio --filename=/tmp/test --direct=1 --iodepth 1 --thread --rw=read --ioengine=psync --bs=512k --size=1G --numjobs=8 --group_reporting --name=mytest
>
> 随机写： 
>
> fio --filename=/tmp/test --direct=1 --iodepth 1 --thread --rw=randwrite --ioengine=psync --bs=512k --size=1G --numjobs=8 --group_reporting --name=mytest 
>
> 顺序写： 
>
> fio --filename=/tmp/test --direct=1 --iodepth 1 --thread --rw=write --ioengine=psync --bs=512k --size=1G --numjobs=8 --group_reporting --name=mytest 

rqy-k8s-3: --direct=1

|        | 随机读1G | 顺序读1G | 随机写1G          | 顺序写1G         |
| ------ | -------- | -------- | ----------------- | ---------------- |
| kata   | 8600     | 9389     | 741(370.9MB/s)    | 756(378.4MB/s)   |
| runc   | 647      | 1464     | 13159(6579.1MB/s) | 6387(3193.8MB/s) |
| 物理机 | 1681     | 4477     | 10850(5425.2MB/s) | 9666(4833.4MB/s) |

> kata:
>
> read : io=8192.0MB, bw=4300.3MB/s, iops=8600, runt=  1905msec
>
> read : io=8192.0MB, bw=4694.6MB/s, iops=9389, runt=  1745msec
>
> write: io=8192.0MB, bw=379850KB/s, iops=741, runt= 22084msec
>
> write: io=8192.0MB, bw=387518KB/s, iops=756, runt= 21647msec
>
> runc: 
>
> read : io=8192.0MB, bw=331762KB/s, iops=647, runt= 25285msec
>
> read : io=8192.0MB, bw=749786KB/s, iops=1464, runt= 11188msec
>
> write: io=8192.0MB, bw=6579.1MB/s, iops=13159, runt=  1245msec
>
> write: io=8192.0MB, bw=3193.8MB/s, iops=6387, runt=  2565msec
>
>  host:
>
> read : io=8192.0MB, bw=860900KB/s, iops=1681, runt=  9744msec
>
> read : io=8192.0MB, bw=2238.9MB/s, iops=4477, runt=  3659msec
>
> write: io=8192.0MB, bw=5425.2MB/s, iops=10850, runt=  1510msec
>
> write: io=8192.0MB, bw=4833.4MB/s, iops=9666, runt=  1695msec

### 大小文件随机读测试

小文件随机读

fio --name=small-file-multi-read --directory=/ --rw=randread --file_service_type=sequential --bs=4k --filesize=10M --nrfiles=100 --runtime=60 --time_based --numjobs=1 

大文件随机读

fio --name=5G-bigfile-rand-read --directory=/  --rw=randread --size=5G --bs=4k  --runtime=60 --time_based --numjobs=1 

root: rqy-k8s-3
| IOPS   | 小文件10M | 大文件5G | 大文件10G |
| ------ | --------- | -------- | -------- |
| kata   | 766021 | 35584    | 38707 |
| runc   | 136547    | 316      | 267    |
| 物理机 | 93862     | 300      | 284    |

nfs:rqy-k8s-3:

|        |    小文件10M    |   大文件5G   |
| ------ | ------ | ------ |
| kata   | 733795 | 35251  |
| runc   | 436715 | 13729  |

--direct=1:

| IOPS   | 小文件10M | 大文件5G | 大文件10G |
| ------ | --------- | -------- | --------- |
| kata   | 41716     |   39963  |      39312     |
| runc   |  22031   |    305    |    276      |
| 物理机 |      21828     |     307     |   231      |



## sysbench

> 生成测试文件；sysbench --test=fileio --file-num=10 --file-total-size=5G prepare 表示生成10个5G的文件
>
> 运行测试： sysbench --test=fileio --file-total-size=5G --file-test-mode=rndrw --max-requests=5000 --num-threads=16 --file-num=10 --file-extra-flags=direct --file-fsync-freq=0 --file-block-size=16384 run
>
> sysbench --test=fileio cleanup
>
> ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/dotnetdr/sysbench:0.5 hfftest sh



seqwr：顺序写

seqrewr：顺序重写

seqrd：顺序读

rndrd：随机读取

rndwr：随机写入

rndrw：混合随机读/写

10.208.11.110:

|      | rndrw：混合随机读/写                                         |
| ---- | ------------------------------------------------------------ |
| runc | 5368709120 bytes written in 44.82 seconds (114.23 MB/sec).<br>Total transferred 78.125Mb  (276.22Mb/sec)<br>total time:                          0.2828s |
| kata | 5368709120 bytes written in 106.52 seconds (48.07 MB/sec).<br>Total transferred 78.125Mb  (29.944Mb/sec)<br/>total time:                          2.6090s |
| 主机 | 5368709120 bytes written in 40.07 seconds (127.77 MiB/sec)<br>total time:                          0.2545s<br/>read, MiB/s:                  182.88<br/>written, MiB/s:               122.64 |

> host:
>
> 5368709120 bytes written in 40.07 seconds (127.77 MiB/sec).
> [root@localhost ~]# sysbench --test=fileio --file-total-size=5G --file-test-mode=rndrw --max-requests=5000 --num-threads=16 --file-num=10 --file-extra-flags=direct --file-fsync-freq=0 --file-block-size=16384 run
> WARNING: the --test option is deprecated. You can pass a script name or path on the command line without any options.
> WARNING: --num-threads is deprecated, use --threads instead
> WARNING: --max-requests is deprecated, use --events instead
> sysbench 1.0.20 (using bundled LuaJIT 2.1.0-beta2)
>
> Running the test with following options:
> Number of threads: 16
> Initializing random number generator from current time
>
>
> Extra file open flags: directio
> 10 files, 512MiB each
> 5GiB total file size
> Block size 16KiB
> Number of IO requests: 5000
> Read/Write ratio for combined random IO test: 1.50
> Calling fsync() at the end of test, Enabled.
> Using synchronous I/O mode
> Doing random r/w test
> Initializing worker threads...
>
> Threads started!
>
>
> File operations:
>  reads/s:                      11704.55
>  writes/s:                     7848.66
>  fsyncs/s:                     625.70
>
> Throughput:
>  read, MiB/s:                  182.88
>  written, MiB/s:               122.64
>
> General statistics:
>  total time:                          0.2545s
>  total number of events:              5000
>
> Latency (ms):
>       min:                                    0.05
>       avg:                                    0.80
>       max:                                    3.60
>       95th percentile:                        1.64
>       sum:                                 4013.26
>
> Threads fairness:
>  events (avg/stddev):           312.5000/8.49
>  execution time (avg/stddev):   0.2508/0.00
>
> 
>
> kata:
>
> 5368709120 bytes written in 106.52 seconds (48.07 MB/sec).
> sh-4.2#  sysbench --test=fileio --file-total-size=5G --file-test-mode=rndrw --max-requests=5000 --num-threads=16 --file-num=10 --file-extra-flags=direct --file-fsync-freq=0 --file-block-size=16384 run
> sysbench 0.5:  multi-threaded system evaluation benchmark
>
> Running the test with following options:
> Number of threads: 16
> Random number generator seed is 0 and will be ignored
>
>
> Extra file open flags: 3
> 10 files, 512Mb each
> 5Gb total file size
> Block size 16Kb
> Number of IO requests: 5000
> Read/Write ratio for combined random IO test: 1.50
> Calling fsync() at the end of test, Enabled.
> Using synchronous I/O mode
> Doing random r/w test
> Threads started!
>
> Operations performed:  2996 reads, 2004 writes, 10 Other = 5010 Total
> Read 46.812Mb  Written 31.312Mb  Total transferred 78.125Mb  (29.944Mb/sec)
> 1916.43 Requests/sec executed
>
> General statistics:
>  total time:                          2.6090s
>  total number of events:              5000
>  total time taken by event execution: 39.8598s
>  response time:
>       min:                                  0.02ms
>       avg:                                  7.97ms
>       max:                                140.80ms
>       approx.  95 percentile:              19.61ms
>
> Threads fairness:
>  events (avg/stddev):           312.5000/26.30
>  execution time (avg/stddev):   2.4912/0.05
>
> 
>
> runc:
>
> 5368709120 bytes written in 44.82 seconds (114.23 MB/sec).
> sh-4.2# sysbench --test=fileio --file-total-size=5G --file-test-mode=rndrw --max-requests=5000 --num-threads=16 --file-num=10 --file-extra-flags=direct --file-fsync-freq=0 --file-block-size=16384 run
> sysbench 0.5:  multi-threaded system evaluation benchmark
>
> Running the test with following options:
> Number of threads: 16
> Random number generator seed is 0 and will be ignored
>
>
> Extra file open flags: 3
> 10 files, 512Mb each
> 5Gb total file size
> Block size 16Kb
> Number of IO requests: 5000
> Read/Write ratio for combined random IO test: 1.50
> Calling fsync() at the end of test, Enabled.
> Using synchronous I/O mode
> Doing random r/w test
> Threads started!
>
> Operations performed:  3002 reads, 1998 writes, 10 Other = 5010 Total
> Read 46.906Mb  Written 31.219Mb  Total transferred 78.125Mb  (276.22Mb/sec)
> 17678.09 Requests/sec executed
>
> General statistics:
>  total time:                          0.2828s
>  total number of events:              5000
>  total time taken by event execution: 4.3543s
>  response time:
>       min:                                  0.05ms
>       avg:                                  0.87ms
>       max:                                  5.37ms
>       approx.  95 percentile:               1.82ms
>
> Threads fairness:
>  events (avg/stddev):           312.5000/7.57
>  execution time (avg/stddev):   0.2721/0.00
>
> 
