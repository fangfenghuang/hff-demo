
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
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: test-runc-httpd
spec:
  selector:
    matchLabels:
      app: test-runc-httpd
  template:
    metadata:
      labels:
        app: test-runc-httpd
    spec:
      nodeName: telecom-k8s-phy02
      containers:
      - name: httpd-runc
        image: httpd
        imagePullPolicy: IfNotPresent
        resources:
          limits:
            memory: "256Gi"
            cpu: "64"
          requests:
            memory: "2Gi"
            cpu: "1"
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: test-kata-httpd
spec:
  selector:
    matchLabels:
      app: test-kata-httpd
  template:
    metadata:
      labels:
        app: test-kata-httpd
    spec:
      nodeName: telecom-k8s-phy02
      containers:
      - name: httpd-kata
        image: httpd
        imagePullPolicy: IfNotPresent
        resources:
          limits:
            memory: "256Gi"
            cpu: "64"
          requests:
            memory: "2Gi"
            cpu: "1"
```

