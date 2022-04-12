 [TOC]

# 前提

- 物理机已开启硬件虚拟化

  [[开启硬件虚拟化]]

- k8s容器运行时使用containerd（推荐）

  containerd配置参考《config.toml》、《0-containerd.conf》

# ansible部署

- 生成configuration.toml分发到各个kata节点（所有节点）

- 创建kata资源：kata-rbac.yaml、kata-deploy.yaml

- 创建runtimeClass： runtimeClass.yaml

- 各个kata节点（所有节点）创建软链接：

``
ln -s /opt/kata/bin/kata-runtime /usr/bin/kata-runtime
ln -s /opt/kata/bin/kata-monitor /usr/bin/kata-monitor
ln -s /opt/kata/bin/containerd-shim-kata-v2 /usr/local/bin/containerd-shim-kata-v2
``

注：可以增加节点标签设置kata-deploy daemonset节点亲和性设置kata节点

## 检查：
>[root@rqy-k8s-1 ~]# kata-runtime check
  WARN[0000] Not running network checks as super user      arch=amd64 name=kata-runtime pid=48176 source=runtime
  System is capable of running Kata Containers
  System can currently create Kata Containers


# 卸载
- 删除所有kata容器

- 删除kata节点软连接

- 删除kata资源：kata-rbac.yaml、kata-deploy.yaml、runtimeClass.yaml

- 删除kata节点configuration.toml文件

# 升级(TODO)
