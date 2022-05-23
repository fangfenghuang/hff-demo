
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






[root@telecom-k8s-phy02 hff]# kubectl get svc
NAME                   TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)
test-kata-httpd        NodePort    10.196.97.40     <none>        80:28394/TCP
test-runc-httpd        NodePort    10.196.26.12     <none>        80:21020/TCP 


```bash
./wrk -t64 -c20000 -d3m http://10.96.0.2:40080/ >> wrk-t64-c20000.log
sleep 3m
./wrk -t64 -c20000 -d3m http://10.96.0.2:21020/ >> wrk-t64-c20000.log
sleep 3m
./wrk -t64 -c20000 -d3m http://10.96.0.2:28394/ >> wrk-t64-c20000.log




```



sed -i -e 's/^sandbox_cgroup_only.*$/sandbox_cgroup_only=false/g' /etc/kata-containers/configuration.toml

kubectl get pod |grep test-kata-httpd| awk '{print $1}' | xargs -I {} kubectl delete pod {}

sleep 3m

./wrk -t64 -c20000 -d3m http://10.96.0.2:28394/

