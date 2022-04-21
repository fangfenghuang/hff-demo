# 存储空间资源限制ephemeral-storage

在每个Kubernetes的节点上，kubelet的根目录(默认是/var/lib/kubelet)和日志目录(/var/log)保存在节点的主分区上，这个分区同时也会被Pod的EmptyDir类型的volume、容器日志、镜像的层、容器的可写层所占用。

```yaml
resources:
  requests:
    cpu: 1
    memory: 2048Mi
    ephemeral-storage: 2Gi
  limits:
    cpu: 2
    memory: 2048Mi
    ephemeral-storage: 5Gi
```

## 使用限制：
- 只要Root Dir和kubelet --root-dir在一个分区，就能起作用
- 如果运行时指定了别的独立的分区，比如修改了docker的镜像层和容器可写层的存储位置(默认是/var/lib/docker)所在的分区，将不再将其计入ephemeral-storage的消耗。


## 效果
在容器写入超过存储限制时就会被驱逐掉
teststorage   0/1       Evicted   0         1m