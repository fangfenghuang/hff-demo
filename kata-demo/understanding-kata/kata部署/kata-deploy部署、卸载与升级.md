 [TOC]

# 前提

- 物理机已开启硬件虚拟化

  [[开启硬件虚拟化]]

- k8s容器运行时使用containerd（推荐）

  containerd配置参考《config.toml》、《0-containerd.conf》

# 部署
- 生成configuration.toml分发到各个kata节点（所有节点）
- 创建kata资源：kata-rbac.yaml、kata-deploy.yaml，创建runtimeClass： runtimeClass.yaml
```bash
$ kubectl apply -f kata-rbac.yaml
$ kubectl apply -f kata-deploy.yaml

$ kubectl -n kube-system wait --timeout=10m --for=condition=Ready -l name=kata-deploy pod

$ kubectl apply -f kata-runtimeClasses.yaml
```

- 各个kata节点（所有节点）创建软链接：
```bash
ln -s /opt/kata/bin/kata-runtime /usr/bin/kata-runtime
ln -s /opt/kata/bin/kata-monitor /usr/bin/kata-monitor
```

注：
- 可以增加节点标签设置kata-deploy daemonset节点亲和性设置kata节点
- 需要注意的是kata-deploy重启可能会导致默认的configuration.toml文件恢复默认配置，因此使用的是优先级更高的/etc/kata-containers/configuration.toml
- 使用ctr run创建的容器默认使用的是-qemu配置（/opt/kata/share/defaults/kata-containers/configuration-qemu.toml），如果需要使用ctr run测试，请同步配置到-qemu配置，或重定向shim链接到新的文件下


## 检查：
```bash
[root@rqy-k8s-1 ~]# kata-runtime check
  WARN[0000] Not running network checks as super user      arch=amd64 name=kata-runtime pid=48176 source=runtime
  System is capable of running Kata Containers
  System can currently create Kata Containers
```

# 卸载
- 删除所有kata容器
- 删除kata节点软连接
- 删除kata资源：kata-rbac.yaml、kata-deploy.yaml、runtimeClass.yaml
- 删除kata节点configuration.toml文件

```bash
$ kubectl delete -f kata-deploy.yaml

$ kubectl -n kube-system wait --timeout=10m --for=delete -l name=kata-deploy pod

$ kubectl apply -f kata-cleanup.yaml
# The cleanup daemon-set will run a single time, cleaning up the node-label, which makes it difficult to check in an automated fashion. This process should take, at most, 5 minutes.
# kubectl get pod -n kube-system | grep kubelet-kata-cleanup

$ kubectl delete -f kata-cleanup.yaml
$ kubectl delete -f kata-rbac.yaml
$ kubectl delete -f kata-runtimeClasses.yaml

```

# 升级(TODO)
