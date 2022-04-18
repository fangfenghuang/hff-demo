https://github.com/kubevirt/kubevirt/releases/
## 版本：
v0.45.0

## 镜像：
image: quay.io/kubevirt/virt-operator:v0.45.0


## 环境检查：
[root@tztest ~]# virt-host-validate qemu

## 修改kube-apiserver配置
--requestheader-client-ca-file=/etc/kubernetes/ssl/ca.pem \


## 部署：
```
# Pick an upstream version of KubeVirt to install 
$ export RELEASE=v0.35.0 
# Deploy the KubeVirt operator 
$ kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/${RELEASE}/kubevirt-operator.yaml 
# Create the KubeVirt CR (instance deployment request) which triggers the actual installation
$ kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/${RELEASE}/kubevirt-cr.yaml 
# wait until all KubeVirt components are up 
$ kubectl -n kubevirt wait kv kubevirt --for condition=Available
```


## 拷贝virtctl:
```
cp virtctl-v0.30.0-linux-x86_64 virtctl
chmod +x virtctl
mv virtctl /usr/local/bin/virtctl
```
