# webhook-demo
pod annotation:
- mutate处理：
annotations:
    kubernetes.io/ingress-bandwidth: 1M
    kubernetes.io/egress-bandwidth: 1M

    如果annotation已存在，则不修改

- validate处理：
验证pod是否有annotation

## build
sh build.sh v0.0.0


## deploy
```bash
cd deploy
# 生成证书：如不修改yaml参数，
sh webhook-create-signed-cert.sh
creating certs in tmpdir /tmp/tmp.Xez03UX35o
Generating RSA private key, 2048 bit long modulus
......................................................................+++
.+++
e is 65537 (0x10001)
certificatesigningrequest.certificates.k8s.io/admission-webhook-demo-svc.default created
NAME                                 AGE   REQUESTOR          CONDITION
admission-webhook-demo-svc.default   0s    kubernetes-admin   Pending
certificatesigningrequest.certificates.k8s.io/admission-webhook-demo-svc.default approved
secret/admission-webhook-demo-certs created
```


# 部署
```
kubectl apply -f deployment.yaml
```

# 填充占位符，生成webhook
```bash
chmod +x webhook-patch-ca-bundle.sh
cat webhook-temp.yaml | ./webhook-patch-ca-bundle.sh > webhook.yaml
kubectl apply -f webhook.yaml
```
# 验证
>kubectl apply -f test-deplyment.yaml 



## 网络限速性能测试

# 启动测试pod（安装iperf）
```bash
[root@tztest iperf]# kubectl get pod -n hffns -o wide
NAME                              READY   STATUS    RESTARTS   AGE     IP               NODE     NOMINATED NODE   READINESS GATES
netperf-server                    1/1     Running   0          43s     10.242.235.107   tztest   
```

# 在宿主机上：
```bash
[root@tztest iperf]# iperf3 -c 10.242.235.107
Connecting to host 10.242.235.107, port 5201
[  4] local 10.19.0.13 port 45376 connected to 10.242.235.107 port 5201
[ ID] Interval           Transfer     Bandwidth       Retr  Cwnd
[  4]   0.00-1.00   sec   750 MBytes  6.29 Gbits/sec    0    397 KBytes
[  4]   1.00-2.00   sec   914 MBytes  7.67 Gbits/sec    0    443 KBytes
[  4]   2.00-3.00   sec   917 MBytes  7.69 Gbits/sec    0    477 KBytes
[  4]   3.00-4.00   sec   935 MBytes  7.84 Gbits/sec    0    503 KBytes
[  4]   4.00-5.00   sec   880 MBytes  7.38 Gbits/sec    0    529 KBytes
[  4]   5.00-6.00   sec   903 MBytes  7.57 Gbits/sec    0    556 KBytes
[  4]   6.00-7.00   sec   936 MBytes  7.86 Gbits/sec    0    581 KBytes
[  4]   7.00-8.00   sec   899 MBytes  7.54 Gbits/sec    0    609 KBytes
[  4]   8.00-9.00   sec   856 MBytes  7.18 Gbits/sec    0    634 KBytes
^C[  4]   9.00-9.34   sec   333 MBytes  8.36 Gbits/sec    0    644 KBytes
[ ID] Interval           Transfer     Bandwidth       Retr
[  4]   0.00-9.34   sec  8.13 GBytes  7.48 Gbits/sec    0             sender
[  4]   0.00-9.34   sec  0.00 Bytes  0.00 bits/sec                  receiver
iperf3: interrupt - the client has terminated

# tc qdisk show
[root@tztest iperf]# tc qdisc show
qdisc htb 20: dev cali90ab6b94883 root refcnt 2 r2q 1 default 0 direct_packets_stat 22237
qdisc ingress ffff: dev cali90ab6b94883 parent ffff:fff1 ----------------
qdisc tbf 1: dev 31b3 root refcnt 2 rate 1000Kbit burst 27917286b lat 1924.2s
```
