operator 是一种 kubernetes 的扩展形式，利用自定义资源对象（Custom Resource）来管理应用和组件，允许用户以 Kubernetes 的声明式 API 风格来管理应用及服务。operator 定义了一组在 Kubernetes 集群中打包和部署复杂业务应用的方法，operator主要是为解决特定应用或服务关于如何运行、部署及出现问题时如何处理提供的一种特定的自定义方式。


# operator工具
- operator SDK —— operator framework
- KUDO (Kubernetes 通用声明式 Operator)
- kubebuilder，kubernetes SIG 在维护的一个项目
- Metacontroller，可与 Webhook 结合使用，以实现自己的功能。


## operator-sdk or Kubebuilder 
https://github.com/kubernetes-sigs/kubebuilder/blob/master/designs/integrating-kubebuilder-and-osdk.md


# Kubebuilder
## 安装
```bash
$ os=$(go env GOOS)
$ arch=$(go env GOARCH)

# 下载 kubebuilder 并解压到 tmp 目录中
$ curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
$ chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
# /usr/local/bin/加入PATH
$ kubebuilder version

```

## 新建项目
```bash
# 创建一个目录，然后在里面运行 kubebuilder init 命令
# 初始化项目
$ mkdir kubebuilder-demo

$ kubebuilder init --domain fangfenghuang.io --owner fangfenghuang --repo github.com/fangfenghuang/kubebuilder-demo

```


## 新建一个 API
创建一个新的 API（组/版本）为 hffapp/v1，并在上面创建新的 Kind(CRD) Hffdemo.
```bash
$ kubebuilder create api --group hffapp --version v1 --kind HffDemo
```

## make install
```bash

$ go run main.go  # 先编译看看

$ make install # 部署CRD
GOBIN=/home/kbuser/hff/kubebuilder-demo/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.8.0
/home/kbuser/hff/kubebuilder-demo/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/home/kbuser/hff/kubebuilder-demo/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/hffdemoes.hffapp.fangfenghuang.io created


$ kubectl get crd|grep hff
hffdemoes.hffapp.fangfenghuang.io                    2022-06-06T07:18:53Z

```
## 修改controller.go，加入自己的业务逻辑。

## 启动controller，执行make run
```bash
$ make run  # 本地运行controller  也可以将其部署到 Kubernetes 上运行

```

## 创建自定义资源对象
```bash
$ kubectl apply -f config/samples/
$ kubectl get hffdemoes.hffapp.fangfenghuang.io
NAME             AGE
hffdemo-sample   11s

$ kubectl get hffdemoes.hffapp.fangfenghuang.io hffdemo-sample -o yaml

```


## 将controller制作成docker镜像


## 部署controller


## 卸载和清理



# 参考资料
https://operatorhub.io/


# 附件
## apiVersion 命名规则
group/version
- group 省略为核心组
- version 版本介绍


# 错误问题
```bash
$ make install # 部署CRD
/home/kbuser/hff/kubebuilder-demo/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash -s -- 3.8.7 /home/kbuser/hff/kubebuilder-demo/bin
make: *** [/home/kbuser/hff/kubebuilder-demo/bin/kustomize] 错误 35
```
单独执行
curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash -s -- 3.8.7 /home/kbuser/hff/kubebuilder-demo/bin
然后继续执行make install



# operator 开发
## 填充CRD： types.go


## 生成webhook框架
operator中的webhook，其作用与上述过滤器类似，外部对CRD资源的变更，在Controller处理之前都会交给webhook提前处理

## 填充controller
- 增加权限注释
Reconcile方法前面有一些+kubebuilder:rbac前缀的注释，这些是用来确保controller在运行时有对应的资源操作权限

- 增加入队逻辑，关联资源对象的触发机制
- 填充业务逻辑：修改reconcile函数，处理工作队列，将逻辑结果反馈回status；需要注意的是Reconcile函数里return  reconcile.Result{}, err会重新入队


kubebuilder 2.0中，构建一个reconciler时，可以用Own,Watch方法来额外监听一些资源，但是For方法必须要有，如果没有For方法，编译出来的程序运行时会报错，类似于"kind.Type should not be empty"

## finalizers
Finalizers 是由字符串组成的列表，当 Finalizers 字段存在时，相关资源不允许被强制删除。存在 Finalizers 字段的的资源对象接收的第一个删除请求设置 metadata.deletionTimestamp 字段的值， 但不删除具体资源，在该字段设置后， finalizer 列表中的对象只能被删除，不能做其他操作。

当 metadata.deletionTimestamp 字段非空时，controller watch 对象并执行对应 finalizers 的动作，当所有动作执行完后，需要清空 finalizers ，之后 k8s 会删除真正想要删除的资源。



