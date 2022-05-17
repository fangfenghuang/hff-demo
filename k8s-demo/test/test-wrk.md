
# wrk2
https://github.com/giltene/wrk2
https://blog.csdn.net/ccccsy99/article/details/105958366?spm=1001.2101.3001.6650.2&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7Edefault-2-105958366-blog-88062573.pc_relevant_baidufeatures_v5&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7Edefault-2-105958366-blog-88062573.pc_relevant_baidufeatures_v5&utm_relevant_index=5



```bash
yum -y install gcc automake autoconf libtool make
yum install gcc gcc-c++
yum install openssl-devel
yum install libssl-dev

ln -s /export/ccc/wrk2-master/wrk /usr/local/bin
```




```bash
---
apiVersion: v1
kind: Pod
metadata:
  name: test-runc-httpd
  labels:
    app: test-runc-httpd
spec:
  nodeName: telecom-k8s-phy02
  containers:
  - name: httpd-runc
    image: httpd
    imagePullPolicy: IfNotPresent
    volumeMounts:
    - mountPath: /test
      name: test-volume
    resources:
      limits:
        memory: "2Gi"
        cpu: "1"
  volumes:
  - name: test-volume
    hostPath:
      path: /hff/test/runc
      type: Directory
---
apiVersion: v1
kind: Pod
metadata:
  name: test-kata-httpd
  labels:
    app: test-kata-httpd
spec:
  nodeName: telecom-k8s-phy02
  runtimeClassName: kata
  containers:
  - name: httpd-kata
    image: httpd
    imagePullPolicy: IfNotPresent
    volumeMounts:
    - mountPath: /test
      name: test-volume
    resources:
      limits:
        memory: "2Gi"
        cpu: "1"
  volumes:
  - name: test-volume
    hostPath:
      path: /hff/test/kata
      type: Directory
```


```bash
[root@telecom-k8s-phy01 hff]# kubectl get pod -o wide
NAME                                 READY   STATUS    RESTARTS   AGE     IP               NODE                NOMINATED NODE   READINESS GATES
qperf-server-kata-5d9bffcf97-gw6xb   1/1     Running   0          4d6h    10.196.192.79    telecom-k8s-phy01   <none>           <none>
qperf-server-z5rwk                   1/1     Running   0          4d22h   10.196.192.101   telecom-k8s-phy01   <none>           <none>
test-kata                            1/1     Running   0          4d4h    10.196.192.195   telecom-k8s-phy01   <none>           <none>
test-kata-httpd                      1/1     Running   0          11s     10.196.142.171   telecom-k8s-phy02   <none>           <none>
test-runc                            1/1     Running   0          4d4h    10.196.192.144   telecom-k8s-phy01   <none>           <none>
test-runc-httpd                      1/1     Running   0          17s     10.196.142.131   telecom-k8s-phy02   <none>           <none>
[root@telecom-k8s-phy01 hff]# kubectl get node
NAME                STATUS   ROLES    AGE     VERSION
telecom-k8s-phy01   Ready    master   26d     v1.17.2
telecom-k8s-phy02   Ready    master   26d     v1.17.2

```