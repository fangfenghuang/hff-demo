 [TOC]

## journalctl -t kata

Kata containerd shimv2 运行时日志通过containerd，其日志将被发送到containerd日志指向的任何地方。

查看shimv2运行时日志：

> $ sudo journalctl -t kata

## journalctl -t containerd



# debug

[https://github.com/kata-containers/kata-containers/blob/main/docs/Developer-Guide.md#troubleshoot-kata-containers](https://github.com/kata-containers/kata-containers/blob/main/docs/Developer-Guide.md#troubleshoot-kata-containers)

[https://github.com/kata-containers/kata-containers/blob/main/docs/Developer-Guide.md#connect-to-debug-console](https://github.com/kata-containers/kata-containers/blob/main/docs/Developer-Guide.md#connect-to-debug-console)
（进入虚拟机需要打开debug_console_enabled）

# tracing

[https://github.com/kata-containers/kata-containers/blob/main/docs/tracing.md](https://github.com/kata-containers/kata-containers/blob/main/docs/tracing.md)

# kata-log-parser

[https://github.com/kata-containers/tests/tree/main/cmd/log-parser](https://github.com/kata-containers/tests/tree/main/cmd/log-parser)

# kata-collect-data.sh

向社区提issue需要添加采集信息

