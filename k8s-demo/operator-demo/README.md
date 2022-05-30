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
```
$ make install

$ kubectl get crd|grep Hffdemo
```
## 修改controller.go，加入自己的业务逻辑。

## 启动controller，执行make run
```bash
$ make run

```

## 创建自定义资源对象
```bash
$ kubectl apply -f config/samples/
$ kubectl get  PVCUpload

```


## 将controller制作成docker镜像
```bash
cd $GOPATH/src/helloworld
make docker-build docker-push IMG=bolingcavalry/guestbook:002
```

## 部署controller
```bash
cd $GOPATH/src/helloworld
make deploy IMG=bolingcavalry/guestbook:002
```

## 卸载和清理
```bash
cd $GOPATH/src/helloworld
make uninstall
```


# 参考资料
https://operatorhub.io/


# 附件
## apiVersion 命名规则
group/version
- group 省略为核心组
- version 版本介绍
