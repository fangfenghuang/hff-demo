[TOC]

# 测试环境
- limit: 2C2G
- 测试节点10.208.11.110、rqy-k8s-1


# 测试项
## fio

### 随机读写：
```
echo 3 > /proc/sys/vm/drop_caches


随机读： 
fio --filename=/test/randread -direct=1 -iodepth 1 -thread -rw=randread  -ioengine=libaio -bs=512k  -numjobs=8 --size=1G -group_reporting -name=randread
顺序读： 
fio --filename=/test/read --direct=1 --iodepth 1 --thread --rw=read --ioengine=psync --bs=512k --size=1G --numjobs=8 --group_reporting --name=read
随机写： 
fio --filename=/test/randwrite --direct=1 --iodepth 1 --thread --rw=randwrite --ioengine=psync --bs=512k --size=1G --numjobs=8 --group_reporting --name=randwrite 
顺序写： 
fio --filename=/test/write --direct=1 --iodepth 1 --thread --rw=write --ioengine=psync --bs=512k --size=1G --numjobs=8 --group_reporting --name=write 
```


rqy-k8s-1: 

|        | 随机读1G | 顺序读1G | 随机写1G | 顺序写1G | 
| ------ | ------ | ------ | ------- | ------ |
| kata   | 779  |  4114  | 370  | 538 |
| runc   | 845  | 858 | 427 | 417 |



