 [TOC]

# 网络IO对比
```bash
netperf-server-default             1/1     Running   0          18d     10.192.181.19    rqy-k8s-1 
netperf-server-kata                1/1     Running   0          18d     10.192.173.212   rqy-k8s-3 
netperf-server-runc                1/1     Running   0          18d     10.192.173.213   rqy-k8s-3 
```


## iperf3

```yaml
apiVersion: v1
kind: Pod
metadata:
 name: netperf-server-kata
spec:
 runtimeClassName: kata-containers
 containers:
 - image: sirot/netperf-latest
 command: ["/bin/sh","-c","netserver -p 4444 -4; iperf3 -s -i 1;"]
 imagePullPolicy: IfNotPresent
 name: netperf
 ports:
 - name: netperf-port
 containerPort: 4444
 - name: iperf-port
 containerPort: 5210
 restartPolicy: Always
```

> iperf3 -c pod_ip

**10.91.0.3**

|      | sender         | receiver       |
| ---- | -------------- | -------------- |
| kata | 9.30 Gbits/sec | 9.30 Gbits/sec |
| runc | 9.31 Gbits/sec | 9.31 Gbits/sec |
| 主机 | 9.27 Gbits/sec | 9.27 Gbits/sec |



## ab压测

> ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/library/httpd:latest hfftest sh
>
> ab -kn 100000 -c 100 [http://10.241.102.146](http://pod_ip)/

**10.91.0.3**

|      | kata           |
| ---- | -------------- |
| kata | 9704.11 req/s  |
| runc | 13740.20 req/s |



## 网络限制kubernetes.io/ingress-bandwidth

 

```yaml
annotations:

  kubernetes.io/egress-bandwidth: 1M
  kubernetes.io/ingress-bandwidth: 1M
```
```bash
[ ID] Interval      Transfer   Bandwidth    Retr
[  4]  0.00-10.00  sec  23.7 MBytes  19.8 Mbits/sec   2       sender
[  4]  0.00-10.00  sec  21.8 MBytes  18.2 Mbits/sec          receiver
```
 