# 当前集群环境
containerd: 1.4.6
k8s: v1.17.2
kernel: 3.10.0-1160.59.1.el7.x86_64



# kata 2.4.0
## 兼容性
Kata Containers 2.4.0 与 contaienrd v1.5.2 兼容
Kata Containers 2.4.0 与 Kubernetes 1.23.1-00 兼容
Kata Containers 2.4.0 建议使用Linux kernel v5.15.26

## 默认Image操作系统
aarch64:
  name: "ubuntu"
  version: "latest"
ppc64le:
  name: "ubuntu"
  version: "latest"
s390x:
  name: "ubuntu"
  version: "latest"
x86_64:
  name: "clearlinux"
  version: "latest"

meta:
image-type: "clearlinux"

## 默认initrd操作系统
aarch64:
  name: "alpine"
  version: "3.15"
ppc64le:
  name: "ubuntu"
  version: "20.04"
s390x:
  name: "ubuntu"
  version: "20.04"
x86_64:
  name: "alpine"
  version: "3.15"