kata
```bash
[root@localhost ~]# kubectl exec -it test-kata sh
/data # fio --filename=/test/randread -direct=1 -iodepth 1 -thread -rw=randread  -ioengine=libaio -bs=512k  -numjobs=8 --size
=1G -group_reporting -name=randread
randread: (g=0): rw=randread, bs=(R) 512KiB-512KiB, (W) 512KiB-512KiB, (T) 512KiB-512KiB, ioengine=libaio, iodepth=1
...
fio-3.13
Starting 8 threads
randread: Laying out IO file (1 file / 1024MiB)
Jobs: 4 (f=2): [_(1),r(1),f(1),_(1),r(1),_(2),f(1)][100.0%][r=507MiB/s][r=1014 IOPS][eta 00m:00s]
randread: (groupid=0, jobs=8): err= 0: pid=10: Thu Apr 21 07:04:04 2022
  read: IOPS=779, BW=390MiB/s (409MB/s)(8192MiB/21014msec)
    slat (usec): min=6, max=58225, avg=3547.23, stdev=3932.50
    clat (usec): min=114, max=62812, avg=6636.91, stdev=4804.31
     lat (usec): min=132, max=81685, avg=10184.80, stdev=6754.24
    clat percentiles (usec):
     |  1.00th=[  404],  5.00th=[ 1696], 10.00th=[ 2376], 20.00th=[ 3359],
     | 30.00th=[ 4146], 40.00th=[ 4883], 50.00th=[ 5604], 60.00th=[ 6390],
     | 70.00th=[ 7439], 80.00th=[ 8717], 90.00th=[11731], 95.00th=[15139],
     | 99.00th=[24773], 99.50th=[31589], 99.90th=[44827], 99.95th=[47973],
     | 99.99th=[56361]
   bw (  KiB/s): min=211914, max=612352, per=99.91%, avg=398811.81, stdev=11257.29, samples=330
   iops        : min=  413, max= 1196, avg=778.74, stdev=22.02, samples=330
  lat (usec)   : 250=0.38%, 500=0.90%, 750=0.67%, 1000=0.13%
  lat (msec)   : 2=5.11%, 4=20.81%, 10=57.45%, 20=12.41%, 50=2.10%
  lat (msec)   : 100=0.04%
  cpu          : usr=0.06%, sys=0.30%, ctx=39194, majf=0, minf=1032
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=16384,0,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
   READ: bw=390MiB/s (409MB/s), 390MiB/s-390MiB/s (409MB/s-409MB/s), io=8192MiB (8590MB), run=21014-21014msec
/data # fio --filename=/test/read --direct=1 --iodepth 1 --thread --rw=read --ioengine=psync --bs=512k --size=1G --numjobs=8
--group_reporting --name=read
read: (g=0): rw=read, bs=(R) 512KiB-512KiB, (W) 512KiB-512KiB, (T) 512KiB-512KiB, ioengine=psync, iodepth=1
...
fio-3.13
Starting 8 threads
read: Laying out IO file (1 file / 1024MiB)
Jobs: 8 (f=8): [R(8)][80.0%][r=1841MiB/s][r=3681 IOPS][eta 00m:01s]
read: (groupid=0, jobs=8): err= 0: pid=20: Thu Apr 21 07:21:10 2022
  read: IOPS=4114, BW=2057MiB/s (2157MB/s)(8192MiB/3982msec)
    clat (usec): min=74, max=25230, avg=1930.07, stdev=1409.08
     lat (usec): min=75, max=25230, avg=1930.23, stdev=1409.08
    clat percentiles (usec):
     |  1.00th=[  204],  5.00th=[  523], 10.00th=[ 1139], 20.00th=[ 1270],
     | 30.00th=[ 1319], 40.00th=[ 1385], 50.00th=[ 1467], 60.00th=[ 1598],
     | 70.00th=[ 2008], 80.00th=[ 2409], 90.00th=[ 3326], 95.00th=[ 4424],
     | 99.00th=[ 7308], 99.50th=[ 9896], 99.90th=[14353], 99.95th=[14877],
     | 99.99th=[22414]
   bw (  MiB/s): min= 1121, max= 2658, per=94.61%, avg=1946.43, stdev=59.63, samples=56
   iops        : min= 2242, max= 5316, avg=3892.57, stdev=119.21, samples=56
  lat (usec)   : 100=0.02%, 250=1.35%, 500=3.45%, 750=2.38%, 1000=1.12%
  lat (msec)   : 2=61.58%, 4=23.94%, 10=5.69%, 20=0.43%, 50=0.03%
  cpu          : usr=0.18%, sys=1.26%, ctx=44470, majf=0, minf=1024
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=16384,0,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
   READ: bw=2057MiB/s (2157MB/s), 2057MiB/s-2057MiB/s (2157MB/s-2157MB/s), io=8192MiB (8590MB), run=3982-3982msec
/data #
/data #
/data #
/data # fio --filename=/test/randwrite --direct=1 --iodepth 1 --thread --rw=randwrite --ioengine=psync --bs=512k --size=1G --
numjobs=8 --group_reporting --name=randwrite
randwrite: (g=0): rw=randwrite, bs=(R) 512KiB-512KiB, (W) 512KiB-512KiB, (T) 512KiB-512KiB, ioengine=psync, iodepth=1
...
fio-3.13
Starting 8 threads
randwrite: Laying out IO file (1 file / 1024MiB)
Jobs: 5 (f=5): [w(2),_(2),w(3),_(1)][97.8%][w=387MiB/s][w=773 IOPS][eta 00m:01s]
randwrite: (groupid=0, jobs=8): err= 0: pid=30: Thu Apr 21 07:22:20 2022
  write: IOPS=370, BW=185MiB/s (194MB/s)(8192MiB/44225msec); 0 zone resets
    clat (usec): min=99, max=456043, avg=21354.62, stdev=34382.05
     lat (usec): min=105, max=456057, avg=21365.49, stdev=34382.10
    clat percentiles (usec):
     |  1.00th=[   135],  5.00th=[   180], 10.00th=[   269], 20.00th=[  1614],
     | 30.00th=[  6783], 40.00th=[ 10028], 50.00th=[ 11338], 60.00th=[ 16057],
     | 70.00th=[ 18744], 80.00th=[ 25822], 90.00th=[ 49546], 95.00th=[ 89654],
     | 99.00th=[160433], 99.50th=[208667], 99.90th=[379585], 99.95th=[417334],
     | 99.99th=[442500]
   bw (  KiB/s): min=14331, max=775501, per=99.60%, avg=188913.72, stdev=22335.67, samples=696
   iops        : min=   27, max= 1513, avg=368.38, stdev=43.62, samples=696
  lat (usec)   : 100=0.01%, 250=9.53%, 500=1.18%, 750=0.20%, 1000=0.48%
  lat (msec)   : 2=12.86%, 4=3.00%, 10=12.57%, 20=35.55%, 50=14.74%
  lat (msec)   : 100=6.44%, 250=3.20%, 500=0.24%
  cpu          : usr=0.08%, sys=0.24%, ctx=35229, majf=0, minf=0
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,16384,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
  WRITE: bw=185MiB/s (194MB/s), 185MiB/s-185MiB/s (194MB/s-194MB/s), io=8192MiB (8590MB), run=44225-44225msec
/data # fio --filename=/test/write --direct=1 --iodepth 1 --thread --rw=write --ioengine=psync --bs=512k --size=1G --numjobs=
8 --group_reporting --name=write
write: (g=0): rw=write, bs=(R) 512KiB-512KiB, (W) 512KiB-512KiB, (T) 512KiB-512KiB, ioengine=psync, iodepth=1
...
fio-3.13
Starting 8 threads
write: Laying out IO file (1 file / 1024MiB)
Jobs: 4 (f=4): [W(2),_(1),W(1),_(2),W(1),_(1)][88.2%][w=501MiB/s][w=1002 IOPS][eta 00m:04s]
write: (groupid=0, jobs=8): err= 0: pid=40: Thu Apr 21 07:23:09 2022
  write: IOPS=538, BW=269MiB/s (282MB/s)(8192MiB/30437msec); 0 zone resets
    clat (usec): min=79, max=659022, avg=14603.08, stdev=39223.62
     lat (usec): min=87, max=659033, avg=14613.37, stdev=39223.86
    clat percentiles (usec):
     |  1.00th=[   113],  5.00th=[   133], 10.00th=[   167], 20.00th=[  1172],
     | 30.00th=[  1303], 40.00th=[  1483], 50.00th=[  1795], 60.00th=[  7832],
     | 70.00th=[ 14353], 80.00th=[ 21103], 90.00th=[ 32900], 95.00th=[ 43254],
     | 99.00th=[210764], 99.50th=[316670], 99.90th=[492831], 99.95th=[583009],
     | 99.99th=[641729]
   bw (  KiB/s): min= 8192, max=1386244, per=98.15%, avg=270500.05, stdev=39219.54, samples=474
   iops        : min=   16, max= 2707, avg=528.02, stdev=76.60, samples=474
  lat (usec)   : 100=0.30%, 250=14.05%, 500=1.54%, 750=1.06%, 1000=0.84%
  lat (msec)   : 2=34.21%, 4=3.68%, 10=6.20%, 20=17.14%, 50=17.18%
  lat (msec)   : 100=2.09%, 250=0.90%, 500=0.72%, 750=0.10%
  cpu          : usr=0.09%, sys=0.29%, ctx=34043, majf=0, minf=0
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,16384,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
  WRITE: bw=269MiB/s (282MB/s), 269MiB/s-269MiB/s (282MB/s-282MB/s), io=8192MiB (8590MB), run=30437-30437msec

```

