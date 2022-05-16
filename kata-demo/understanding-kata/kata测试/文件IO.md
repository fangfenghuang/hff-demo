[TOC]

裸机 vs runc容器 vs kata容器磁盘IO对比

# 测试环境与配置：

说明：
- 限制容器request/limit 1C2G
- kata设置debug_console_enabled=true（虚拟机开销占用业务开销）
- kata设置debug_console_enabled=false（虚拟机开销不限制）



# fio 
使用fio 对磁盘做性能测试：(4k块大小1G文件)
```bash
[global]
filename=/test/test.file
ioengine=libaio
iodepth=1
thread
group_reporting
direct=1 
bs=4K
size=1G
numjobs=8
runtime=60

[seq-read-4k-1G] 
rw=read

[random-write-4k-1G] 
rw=randwrite

[seq-write-4k-1G] 
rw=write

[rwmix-50r-50w-4k-1G] 
rw=randrw
rwmixread=50
```

```bash
fio jobfile.fio
```

|  iops    |host|runc|kata（true）
|----------|-----------|-------------|-------------|
|顺序读     | 80396(80.4k) | 55.9k | 55.2k
|随机写     | 79278(79.3k) | 60.6k | 49.3k
|顺序写     |  96691(96.7k) | 46.3k | 50.5k
|混合随机读写| 40594/40530(40.5k) | 22.7k/22.7k | 21.4k/21.4k


