# kubernetes webshell实现
k8s.io/client-go/tools/remotecommand kubernetes client-go 提供的 remotecommand 包，提供了方法与集群中的容器建立长连接，并设置容器的 stdin，stdout 等。
remotecommand 包提供基于 SPDY 协议的 Executor interface，进行和 pod 终端的流的传输。初始化一个 Executor 很简单，只需要调用 remotecommand 的 NewSPDYExecutor 并传入对应参数。
Executor 的 Stream 方法，会建立一个流传输的连接，直到服务端和调用端一端关闭连接，才会停止传输。常用的做法是定义一个如下 PtyHandler 的 interface，然后使用你想用的客户端实现该 interface 对应的Read(p []byte) (int, error)和Write(p []byte) (int, error)方法即可，调用 Stream 方法时，只要将 StreamOptions 的 Stdin Stdout 都设置为 ptyHandler，Executor 就会通过你定义的 write 和 read 方法来传输数据。

```bash
// PtyHandler
type PtyHandler interface {
    io.Reader
    io.Writer
    remotecommand.TerminalSizeQueue
}

// NewSPDYExecutor
req := kubeClient.CoreV1().RESTClient().Post().
        Resource("pods").
        Name(podName).
        Namespace(namespace).
        SubResource("exec")
req.VersionedParams(&corev1.PodExecOptions{
    Container: containerName,
    Command:   cmd,
    Stdin:     true,
    Stdout:    true,
    Stderr:    true,
    TTY:       true,
}, scheme.ParameterCodec)
executor, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
if err != nil {
    log.Printf("NewSPDYExecutor err: %v", err)
    return err
}

// Stream
err = executor.Stream(remotecommand.StreamOptions{
        Stdin:             ptyHandler,
        Stdout:            ptyHandler,
        Stderr:            ptyHandler,
        TerminalSizeQueue: ptyHandler,
        Tty:               true,
    })
```

# websocket
github.com/gorilla/websocket 是 go 的一个 websocket 实现，提供了全面的 websocket 相关的方法，这里使用它来实现上面所说的PtyHandler接口。
首先定义一个 TerminalSession 类，该类包含一个 *websocket.Conn，通过 websocket 连接实现PtyHandler接口的读写方法，Next 方法在 remotecommand 执行过程中会被调用。

# file browser

- 从本地拷贝到容器内
tar file ---> io.Writer ---> kubernetes api ---> io.Reader ---> unTar file
- 从容器内拷贝到本地
tar file ---> io.Writer ---> kubernetes api ---> io.Reader ---> unTar file
```




# 参考
https://www.yfdou.com/archives/kuberneteszhi-kubectlexeczhi-ling-gong-zuo-yuan-li-shi-xian-copyhe-webshellyi-ji-filebrowser.html

https://github.com/fangfenghuang/webssh