runc:
```bash
[root@localhost hff]# kubectl exec -it test-runc sh
/data # fio --filename=/test/randread -direct=1 -iodepth 1 -thread -rw=randread  -ioengine=libaio -bs=512k  -numjobs=8 --size
=1G -group_reporting -name=randread
randread: (g=0): rw=randread, bs=(R) 512KiB-512KiB, (W) 512KiB-512KiB, (T) 512KiB-512KiB, ioengine=libaio, iodepth=1
...
fio-3.13
Starting 8 threads
randread: Laying out IO file (1 file / 1024MiB)
Jobs: 8 (f=8): [r(8)][100.0%][r=424MiB/s][r=847 IOPS][eta 00m:00s]
randread: (groupid=0, jobs=8): err= 0: pid=24: Thu Apr 21 07:28:56 2022
  read: IOPS=845, BW=423MiB/s (443MB/s)(8192MiB/19379msec)
    slat (usec): min=17, max=251, avg=30.18, stdev=10.86
    clat (usec): min=1555, max=16290, avg=9426.12, stdev=290.55
     lat (usec): min=1722, max=16325, avg=9456.52, stdev=289.00
    clat percentiles (usec):
     |  1.00th=[ 9110],  5.00th=[ 9241], 10.00th=[ 9241], 20.00th=[ 9372],
     | 30.00th=[ 9372], 40.00th=[ 9372], 50.00th=[ 9372], 60.00th=[ 9372],
     | 70.00th=[ 9503], 80.00th=[ 9503], 90.00th=[ 9503], 95.00th=[ 9634],
     | 99.00th=[10290], 99.50th=[10683], 99.90th=[12518], 99.95th=[13698],
     | 99.99th=[14877]
   bw (  KiB/s): min=417792, max=441233, per=99.99%, avg=432832.58, stdev=514.74, samples=304
   iops        : min=  816, max=  861, avg=845.21, stdev= 1.01, samples=304
  lat (msec)   : 2=0.01%, 4=0.02%, 10=98.26%, 20=1.71%
  cpu          : usr=0.07%, sys=0.45%, ctx=16484, majf=13, minf=1043
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=16384,0,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
   READ: bw=423MiB/s (443MB/s), 423MiB/s-423MiB/s (443MB/s-443MB/s), io=8192MiB (8590MB), run=19379-19379msec

Disk stats (read/write):
    dm-0: ios=16623/69, merge=0/0, ticks=156148/838, in_queue=157018, util=99.54%, aggrios=16636/57, aggrmerge=0/12, aggrticks=156399/675, aggrin_queue=157066, aggrutil=99.40%
  sda: ios=16636/57, merge=0/12, ticks=156399/675, in_queue=157066, util=99.40%
/data # fio --filename=/test/read --direct=1 --iodepth 1 --thread --rw=read --ioengine=psync --bs=512k --size=1G --numjobs=8
--group_reporting --name=read
read: (g=0): rw=read, bs=(R) 512KiB-512KiB, (W) 512KiB-512KiB, (T) 512KiB-512KiB, ioengine=psync, iodepth=1
...
fio-3.13
Starting 8 threads
read: Laying out IO file (1 file / 1024MiB)
Jobs: 8 (f=8): [R(8)][100.0%][r=432MiB/s][r=863 IOPS][eta 00m:00s]
read: (groupid=0, jobs=8): err= 0: pid=42: Thu Apr 21 07:29:33 2022
  read: IOPS=858, BW=429MiB/s (450MB/s)(8192MiB/19094msec)
    clat (usec): min=2823, max=20339, avg=9317.96, stdev=270.25
     lat (usec): min=2823, max=20339, avg=9318.19, stdev=270.24
    clat percentiles (usec):
     |  1.00th=[ 8979],  5.00th=[ 8979], 10.00th=[ 9110], 20.00th=[ 9110],
     | 30.00th=[ 9241], 40.00th=[ 9241], 50.00th=[ 9372], 60.00th=[ 9372],
     | 70.00th=[ 9372], 80.00th=[ 9503], 90.00th=[ 9503], 95.00th=[ 9634],
     | 99.00th=[ 9896], 99.50th=[10290], 99.90th=[10945], 99.95th=[11338],
     | 99.99th=[20317]
   bw (  KiB/s): min=432825, max=444416, per=99.99%, avg=439299.74, stdev=515.89, samples=304
   iops        : min=  843, max=  868, avg=857.84, stdev= 1.02, samples=304
  lat (msec)   : 4=0.02%, 10=99.26%, 20=0.71%, 50=0.01%
  cpu          : usr=0.05%, sys=0.45%, ctx=16428, majf=0, minf=1037
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=16384,0,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
   READ: bw=429MiB/s (450MB/s), 429MiB/s-429MiB/s (450MB/s-450MB/s), io=8192MiB (8590MB), run=19094-19094msec

Disk stats (read/write):
    dm-0: ios=16513/5, merge=0/0, ticks=152882/52, in_queue=152966, util=99.51%, aggrios=16710/5, aggrmerge=0/0, aggrticks=154761/52, aggrin_queue=154802, aggrutil=99.43%
  sda: ios=16710/5, merge=0/0, ticks=154761/52, in_queue=154802, util=99.43%
/data #
/data #
/data # fio --filename=/test/randwrite --direct=1 --iodepth 1 --thread --rw=randwrite --ioengine=psync --bs=512k --size=1G --
numjobs=8 --group_reporting --name=randwrite
randwrite: (g=0): rw=randwrite, bs=(R) 512KiB-512KiB, (W) 512KiB-512KiB, (T) 512KiB-512KiB, ioengine=psync, iodepth=1
...
fio-3.13
Starting 8 threads
randwrite: Laying out IO file (1 file / 1024MiB)
Jobs: 8 (f=8): [w(8)][97.4%][w=484MiB/s][w=968 IOPS][eta 00m:01s]
randwrite: (groupid=0, jobs=8): err= 0: pid=60: Thu Apr 21 07:30:44 2022
  write: IOPS=427, BW=214MiB/s (224MB/s)(8192MiB/38316msec); 0 zone resets
    clat (usec): min=1122, max=487875, avg=18682.47, stdev=32710.71
     lat (usec): min=1132, max=487891, avg=18703.72, stdev=32710.73
    clat percentiles (msec):
     |  1.00th=[    8],  5.00th=[    8], 10.00th=[    8], 20.00th=[    8],
     | 30.00th=[    9], 40.00th=[    9], 50.00th=[    9], 60.00th=[    9],
     | 70.00th=[    9], 80.00th=[    9], 90.00th=[   41], 95.00th=[   85],
     | 99.00th=[  153], 99.50th=[  165], 99.90th=[  372], 99.95th=[  422],
     | 99.99th=[  489]
   bw (  KiB/s): min=15360, max=507904, per=98.93%, avg=216581.34, stdev=26362.63, samples=608
   iops        : min=   30, max=  992, avg=422.79, stdev=51.50, samples=608
  lat (msec)   : 2=0.01%, 4=0.03%, 10=81.82%, 20=2.13%, 50=6.54%
  lat (msec)   : 100=6.49%, 250=2.77%, 500=0.20%
  cpu          : usr=0.13%, sys=0.18%, ctx=16462, majf=0, minf=7
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,16384,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
  WRITE: bw=214MiB/s (224MB/s), 214MiB/s-214MiB/s (224MB/s-224MB/s), io=8192MiB (8590MB), run=38316-38316msec

Disk stats (read/write):
    dm-0: ios=126/16340, merge=0/0, ticks=9757/320554, in_queue=330358, util=99.80%, aggrios=126/16497, aggrmerge=0/1, aggrticks=9757/321798, aggrin_queue=331542, aggrutil=99.76%
  sda: ios=126/16497, merge=0/1, ticks=9757/321798, in_queue=331542, util=99.76%
/data # fio --filename=/test/write --direct=1 --iodepth 1 --thread --rw=write --ioengine=psync --bs=512k --size=1G --numjobs=
8 --group_reporting --name=write
write: (g=0): rw=write, bs=(R) 512KiB-512KiB, (W) 512KiB-512KiB, (T) 512KiB-512KiB, ioengine=psync, iodepth=1
...
fio-3.13
Starting 8 threads
write: Laying out IO file (1 file / 1024MiB)
Jobs: 8 (f=8): [W(8)][100.0%][w=342MiB/s][w=684 IOPS][eta 00m:00s]
write: (groupid=0, jobs=8): err= 0: pid=78: Thu Apr 21 07:31:43 2022
  write: IOPS=417, BW=209MiB/s (219MB/s)(8192MiB/39259msec); 0 zone resets
    clat (usec): min=1713, max=633093, avg=19144.06, stdev=35066.38
     lat (usec): min=1719, max=633110, avg=19165.90, stdev=35066.17
    clat percentiles (msec):
     |  1.00th=[    8],  5.00th=[    8], 10.00th=[    8], 20.00th=[    8],
     | 30.00th=[    8], 40.00th=[    9], 50.00th=[    9], 60.00th=[    9],
     | 70.00th=[    9], 80.00th=[    9], 90.00th=[   41], 95.00th=[   87],
     | 99.00th=[  155], 99.50th=[  171], 99.90th=[  418], 99.95th=[  527],
     | 99.99th=[  575]
   bw (  KiB/s): min= 8192, max=507904, per=99.70%, avg=213032.02, stdev=26253.16, samples=621
   iops        : min=   16, max=  992, avg=415.81, stdev=51.27, samples=621
  lat (msec)   : 2=0.01%, 4=0.01%, 10=81.70%, 20=1.95%, 50=6.77%
  lat (msec)   : 100=6.52%, 250=2.80%, 500=0.20%, 750=0.05%
  cpu          : usr=0.14%, sys=0.17%, ctx=16433, majf=0, minf=7
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,16384,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
  WRITE: bw=209MiB/s (219MB/s), 209MiB/s-209MiB/s (219MB/s-219MB/s), io=8192MiB (8590MB), run=39259-39259msec

Disk stats (read/write):
    dm-0: ios=105/16393, merge=0/0, ticks=8635/323302, in_queue=331990, util=99.80%, aggrios=105/16502, aggrmerge=0/3, aggrticks=8635/324215, aggrin_queue=332852, aggrutil=99.77%
  sda: ios=105/16502, merge=0/3, ticks=8635/324215, in_queue=332852, util=99.77%

```