## 测试数据
### host
```bash
[root@telecom-k8s-phy01 test]# fio jobfile.fio
seq-read-4k-1G: (g=0): rw=read, bs=4K-4K/4K-4K/4K-4K, ioengine=libaio, iodepth=1
...
random-write-4k-1G: (g=1): rw=randwrite, bs=4K-4K/4K-4K/4K-4K, ioengine=libaio, iodepth=1
...
seq-write-4k-1G: (g=2): rw=write, bs=4K-4K/4K-4K/4K-4K, ioengine=libaio, iodepth=1
...
rwmix-50r-50w-4k-1G: (g=3): rw=randrw, bs=4K-4K/4K-4K/4K-4K, ioengine=libaio, iodepth=1
...
fio-2.2.10
Starting 32 threads
seq-read-4k-1G: Laying out IO file(s) (1 file(s) / 1024MB)
Jobs: 8 (f=8): [_(24),m(8)] [100.0% done] [162.4MB/159.4MB/0KB /s] [41.6K/40.8K/0 iops] [eta 00m:00s]
seq-read-4k-1G: (groupid=0, jobs=8): err= 0: pid=26216: Mon May 16 10:03:13 2022
  read : io=8192.0MB, bw=321587KB/s, iops=80396, runt= 26085msec
    slat (usec): min=3, max=958, avg= 8.23, stdev= 3.31
    clat (usec): min=1, max=60166, avg=90.65, stdev=485.64
     lat (usec): min=29, max=60179, avg=98.97, stdev=485.62
    clat percentiles (usec):
     |  1.00th=[   60],  5.00th=[   66], 10.00th=[   68], 20.00th=[   71],
     | 30.00th=[   72], 40.00th=[   74], 50.00th=[   77], 60.00th=[   87],
     | 70.00th=[   91], 80.00th=[   95], 90.00th=[  101], 95.00th=[  104],
     | 99.00th=[  126], 99.50th=[  135], 99.90th=[  342], 99.95th=[ 2480],
     | 99.99th=[31616]
    bw (KB  /s): min=13672, max=45112, per=12.50%, avg=40196.34, stdev=7555.91
    lat (usec) : 2=0.01%, 4=0.01%, 10=0.01%, 20=0.01%, 50=0.21%
    lat (usec) : 100=85.42%, 250=14.23%, 500=0.06%, 750=0.01%, 1000=0.01%
    lat (msec) : 2=0.01%, 4=0.01%, 10=0.02%, 20=0.01%, 50=0.01%
    lat (msec) : 100=0.01%
  cpu          : usr=2.62%, sys=12.04%, ctx=2097396, majf=0, minf=314
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued    : total=r=2097152/w=0/d=0, short=r=0/w=0/d=0, drop=r=0/w=0/d=0
     latency   : target=0, window=0, percentile=100.00%, depth=1
random-write-4k-1G: (groupid=1, jobs=8): err= 0: pid=27722: Mon May 16 10:03:13 2022
  write: io=8192.0MB, bw=317114KB/s, iops=79278, runt= 26453msec
    slat (usec): min=4, max=39819, avg= 7.85, stdev=27.68
    clat (usec): min=0, max=126293, avg=92.25, stdev=456.20
     lat (usec): min=28, max=126302, avg=100.19, stdev=457.15
    clat percentiles (usec):
     |  1.00th=[   69],  5.00th=[   76], 10.00th=[   78], 20.00th=[   80],
     | 30.00th=[   81], 40.00th=[   82], 50.00th=[   83], 60.00th=[   84],
     | 70.00th=[   84], 80.00th=[   85], 90.00th=[   89], 95.00th=[   98],
     | 99.00th=[  145], 99.50th=[  167], 99.90th=[ 2096], 99.95th=[ 2640],
     | 99.99th=[11456]
    bw (KB  /s): min=  347, max=44752, per=12.49%, avg=39606.35, stdev=11008.55
    lat (usec) : 2=0.01%, 4=0.01%, 10=0.01%, 20=0.01%, 50=0.23%
    lat (usec) : 100=95.51%, 250=3.99%, 500=0.01%, 750=0.01%, 1000=0.02%
    lat (msec) : 2=0.11%, 4=0.08%, 10=0.01%, 20=0.01%, 50=0.01%
    lat (msec) : 100=0.01%, 250=0.01%
  cpu          : usr=2.88%, sys=10.93%, ctx=2104963, majf=0, minf=490
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued    : total=r=0/w=2097152/d=0, short=r=0/w=0/d=0, drop=r=0/w=0/d=0
     latency   : target=0, window=0, percentile=100.00%, depth=1
seq-write-4k-1G: (groupid=2, jobs=8): err= 0: pid=29631: Mon May 16 10:03:13 2022
  write: io=8192.0MB, bw=386768KB/s, iops=96691, runt= 21689msec
    slat (usec): min=4, max=1392, avg= 7.07, stdev= 2.51
    clat (usec): min=1, max=35541, avg=75.00, stdev=40.70
     lat (usec): min=28, max=35562, avg=82.15, stdev=40.68
    clat percentiles (usec):
     |  1.00th=[   64],  5.00th=[   69], 10.00th=[   70], 20.00th=[   72],
     | 30.00th=[   73], 40.00th=[   74], 50.00th=[   74], 60.00th=[   75],
     | 70.00th=[   75], 80.00th=[   76], 90.00th=[   78], 95.00th=[   81],
     | 99.00th=[  119], 99.50th=[  135], 99.90th=[  153], 99.95th=[  159],
     | 99.99th=[  205]
    bw (KB  /s): min=44574, max=49264, per=12.51%, avg=48377.10, stdev=509.36
    lat (usec) : 2=0.01%, 4=0.01%, 10=0.01%, 20=0.01%, 50=0.15%
    lat (usec) : 100=98.53%, 250=1.32%, 500=0.01%, 750=0.01%, 1000=0.01%
    lat (msec) : 2=0.01%, 4=0.01%, 10=0.01%, 20=0.01%, 50=0.01%
  cpu          : usr=3.43%, sys=11.82%, ctx=2104891, majf=0, minf=238
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued    : total=r=0/w=2097152/d=0, short=r=0/w=0/d=0, drop=r=0/w=0/d=0
     latency   : target=0, window=0, percentile=100.00%, depth=1
rwmix-50r-50w-4k-1G: (groupid=3, jobs=8): err= 0: pid=30632: Mon May 16 10:03:13 2022
  read : io=4099.3MB, bw=162378KB/s, iops=40594, runt= 25851msec
    slat (usec): min=3, max=1060, avg= 7.79, stdev= 3.48
    clat (usec): min=1, max=3922, avg=97.74, stdev=20.14
     lat (usec): min=28, max=3933, avg=105.61, stdev=19.83
    clat percentiles (usec):
     |  1.00th=[   68],  5.00th=[   74], 10.00th=[   77], 20.00th=[   81],
     | 30.00th=[   85], 40.00th=[   93], 50.00th=[  100], 60.00th=[  104],
     | 70.00th=[  107], 80.00th=[  111], 90.00th=[  116], 95.00th=[  122],
     | 99.00th=[  153], 99.50th=[  169], 99.90th=[  211], 99.95th=[  225],
     | 99.99th=[  294]
    bw (KB  /s): min=19184, max=21456, per=12.49%, avg=20287.92, stdev=461.09
  write: io=4092.8MB, bw=162120KB/s, iops=40530, runt= 25851msec
    slat (usec): min=3, max=995, avg= 7.94, stdev= 3.55
    clat (usec): min=1, max=7102, avg=82.09, stdev=18.03
     lat (usec): min=27, max=7107, avg=90.11, stdev=18.00
    clat percentiles (usec):
     |  1.00th=[   63],  5.00th=[   71], 10.00th=[   73], 20.00th=[   76],
     | 30.00th=[   78], 40.00th=[   79], 50.00th=[   81], 60.00th=[   82],
     | 70.00th=[   83], 80.00th=[   85], 90.00th=[   90], 95.00th=[   97],
     | 99.00th=[  131], 99.50th=[  147], 99.90th=[  189], 99.95th=[  205],
     | 99.99th=[  251]
    bw (KB  /s): min=19384, max=21072, per=12.50%, avg=20260.06, stdev=352.92
    lat (usec) : 2=0.01%, 4=0.01%, 10=0.01%, 20=0.01%, 50=0.11%
    lat (usec) : 100=72.62%, 250=27.25%, 500=0.01%, 750=0.01%, 1000=0.01%
    lat (msec) : 2=0.01%, 4=0.01%, 10=0.01%
  cpu          : usr=3.00%, sys=11.25%, ctx=2111174, majf=0, minf=480
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued    : total=r=1049408/w=1047744/d=0, short=r=0/w=0/d=0, drop=r=0/w=0/d=0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
   READ: io=8192.0MB, aggrb=321587KB/s, minb=321587KB/s, maxb=321587KB/s, mint=26085msec, maxt=26085msec

Run status group 1 (all jobs):
  WRITE: io=8192.0MB, aggrb=317113KB/s, minb=317113KB/s, maxb=317113KB/s, mint=26453msec, maxt=26453msec

Run status group 2 (all jobs):
  WRITE: io=8192.0MB, aggrb=386767KB/s, minb=386767KB/s, maxb=386767KB/s, mint=21689msec, maxt=21689msec

Run status group 3 (all jobs):
   READ: io=4099.3MB, aggrb=162377KB/s, minb=162377KB/s, maxb=162377KB/s, mint=25851msec, maxt=25851msec
  WRITE: io=4092.8MB, aggrb=162120KB/s, minb=162120KB/s, maxb=162120KB/s, mint=25851msec, maxt=25851msec

Disk stats (read/write):
    dm-0: ios=3137975/5237442, merge=0/0, ticks=289341/431921, in_queue=725766, util=99.67%, aggrios=3146559/5246382, aggrmerge=0/66, aggrticks=290576/435171, aggrin_queue=724207, aggrutil=99.22%
  sda: ios=3146559/5246382, merge=0/66, ticks=290576/435171, in_queue=724207, util=99.22%


```

