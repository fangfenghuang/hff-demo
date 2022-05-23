# kubefileborwser

https://github.com/xmapst/kubefilebrowser

kubernetes container file browser
通过kube-apiserver(kubectl exec)访问容器，进行文件传输及执行文件查询创建等操作；
由于pod镜像可能缺少需要的文件操作命令，因此在文件浏览器功能接口里需要先上传一个封装的二进制工具包（kftools）到目标容器，通过kftools执行ls、cat、touch、zip等命令。
对于pvc上传需要考虑未挂载pod的情况，后端会创建一个upload pod用于执行文件上传，如果pvc已经挂载pod，则需要将新建的upload pod绑定到pvc挂载pod的某个节点。

https://www.yfdou.com/archives/kuberneteszhi-kubectlexeczhi-ling-gong-zuo-yuan-li-shi-xian-copyhe-webshellyi-ji-filebrowser.html#Exec-%E6%A0%B8%E5%BF%83%E6%BA%90%E7%A0%81

## 启动可选环境变量


+ [golang 1.16 gin static embed](https://mojotv.cn/golang/golang-html5-websocket-remote-desktop)
+ [vue](https://cli.vuejs.org/config/)
+ [kubectl copy & shell 原理讲解](https://www.yfdou.com/archives/kuberneteszhi-kubectlexeczhi-ling-gong-zuo-yuan-li-shi-xian-copyhe-webshellyi-ji-filebrowser.html)


# 本地调试
## 后端
```bash
[root@tztest ~]# ls ~/.kube/config
[root@tztest server]# cd kubefilebrowser/cmd/server
[root@tztest server]# go run main.go
2022/04/22 10:51:55 INFO main.go:69 kubernetes file browser is running ...
2022/04/22 10:51:55 INFO main.go:62 listen address [0.0.0.0:9999]
```

## 前端
```bash
# 安装npm yarn 
[root@tztest webui]# cd /opt;wget https://nodejs.org/dist/v14.15.4/node-v14.15.4-linux-x64.tar.xz
[root@tztest webui]# tar -xvf node-v14.15.4-linux-x64.tar.xz
[root@tztest webui]# ln -s /opt/node-v14.15.4-linux-x64/bin/npm /usr/bin/
[root@tztest webui]# ln -s /opt/node-v14.15.4-linux-x64/bin/node /usr/bin/
[root@tztest webui]# ln -s /opt/node-v14.15.4-linux-x64/bin/npx /usr/bin/
[root@tztest webui]# npm config set registry https://registry.npm.taobao.org
[root@tztest webui]# npm config get registry
[root@tztest webui]# npm install --global yarn
[root@tztest webui]# ln -s /opt/node-v14.15.4-linux-x64/bin/yarn /usr/bin/
[root@tztest webui]# cd webui/
[root@tztest webui]# yarn install
[root@tztest webui]# yarn run serve
[root@tztest webui]# yarn run build
```


# 编译
```bash
docker build --network host --rm --build-arg APP_ROOT=/go/src/kubefilebrowser -t kubefilebrowser:latest -f Dockerfile .
```


# 部署
## RUN_MODE=debug: 使用KUBECONFIG
```yaml
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: kubefilebrowser
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubefilebrowser
  template:
    metadata:
      labels:
        app: kubefilebrowser
    spec:
      volumes:
        - name: k8s-config
          hostPath:
            path: /root/.kube
            type: ''
      containers:
        - name: kubefilebrowser
          image: kubefilebrowser:latest
          volumeMounts:
            - name: k8s-config
              mountPath: /root/.kube
          imagePullPolicy: IfNotPresent
```


## RUN_MODE=release: 使用serviceaccount
```yaml
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kubefilebrowser-clusterrole
rules:
  - apiGroups:
      - ""
    resources:
      - nodes
      - pods
      - pods/exec
      - persistentvolumeclaims
      - namespaces
    verbs:
      - "*"
  - apiGroups:
      - "apps"
    resources:
      - deployments
    verbs:
      - "*"

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubefilebrowser-sa
  namespace: default

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kubefilebrowser-clusterrolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kubefilebrowser-clusterrole
subjects:
- kind: ServiceAccount
  name: kubefilebrowser-sa
  namespace: default


---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: kubefilebrowser
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubefilebrowser
  template:
    metadata:
      labels:
        app: kubefilebrowser
    spec:
      serviceAccount: kubefilebrowser-sa
      containers:
        - name: kubefilebrowser
          image: kubefilebrowser:latest
          env:
            - name: RUN_MODE
              value: release
          imagePullPolicy: IfNotPresent

```