### fio_jobfile.fio
参考用例：https://os.51cto.com/article/618526.html
原文： https://www.stackhpc.com/kata-io-1.html


**测试命令：**
>echo 3 > /proc/sys/vm/drop_caches
export FIO_RW=write/read/randwrite/randread
fio fio_jobfile.fio --fallocate=none --runtime=30 --directory=/test --output-format=json+ --blocksize=65536 --output=${FIO_RW}.json


**fio_jobfile.fio**
```bash
[global] 
; Parameters common to all test environments 
; Ensure that jobs run for a specified time limit, not I/O quantity 
time_based=1 
; To model application load at greater scale, each test client will maintain 
; a number of concurrent I/Os. 
ioengine=libaio 
iodepth=8 
; Note: these two settings are mutually exclusive 
; (and may not apply for Windows test clients) 
direct=1 
buffered=0 
; Set a number of workers on this client 
thread=0 
numjobs=4 
group_reporting=1 
; Each file for each job thread is this size 
filesize=10g 
size=10g 
filename_format=$jobname.$jobnum.$filenum
[fio-job] 
; FIO_RW is read, write, randread or randwrite 
rw=${FIO_RW} 
```

rqy-k8s-1: 

|        | write | read | randwrite | randread | 
| ------ | ------ | ------ | ------- | ------ |
| kata   | 780 |   |    |   |
| runc   | 448 | 6710 | 443 | 6839 |

