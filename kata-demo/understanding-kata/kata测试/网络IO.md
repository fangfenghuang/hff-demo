 [TOC]

# 注意
qperf测试service有问题：
1. 对于runc容器，需要修改qperf服务监听端口，否则跨节点无法测试
2. 对于kata容器，本节点和跨节点都无法测试


# 测试工具与指标
测试指标：用qperf测试带宽和延迟
- 循环测试1bytes-64KiB的带宽和延迟
```bash
qperf <ip> -oo msg_size:1:256K:*64 -vu tcp_bw tcp_lat
```
# 测试方法及测试步骤：

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

测试pod:
|               |msg_size   |宿主机服务端|runc容器服务端|kata容器服务端|
|---------------|-----------|-----------|-------------|-------------|
| 从主机A到服务  |1bytes     |tcp_lat:13 us<br>tcp_bw:1.13 MB/sec |tcp_lat:12.5 us<br>tcp_bw:1.14 MB/sec | tcp_lat:20.2 us<br>tcp_bw:1.05 MB/sec
|               |64bytes    |tcp_lat:12.9 us<br>tcp_bw:77 MB/sec |tcp_lat:12.6 us<br>tcp_bw:42.6 MB/sec | tcp_lat:19.3 us<br>tcp_bw:45.6 MB/sec
|               |4bytes     |tcp_lat:14.1 us<br>tcp_bw:1.71 GB/sec |tcp_lat:20.2 us<br>tcp_bw:335 MB/sec | tcp_lat:25.7 us<br>tcp_bw:355 MB/sec
|               |256bytes   |tcp_lat:68.8 us<br>tcp_bw:3.48 GB/sec |tcp_lat:92.1 us<br>tcp_bw:3.39 GB/sec | tcp_lat:167 us<br>tcp_bw:3.13 GB/sec
|跨节点主机到服务|1bytes     |tcp_lat:15.9 us<br>tcp_bw:1.19 MB/sec|tcp_lat:19.2 us<br>tcp_bw:1.17 MB/sec   | tcp_lat:28.9 us<br>tcp_bw:1.08 MB/sec
|               |64bytes    |tcp_lat:16.2 us<br>tcp_bw:54.2 MB/sec|tcp_lat:18.8 us<br>tcp_bw:69 MB/sec    | tcp_lat:29.7 us<br>tcp_bw:55 MB/sec
|               |4bytes     |tcp_lat:28.8 us<br>tcp_bw:1.13 GB/sec|tcp_lat:34.1 us<br>tcp_bw:1.16 GB/sec  | tcp_lat:45.7 us<br>tcp_bw:1.13 GB/sec
|               |256bytes   |tcp_lat:315 us<br>tcp_bw:1.15 GB/sec|tcp_lat:397 us<br>tcp_bw:1.17 GB/sec  | tcp_lat:338 us<br>tcp_bw:1.16 GB/sec


 

# 附件

```bash
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: qperf-server-kata
spec:
  replicas: 1
  selector:
    matchLabels:
      app: qperf-server-kata
  template:
    metadata:
      labels:
        app: qperf-server-kata
    spec:
      runtimeClassName: kata
      containers:
      - image: "arjanschaaf/centos-qperf"
        args: ["-lp", "4000"]
        imagePullPolicy: "IfNotPresent"
        name: "qperf-server-kata"
        ports:
        - containerPort: 4000
          name: "p1udp"
          protocol: UDP
        - containerPort: 4000
          name: "p1tcp"
          protocol: TCP
        - containerPort: 4001
          name: "p2tcp"
          protocol: TCP
        - containerPort: 4001
          name: "p2udp"
          protocol: UDP
      restartPolicy: Always

---
apiVersion: v1
kind: Service
metadata:
  name: "qperf-server-kata"
spec:
  selector:
    k8s-app: "qperf-server-kata"
  ports:
  - name: "p1udp"
    port: 4000
    targetPort: 4000
    protocol: UDP
  - name: "p1tcp"
    port: 4000
    targetPort: 4000
    protocol: TCP
  - name: "p2udp"
    port: 4001
    targetPort: 4001
    protocol: UDP
  - name: "p2tcp"
    port: 4001
    targetPort: 4001
    protocol: TCP
```

```bash
[root@telecom-k8s-phy01 kbuser]# kubectl get node -o wide
NAME                STATUS   ROLES    AGE     VERSION                               INTERNAL-IP   EXTERNAL-IP   OS-IMAGE                      KERNEL-VERSION                CONTAINER-RUNTIME
telecom-k8s-phy01   Ready    master   24d     v1.17.2                               10.96.0.1     <none>        CentOS Linux 7 (Core)         3.10.0-1160.59.1.el7.x86_64   containerd://1.4.6
telecom-k8s-phy02   Ready    master   24d     v1.17.2                               10.96.0.2     <none>        CentOS Linux 7 (Core)         3.10.0-1160.59.1.el7.x86_64   containerd://1.4.6

[root@telecom-k8s-phy01 kbuser]# kubectl get pod -o wide
NAME                                 READY   STATUS    RESTARTS   AGE     IP               NODE                NOMINATED NODE   READINESS GATES
qperf-server-kata-5d9bffcf97-gw6xb   1/1     Running   0          2d23h   10.196.192.79    telecom-k8s-phy01   <none>           <none>
qperf-server-z5rwk                   1/1     Running   0          3d15h   10.196.192.101   telecom-k8s-phy01   <none>           <none>
```

## 宿主机
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
## runc
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

## kata
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