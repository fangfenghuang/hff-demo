 [TOC]

## journalctl -t kata

Kata containerd shimv2 运行时日志通过containerd，其日志将被发送到containerd日志指向的任何地方。

查看shimv2运行时日志：

> $ sudo journalctl -t kata

## journalctl -t containerd



# 开启 debug log
开启 debug log 可以帮助我们获得更详细的 log，除了 runtime 的 log，而且还能看到 agent 的 log，以及 guest OS 中 kernel 的 log（dmesg命令的输出）。

开启 debug log，需要修改两个配置文件：

- containerd 配置文件按如下修改即可：
[debug]
  level = "debug"

- Kata Containers需要开启 runtime 和 agent 的 debug log
```bash
sudo sed -i -e 's/^# *\(enable_debug\).*=.*$/\1 = true/g' /etc/kata-containers/configuration.toml
$ sudo sed -i -e 's/^kernel_params = "\(.*\)"/kernel_params = "\1 agent.log=debug"/g' /etc/kata-containers/configuration.toml
```


[https://github.com/kata-containers/kata-containers/blob/main/docs/Developer-Guide.md#troubleshoot-kata-containers](https://github.com/kata-containers/kata-containers/blob/main/docs/Developer-Guide.md#troubleshoot-kata-containers)

[https://github.com/kata-containers/kata-containers/blob/main/docs/Developer-Guide.md#connect-to-debug-console](https://github.com/kata-containers/kata-containers/blob/main/docs/Developer-Guide.md#connect-to-debug-console)
（进入虚拟机需要打开debug_console_enabled）

# tracing

[https://github.com/kata-containers/kata-containers/blob/main/docs/tracing.md](https://github.com/kata-containers/kata-containers/blob/main/docs/tracing.md)

# kata-log-parser

[https://github.com/kata-containers/tests/tree/main/cmd/log-parser](https://github.com/kata-containers/tests/tree/main/cmd/log-parser)

# kata-collect-data.sh

向社区提issue需要添加采集信息