kata
```


```

runc
```bash
/test # export FIO_RW=write
/test # fio fio_jobfile.fio --fallocate=none --runtime=30 --directory=/test --output-format=json+ --blocksize=65536 --output=
${FIO_RW}.json
/test # (f=4): [W(4)][100.0%][w=28.0MiB/s][w=448 IOPS][eta 00m:00s]
/test # export FIO_RW=read
/test # fio fio_jobfile.fio --fallocate=none --runtime=30 --directory=/test --output-format=json+ --blocksize=65536 --output=
${FIO_RW}.json
/test # (f=3): [f(1),R(3)][100.0%][r=419MiB/s][r=6710 IOPS][eta 00m:00s]
/test # export FIO_RW=randwrite
/test # fio fio_jobfile.fio --fallocate=none --runtime=30 --directory=/test --output-format=json+ --blocksize=65536 --output=
${FIO_RW}.json
/test # (f=4): [w(4)][100.0%][w=27.7MiB/s][w=443 IOPS][eta 00m:00s]
/test # export FIO_RW=randread
/test # fio fio_jobfile.fio --fallocate=none --runtime=30 --directory=/test --output-format=json+ --blocksize=65536 --output=
${FIO_RW}.json
/test # (f=4): [r(4)][100.0%][r=427MiB/s][r=6839 IOPS][eta 00m:00s]
```