### runc
```bash
/test # fio jobfile.fio
seq-read-4k-1G: (g=0): rw=read, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=1
...
random-write-4k-1G: (g=1): rw=randwrite, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=1
...
seq-write-4k-1G: (g=2): rw=write, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=1
...
rwmix-50r-50w-4k-1G: (g=3): rw=randrw, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=1
...
fio-3.13
Starting 32 threads
seq-read-4k-1G: Laying out IO file (1 file / 1024MiB)
Jobs: 8 (f=8): [_(24),m(8)][99.4%][r=101MiB/s,w=103MiB/s][r=25.0k,w=26.3k IOPS][eta 00m:01s]
seq-read-4k-1G: (groupid=0, jobs=8): err= 0: pid=222: Mon May 16 02:16:49 2022
  read: IOPS=55.9k, BW=218MiB/s (229MB/s)(8192MiB/37538msec)
    slat (usec): min=3, max=42140, avg=10.35, stdev=105.15
    clat (nsec): min=937, max=91316k, avg=132216.73, stdev=1304684.96
     lat (usec): min=28, max=91326, avg=142.65, stdev=1308.85
    clat percentiles (usec):
     |  1.00th=[   44],  5.00th=[   62], 10.00th=[   67], 20.00th=[   71],
     | 30.00th=[   74], 40.00th=[   79], 50.00th=[   85], 60.00th=[   89],
     | 70.00th=[   93], 80.00th=[   97], 90.00th=[  103], 95.00th=[  113],
     | 99.00th=[  141], 99.50th=[  159], 99.90th=[31851], 99.95th=[37487],
     | 99.99th=[41681]
   bw (  KiB/s): min=133468, max=229296, per=74.61%, avg=166725.45, stdev=2611.31, samples=592
   iops        : min=33365, max=57320, avg=41678.24, stdev=652.81, samples=592
  lat (nsec)   : 1000=0.01%
  lat (usec)   : 2=0.01%, 4=0.01%, 10=0.01%, 20=0.01%, 50=1.63%
  lat (usec)   : 100=83.37%, 250=14.81%, 500=0.02%, 750=0.01%, 1000=0.01%
  lat (msec)   : 2=0.01%, 4=0.01%, 10=0.02%, 20=0.01%, 50=0.12%
  lat (msec)   : 100=0.01%
  cpu          : usr=1.64%, sys=10.69%, ctx=2097301, majf=0, minf=411
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=2097152,0,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1
random-write-4k-1G: (groupid=1, jobs=8): err= 0: pid=230: Mon May 16 02:16:49 2022
  write: IOPS=60.6k, BW=237MiB/s (248MB/s)(8192MiB/34591msec); 0 zone resets
    slat (usec): min=3, max=37045, avg= 9.00, stdev=89.92
    clat (nsec): min=965, max=50633k, avg=121966.97, stdev=957260.48
     lat (usec): min=28, max=50637, avg=131.06, stdev=961.41
    clat percentiles (usec):
     |  1.00th=[   47],  5.00th=[   73], 10.00th=[   77], 20.00th=[   80],
     | 30.00th=[   81], 40.00th=[   82], 50.00th=[   83], 60.00th=[   84],
     | 70.00th=[   85], 80.00th=[   86], 90.00th=[   91], 95.00th=[  103],
     | 99.00th=[  157], 99.50th=[  192], 99.90th=[22938], 99.95th=[27919],
     | 99.99th=[32900]
   bw (  KiB/s): min=49414, max=259189, per=77.81%, avg=188701.75, stdev=4542.88, samples=544
   iops        : min=12349, max=64794, avg=47172.47, stdev=1135.72, samples=544
  lat (nsec)   : 1000=0.01%
  lat (usec)   : 2=0.01%, 4=0.01%, 10=0.01%, 20=0.01%, 50=1.23%
  lat (usec)   : 100=92.82%, 250=5.56%, 500=0.02%, 750=0.02%, 1000=0.02%
  lat (msec)   : 2=0.11%, 4=0.06%, 10=0.04%, 20=0.01%, 50=0.11%
  lat (msec)   : 100=0.01%
  cpu          : usr=1.93%, sys=9.84%, ctx=2107324, majf=0, minf=1072
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,2097152,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1
seq-write-4k-1G: (groupid=2, jobs=8): err= 0: pid=238: Mon May 16 02:16:49 2022
  write: IOPS=46.3k, BW=181MiB/s (190MB/s)(8192MiB/45259msec); 0 zone resets
    slat (usec): min=3, max=57303, avg=16.17, stdev=283.17
    clat (nsec): min=746, max=61181k, avg=153902.91, stdev=1964543.40
     lat (usec): min=27, max=61189, avg=170.30, stdev=1987.10
    clat percentiles (usec):
     |  1.00th=[   35],  5.00th=[   49], 10.00th=[   57], 20.00th=[   63],
     | 30.00th=[   67], 40.00th=[   69], 50.00th=[   71], 60.00th=[   73],
     | 70.00th=[   76], 80.00th=[   79], 90.00th=[   87], 95.00th=[   96],
     | 99.00th=[  186], 99.50th=[  334], 99.90th=[47449], 99.95th=[51643],
     | 99.99th=[55313]
   bw (  KiB/s): min=37898, max=236720, per=94.47%, avg=175087.37, stdev=3238.06, samples=714
   iops        : min= 9473, max=59178, avg=43768.87, stdev=809.50, samples=714
  lat (nsec)   : 750=0.01%, 1000=0.01%
  lat (usec)   : 2=0.07%, 4=0.01%, 10=0.01%, 20=0.01%, 50=5.37%
  lat (usec)   : 100=90.19%, 250=3.68%, 500=0.27%, 750=0.06%, 1000=0.07%
  lat (msec)   : 2=0.07%, 4=0.02%, 10=0.01%, 20=0.01%, 50=0.09%
  lat (msec)   : 100=0.07%
  cpu          : usr=1.79%, sys=10.59%, ctx=2108770, majf=0, minf=442
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,2097152,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1
rwmix-50r-50w-4k-1G: (groupid=3, jobs=8): err= 0: pid=246: Mon May 16 02:16:49 2022
  read: IOPS=22.7k, BW=88.7MiB/s (92.0MB/s)(4094MiB/46175msec)
    slat (usec): min=3, max=53317, avg=14.58, stdev=278.36
    clat (nsec): min=818, max=56805k, avg=165871.71, stdev=1727900.11
     lat (usec): min=28, max=56823, avg=180.61, stdev=1750.60
    clat percentiles (usec):
     |  1.00th=[   39],  5.00th=[   62], 10.00th=[   69], 20.00th=[   76],
     | 30.00th=[   82], 40.00th=[   88], 50.00th=[   93], 60.00th=[   98],
     | 70.00th=[  103], 80.00th=[  109], 90.00th=[  119], 95.00th=[  133],
     | 99.00th=[  225], 99.50th=[  322], 99.90th=[41157], 99.95th=[45351],
     | 99.99th=[50070]
   bw (  KiB/s): min=67920, max=133040, per=100.00%, avg=90899.25, stdev=1392.25, samples=732
   iops        : min=16980, max=33260, avg=22724.74, stdev=348.07, samples=732
  write: IOPS=22.7k, BW=88.8MiB/s (93.1MB/s)(4098MiB/46175msec); 0 zone resets
    slat (usec): min=3, max=52129, avg=14.59, stdev=265.02
    clat (nsec): min=783, max=56463k, avg=152802.40, stdev=1742387.77
     lat (usec): min=27, max=56697, avg=167.59, stdev=1763.50
    clat percentiles (usec):
     |  1.00th=[   38],  5.00th=[   55], 10.00th=[   62], 20.00th=[   69],
     | 30.00th=[   73], 40.00th=[   76], 50.00th=[   78], 60.00th=[   81],
     | 70.00th=[   84], 80.00th=[   88], 90.00th=[   97], 95.00th=[  112],
     | 99.00th=[  204], 99.50th=[  297], 99.90th=[41681], 99.95th=[45351],
     | 99.99th=[50070]
   bw (  KiB/s): min=67615, max=130688, per=100.00%, avg=90986.98, stdev=1376.89, samples=732
   iops        : min=16903, max=32672, avg=22746.70, stdev=344.23, samples=732
  lat (nsec)   : 1000=0.01%
  lat (usec)   : 2=0.06%, 4=0.01%, 10=0.01%, 20=0.01%, 50=2.60%
  lat (usec)   : 100=75.12%, 250=21.48%, 500=0.42%, 750=0.07%, 1000=0.03%
  lat (msec)   : 2=0.02%, 4=0.01%, 10=0.01%, 20=0.01%, 50=0.16%
  lat (msec)   : 100=0.01%
  cpu          : usr=1.93%, sys=10.55%, ctx=2114234, majf=0, minf=539
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=1047959,1049193,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
   READ: bw=218MiB/s (229MB/s), 218MiB/s-218MiB/s (229MB/s-229MB/s), io=8192MiB (8590MB), run=37538-37538msec

Run status group 1 (all jobs):
  WRITE: bw=237MiB/s (248MB/s), 237MiB/s-237MiB/s (248MB/s-248MB/s), io=8192MiB (8590MB), run=34591-34591msec

Run status group 2 (all jobs):
  WRITE: bw=181MiB/s (190MB/s), 181MiB/s-181MiB/s (190MB/s-190MB/s), io=8192MiB (8590MB), run=45259-45259msec

Run status group 3 (all jobs):
   READ: bw=88.7MiB/s (92.0MB/s), 88.7MiB/s-88.7MiB/s (92.0MB/s-92.0MB/s), io=4094MiB (4292MB), run=46175-46175msec
  WRITE: bw=88.8MiB/s (93.1MB/s), 88.8MiB/s-88.8MiB/s (93.1MB/s-93.1MB/s), io=4098MiB (4297MB), run=46175-46175msec

Disk stats (read/write):
    dm-0: ios=3142171/5248417, merge=0/0, ticks=281861/420565, in_queue=706662, util=64.01%, aggrios=3145111/5251206, aggrmerge=0/151, aggrticks=283023/423518, aggrin_queue=704727, aggrutil=63.46%
  sda: ios=3145111/5251206, merge=0/151, ticks=283023/423518, in_queue=704727, util=63.46%


```

