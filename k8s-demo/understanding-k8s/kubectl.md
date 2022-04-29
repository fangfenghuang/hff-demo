# kubectl cp
从Pod容器中copy文件至本地
kubectl cp <file-spec-src> <file-spec-dest> [options]
kubectl cp -c cloud-loan-gate uat/cloud-786d84c554-p7jz7:app/logs/app/cloud.log cloud.log
-c :指定容器（因pod中有多个container，默认从第一个，有可能报错找不到文件和目录）
源目录参数时，：后只能跟文件名，不能是以“/”开头的路径（eg:app/logs/cloud.log）
目标参数时，必须为文件不能是一个目录(eg:cloud.log)

--no-preserve=false

## !!!Important Note!!!
要求容器镜像中存在“tar”二进制文件。 如果“tar”不存在，“kubectl cp”将失败。  


## 对于高级用例，如符号链接、通配符扩展或文件模式保存，考虑使用'kubectl exec'  
使用 tar -cf - 将具有文件夹结构的数据转换成数据流，再通过 linux 管道接收这个数据流；通过 tar -xf - 将数据流转换成 linux 文件系统。

# Copy /tmp/foo local file to /tmp/bar in a remote pod in namespace <some-namespace>
tar cf - /tmp/foo | kubectl exec -i -n <some-namespace> <some-pod> -- tar xf - -C /tmp/ba

# Copy /tmp/foo from a remote pod to /tmp/bar locally
kubectl exec -n <some-namespace> <some-pod> -- tar cf - /tmp/foo | tar xf - -C /tmp/ba

```

```
# Copy /tmp/foo_dir local directory to /tmp/bar_dir in a remote pod in the default namespace
kubectl cp /tmp/foo_dir <some-pod>:/tmp/bar_di

# Copy /tmp/foo local file to /tmp/bar in a remote pod in a specific container
kubectl cp /tmp/foo <some-pod>:/tmp/bar -c <specific-container

# Copy /tmp/foo local file to /tmp/bar in a remote pod in namespace <some-namespace>
kubectl cp /tmp/foo <some-namespace>/<some-pod>:/tmp/ba

# Copy /tmp/foo from a remote pod to /tmp/bar locally
kubectl cp <some-namespace>/<some-pod>:/tmp/foo /tmp/bar

```
# kubectl exec



## 参考：
https://www.yfdou.com/archives/kuberneteszhi-kubectlexeczhi-ling-gong-zuo-yuan-li-shi-xian-copyhe-webshellyi-ji-filebrowser.html