## sysbench
```
# 生成测试文件；
sysbench --test=fileio --file-num=10 --file-total-size=5G prepare 表示生成10个5G的文件

# 运行测试： 
sysbench --test=fileio --file-total-size=5G --file-test-mode=rndrw --max-requests=5000 --num-threads=16 --file-num=10 --file-extra-flags=direct --file-fsync-freq=0 --file-block-size=16384 run

sysbench --test=fileio cleanup


seqwr：顺序写
seqrewr：顺序重写
seqrd：顺序读
rndrd：随机读取
rndwr：随机写入
rndrw：混合随机读/写
```



10.208.11.110:

|      | rndrw：混合随机读/写                                         |
| ---- | ------------------------------------------------------------ |
| runc | 5368709120 bytes written in 44.82 seconds (114.23 MB/sec).<br>Total transferred 78.125Mb  (276.22Mb/sec)<br>total time:                          0.2828s |
| kata | 5368709120 bytes written in 106.52 seconds (48.07 MB/sec).<br>Total transferred 78.125Mb  (29.944Mb/sec)<br/>total time:                          2.6090s |
| 主机 | 5368709120 bytes written in 40.07 seconds (127.77 MiB/sec)<br>total time:                          0.2545s<br/>read, MiB/s:                  182.88<br/>written, MiB/s:               122.64 |