### kata
```bash
/test # fio jobfile.fio
seq-read-4k-1G: (g=0): rw=read, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=1
...
random-write-4k-1G: (g=1): rw=randwrite, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=1
...
seq-write-4k-1G: (g=2): rw=write, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=1
...
rwmix-50r-50w-4k-1G: (g=3): rw=randrw, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=1
...
fio-3.13
Starting 32 threads
seq-read-4k-1G: Laying out IO file (1 file / 1024MiB)
Jobs: 6 (f=5): [_(24),m(2),_(1),m(2),_(1),m(1),f(1)][99.4%][r=94.2MiB/s,w=94.6MiB/s][r=24.1k,w=24.2k IOPS][eta 00m:01s]
seq-read-4k-1G: (groupid=0, jobs=8): err= 0: pid=10: Mon May 16 02:12:57 2022
  read: IOPS=55.2k, BW=216MiB/s (226MB/s)(8192MiB/37993msec)
    slat (nsec): min=1143, max=151112k, avg=16299.62, stdev=582336.19
    clat (nsec): min=278, max=113995k, avg=125933.02, stdev=1751379.58
     lat (usec): min=15, max=151113, avg=142.61, stdev=1848.62
    clat percentiles (nsec):
     |  1.00th=[     446],  5.00th=[   28288], 10.00th=[   34560],
     | 20.00th=[   41216], 30.00th=[   46336], 40.00th=[   50944],
     | 50.00th=[   55552], 60.00th=[   60672], 70.00th=[   67072],
     | 80.00th=[   74240], 90.00th=[   88576], 95.00th=[  101888],
     | 99.00th=[  140288], 99.50th=[  191488], 99.90th=[42205184],
     | 99.95th=[51118080], 99.99th=[51118080]
   bw (  KiB/s): min=121445, max=286922, per=91.66%, avg=202373.07, stdev=4092.28, samples=597
   iops        : min=30358, max=71727, avg=50590.33, stdev=1023.06, samples=597
  lat (nsec)   : 500=1.34%, 750=0.67%, 1000=0.13%
  lat (usec)   : 2=0.02%, 4=0.01%, 10=0.03%, 20=0.38%, 50=35.09%
  lat (usec)   : 100=56.83%, 250=5.14%, 500=0.13%, 750=0.04%, 1000=0.01%
  lat (msec)   : 2=0.02%, 4=0.01%, 10=0.01%, 20=0.02%, 50=0.05%
  lat (msec)   : 100=0.08%, 250=0.01%
  cpu          : usr=1.39%, sys=10.19%, ctx=2161974, majf=0, minf=13
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=2097152,0,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1
random-write-4k-1G: (groupid=1, jobs=8): err= 0: pid=18: Mon May 16 02:12:57 2022
  write: IOPS=49.3k, BW=193MiB/s (202MB/s)(8192MiB/42522msec); 0 zone resets
    slat (nsec): min=1388, max=151075k, avg=33013.00, stdev=816788.46
    clat (nsec): min=281, max=113246k, avg=124904.46, stdev=1660676.13
     lat (usec): min=16, max=151077, avg=158.40, stdev=1853.42
    clat percentiles (nsec):
     |  1.00th=[     330],  5.00th=[     486], 10.00th=[   28800],
     | 20.00th=[   39680], 30.00th=[   45824], 40.00th=[   50944],
     | 50.00th=[   56576], 60.00th=[   61696], 70.00th=[   68096],
     | 80.00th=[   76288], 90.00th=[   90624], 95.00th=[  104960],
     | 99.00th=[  154624], 99.50th=[  234496], 99.90th=[34865152],
     | 99.95th=[50069504], 99.99th=[51118080]
   bw (  KiB/s): min=89578, max=215358, per=68.80%, avg=135729.26, stdev=2807.15, samples=666
   iops        : min=22391, max=53836, avg=33929.39, stdev=701.79, samples=666
  lat (nsec)   : 500=5.14%, 750=0.79%, 1000=0.17%
  lat (usec)   : 2=0.03%, 4=0.03%, 10=0.27%, 20=0.89%, 50=30.45%
  lat (usec)   : 100=55.91%, 250=5.86%, 500=0.18%, 750=0.05%, 1000=0.02%
  lat (msec)   : 2=0.03%, 4=0.01%, 10=0.01%, 20=0.04%, 50=0.08%
  lat (msec)   : 100=0.06%, 250=0.01%
  cpu          : usr=1.95%, sys=10.14%, ctx=2483253, majf=0, minf=8
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,2097152,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1
seq-write-4k-1G: (groupid=2, jobs=8): err= 0: pid=26: Mon May 16 02:12:57 2022
  write: IOPS=50.5k, BW=197MiB/s (207MB/s)(8192MiB/41526msec); 0 zone resets
    slat (nsec): min=1384, max=154953k, avg=31843.80, stdev=787291.45
    clat (nsec): min=277, max=114869k, avg=121220.53, stdev=1576576.18
     lat (usec): min=15, max=154955, avg=153.54, stdev=1763.48
    clat percentiles (nsec):
     |  1.00th=[     322],  5.00th=[     462], 10.00th=[   24960],
     | 20.00th=[   35072], 30.00th=[   41216], 40.00th=[   46848],
     | 50.00th=[   52480], 60.00th=[   58624], 70.00th=[   65280],
     | 80.00th=[   74240], 90.00th=[   90624], 95.00th=[  109056],
     | 99.00th=[  216064], 99.50th=[  366592], 99.90th=[31850496],
     | 99.95th=[50069504], 99.99th=[51118080]
   bw (  KiB/s): min=43983, max=259537, per=74.72%, avg=150937.58, stdev=4313.70, samples=646
   iops        : min=10993, max=64880, avg=37731.33, stdev=1078.40, samples=646
  lat (nsec)   : 500=5.54%, 750=0.47%, 1000=0.27%
  lat (usec)   : 2=0.04%, 4=0.04%, 10=0.32%, 20=1.28%, 50=38.07%
  lat (usec)   : 100=47.10%, 250=6.09%, 500=0.40%, 750=0.06%, 1000=0.04%
  lat (msec)   : 2=0.05%, 4=0.05%, 10=0.03%, 20=0.04%, 50=0.08%
  lat (msec)   : 100=0.05%, 250=0.01%
  cpu          : usr=1.79%, sys=10.09%, ctx=2510723, majf=0, minf=8
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,2097152,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1
rwmix-50r-50w-4k-1G: (groupid=3, jobs=8): err= 0: pid=34: Mon May 16 02:12:57 2022
  read: IOPS=21.4k, BW=83.7MiB/s (87.7MB/s)(4094MiB/48931msec)
    slat (nsec): min=1190, max=112962k, avg=125989.36, stdev=1718031.14
    clat (nsec): min=280, max=112874k, avg=105187.94, stdev=1574177.07
     lat (usec): min=14, max=112977, avg=231.61, stdev=2334.73
    clat percentiles (nsec):
     |  1.00th=[     370],  5.00th=[   19072], 10.00th=[   28288],
     | 20.00th=[   35584], 30.00th=[   40192], 40.00th=[   44800],
     | 50.00th=[   48896], 60.00th=[   52992], 70.00th=[   58112],
     | 80.00th=[   64768], 90.00th=[   76288], 95.00th=[   89600],
     | 99.00th=[  130560], 99.50th=[  177152], 99.90th=[30277632],
     | 99.95th=[50069504], 99.99th=[51118080]
   bw (  KiB/s): min=59535, max=121806, per=100.00%, avg=85988.84, stdev=1551.16, samples=770
   iops        : min=14883, max=30450, avg=21496.91, stdev=387.79, samples=770
  write: IOPS=21.4k, BW=83.8MiB/s (87.8MB/s)(4098MiB/48931msec); 0 zone resets
    slat (nsec): min=1354, max=112841k, avg=22445.97, stdev=713692.80
    clat (nsec): min=277, max=112876k, avg=111316.25, stdev=1680024.45
     lat (usec): min=15, max=112884, avg=134.05, stdev=1827.46
    clat percentiles (nsec):
     |  1.00th=[     350],  5.00th=[   12352], 10.00th=[   27264],
     | 20.00th=[   34560], 30.00th=[   39680], 40.00th=[   43776],
     | 50.00th=[   47872], 60.00th=[   51968], 70.00th=[   57600],
     | 80.00th=[   65280], 90.00th=[   80384], 95.00th=[   97792],
     | 99.00th=[  146432], 99.50th=[  185344], 99.90th=[35913728],
     | 99.95th=[51118080], 99.99th=[51642368]
   bw (  KiB/s): min=60001, max=122381, per=100.00%, avg=86081.42, stdev=1568.73, samples=770
   iops        : min=14999, max=30594, avg=21520.03, stdev=392.18, samples=770
  lat (nsec)   : 500=3.27%, 750=0.47%, 1000=0.36%
  lat (usec)   : 2=0.04%, 4=0.02%, 10=0.31%, 20=1.12%, 50=48.60%
  lat (usec)   : 100=42.10%, 250=3.42%, 500=0.12%, 750=0.02%, 1000=0.01%
  lat (msec)   : 2=0.01%, 4=0.01%, 10=0.01%, 20=0.02%, 50=0.05%
  lat (msec)   : 100=0.07%, 250=0.01%
  cpu          : usr=1.73%, sys=9.98%, ctx=3435262, majf=0, minf=8
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=1047959,1049193,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
   READ: bw=216MiB/s (226MB/s), 216MiB/s-216MiB/s (226MB/s-226MB/s), io=8192MiB (8590MB), run=37993-37993msec

Run status group 1 (all jobs):
  WRITE: bw=193MiB/s (202MB/s), 193MiB/s-193MiB/s (202MB/s-202MB/s), io=8192MiB (8590MB), run=42522-42522msec

Run status group 2 (all jobs):
  WRITE: bw=197MiB/s (207MB/s), 197MiB/s-197MiB/s (207MB/s-207MB/s), io=8192MiB (8590MB), run=41526-41526msec

Run status group 3 (all jobs):
   READ: bw=83.7MiB/s (87.7MB/s), 83.7MiB/s-83.7MiB/s (87.7MB/s-87.7MB/s), io=4094MiB (4292MB), run=48931-48931msec
  WRITE: bw=83.8MiB/s (87.8MB/s), 83.8MiB/s-83.8MiB/s (87.8MB/s-87.8MB/s), io=4098MiB (4297MB), run=48931-48931msec
/test #


```