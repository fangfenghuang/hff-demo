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
[root@localhost hff]# cat /etc/kata-containers/configuration.toml | grep enable_debug
#enable_debug = true
#enable_debug = true
#enable_debug = true
[root@localhost hff]# cat /etc/kata-containers/configuration.toml | grep kernel_params
kernel_params = ""


$ sudo sed -i -e 's/^# *\(enable_debug\).*=.*$/\1 = true/g' /etc/kata-containers/configuration.toml
$ sudo sed -i -e 's/^kernel_params = "\(.*\)"/kernel_params = "\1 agent.log=debug"/g' /etc/kata-containers/configuration.toml

[root@localhost hff]# cat /etc/kata-containers/configuration.toml | grep enable_debug
enable_debug = true
enable_debug = true
enable_debug = true
[root@localhost hff]# cat /etc/kata-containers/configuration.toml | grep kernel_params
kernel_params = " agent.log=debug"
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



# 开启debug的开销
开启前，已开启了debug_console_enabled，好像并不会因为pod增加而增加？
```bash
[root@localhost hff]# systemd-cgtop | grep kata
/kata_overhead                                                                                                                -      -     2.4M        -        -
/kubepods/podfa152857-05d1-44fc-9cdc-d448b2c98941/kata_f40286f09e1ef5de468d894f1519c7ca6d30962653e7dce8daf90681802a0dde       7      -   167.8M        -        -
/system.slice/kata-monitor.service                                                                                            1      -    22.2M        -        -
```

按上述开启 debug log设置后
```bash
[root@localhost hff]# systemd-cgtop | grep kata
/kata_overhead                                                                                                                -      -     2.4M        -        -
/kubepods/podfa152857-05d1-44fc-9cdc-d448b2c98941/kata_f40286f09e1ef5de468d894f1519c7ca6d30962653e7dce8daf90681802a0dde       7      -   166.6M        -        -
/system.slice/kata-monitor.service                                                                                            1      -    21.8M        -        -
```
 dmesg日志
```bash
[root@localhost hff]# dmesg | grep kata
[2056734.849647] containerd-shim cpuset=kata_5051ee8a1c623152246de3245c514de818127663500f71c01ae4f2952dbdc73a mems_allowed=0
[2056734.849679] Task in /kubepods/burstable/podf4bff02f-469b-4fbb-8152-5daafbe2cb3a/kata_5051ee8a1c623152246de3245c514de818127663500f71c01ae4f2952dbdc73a killed as a result of limit of /kubepods/burstable/podf4bff02f-469b-4fbb-8152-5daafbe2cb3a
[2056734.849690] Memory cgroup stats for /kubepods/burstable/podf4bff02f-469b-4fbb-8152-5daafbe2cb3a/kata_5051ee8a1c623152246de3245c514de818127663500f71c01ae4f2952dbdc73a: cache:188900KB rss:15900KB rss_huge:0KB mapped_file:188876KB swap:0KB inactive_anon:188372KB active_anon:16428KB inactive_file:0KB active_file:0KB unevictable:0KB
[2056751.353176] pool cpuset=kata_d888b88f19234d8148f43414242b99ec4de2d25e881d0453a920622da946c2e1 mems_allowed=0
[2056751.353223] Task in /kubepods/burstable/podf4bff02f-469b-4fbb-8152-5daafbe2cb3a/kata_d888b88f19234d8148f43414242b99ec4de2d25e881d0453a920622da946c2e1 killed as a result of limit of /kubepods/burstable/podf4bff02f-469b-4fbb-8152-5daafbe2cb3a

```