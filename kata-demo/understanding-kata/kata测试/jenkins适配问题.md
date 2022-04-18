 

# jenkins升级

# docker重启导致docker.sock不可用问题

**通过挂载docker.dock在容器内使用docker:**

切换成containerd后，重启docker后容器不感知，但是里面挂载的docker.sock不能用了，这样就需要docker如果被重启，依赖docker.sock的容器也需要手动重启；

# kata容器不支持挂载主机路径的docker.sock

```bash
10:54:05  Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?

bash-4.4# docker images

Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
```

kataShared挂载类型由 [ virtio-fs][virtio-fs] 提供支持，虽然virtio-fs通常是一个很好的选择，但在 DinD 工作负载的情况下virtio-fs会导致问题-[it cannot be used as a "upper layer" of overlayfs without a custom patch](http://lists.katacontainers.io/pipermail/kata-dev/2020-January/001216.html)。
挂载docker时，docker无法启动，或者使用vfs的存储格式
vfs在 Kata Containers 中运行基于 DinD 的工作负载时必须采取特殊措施。


# 镜像管理命令问题
虽然流水线已抛弃了docker，但是backend提供了镜像上传下载的功能，主要依赖的命令如下：
>dockerLoadStr := "docker load -i /home/upload/" + header.Filename
dockerTagStr := "docker tag "
goPushNewImageCommandStr := "docker push "
 dockerPullStr := "docker pull "
 dockerSaveStr := "cd /home/download && docker save -o /home/download/"


## 解决方案
1. DinD
```yaml
---

apiVersion: apps/v1

kind: Deployment

metadata:

 name: dind-test

spec:

 replicas: 1

 selector:

 matchLabels:

 app: dind-test

 template:

 metadata:

 labels:

 app: dind-test

 spec:

 runtimeClassName: kata-containers

 containers:

 - name: docker

 image: docker:20.10-dind

 command: ["sh", "-c"]

 args:

 - if [[ $(df -PT /var/lib/docker | awk 'NR==2 {print $2}') == virtiofs ]]; then

 apk add e2fsprogs &&

 truncate -s 20G /tmp/disk.img &&

 mkfs.ext4 /tmp/disk.img &&

 mount /tmp/disk.img /var/lib/docker; fi &&

 dockerd-entrypoint.sh;

 securityContext:

 privileged: true

---

apiVersion: apps/v1

kind: Deployment

metadata:

 name: dind-test

spec:

 replicas: 1

 selector:

 matchLabels:

 app: dind-test

 template:

 metadata:

 labels:

 app: dind-test

 spec:

 runtimeClassName: kata-containers

 containers:

 - name: dind

 image: 'docker:stable-dind'

 command:

 - dockerd

 - --host=unix:///var/run/docker.sock

 - --host=tcp://0.0.0.0:8000

 securityContext:

 privileged: true

 volumeMounts:

 - mountPath: /var/run

 name: cache-dir

 - name: clean-ci

 image: 'docker:stable'

 volumeMounts:

 - mountPath: /var/run

 name: cache-dir

 volumes:

 - name: cache-dir

 emptyDir: {}
```

2. 命令替换 在容器内挂载主机路径后可以使用crictl和ctr（建议）

docker load: ctr -n k8s.io image import app.tar
docker tag:  ctr -n k8s.io image tag imgname newname
docker push: ctr -n k8s.io image push imgname 或 crictl push
docker pull : ctr -n k8s.io image pull imgname 或 crictl pull
docker save: ctr -n k8s.io image export app.tar imgname


```yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: test-crictl
spec:
  selector:
    matchLabels:
      app: test-crictl
  template:
    metadata:
      labels:
        app: test-crictl
    spec:
      volumes:
        - name: a
          hostPath:
            path: /usr/bin/crictl
            type: File
        - name: b
          hostPath:
            path: /usr/bin/ctr
            type: File
        - name: c
          hostPath:
            path: /run/containerd/containerd.sock
            type: File
      containers:
        - name: test
          image: nginx
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - name: a
              mountPath: /usr/bin/crictl
            - name: b
              mountPath: /usr/bin/ctr
            - name: c
              mountPath: /run/containerd/containerd.sock
```
注意：需要显示配置默认的endpoints:
```bash
crictl config runtime-endpoint unix:///run/containerd/containerd.sock
crictl config image-endpoint unix:///run/containerd/containerd.sock
```

建议：
    镜像上传下载完成后清楚本地缓存，否则可能会占用较多临时存储空间


