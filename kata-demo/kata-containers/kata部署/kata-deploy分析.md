 #  daemonset：
 ```yaml
image: quay.io/kata-containers/kata-deploy:latest
  command:
       - bash
       - '-c'
       - /opt/kata-artifacts/scripts/kata-deploy.shinstall
   lifecycle:
       preStop:
          exec:
              command:
                 - bash
                 - '-c'
                - /opt/kata-artifacts/scriptskata-deploy.sh cleanup
   securityContext:
       privileged: false
  volumeMounts:
  -  name: systemd
      mountPath: /run/systemd
```

# 部署kata-deploy之后
```bash
[root@localhost ~]# ps -ef | grep kata
root     16533 16031  0 14:26 pts/1    00:00:00 grep --color=auto kata
root     24526 20728  0 11:09 ?        00:00:00 bash /opt/kata-artifacts/scripts/kata-deploy.sh install
[root@localhost ~]# ps -ef | grep qemu
root     18834 16031  0 14:27 pts/1    00:00:00 grep --color=auto qemu
[root@localhost ~]#
[root@localhost ~]#
[root@localhost ~]# ps -ef | grep kvm
root       756     2  0 Mar25 ?        00:00:00 [kvm-irqfd-clean]
root     18975 16031  0 14:27 pts/1    00:00:00 grep --color=auto kvm
[root@localhost ~]# ls /opt/kata/
bin  libexec  share
[root@localhost ~]# ls /opt/kata/bin/
cloud-hypervisor         firecracker  kata-collect-data.sh  kata-runtime
containerd-shim-kata-v2  jailer       kata-monitor          qemu-system-x86_64
```

## kata-deploy内部
```bash
[root@kata-deploy-6r86s /]# ps -ef 
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 03:09 ?        00:00:00 bash /opt/kata-artifacts/scripts/kata-deploy.sh install
root        80     1  0 03:09 ?        00:00:00 sleep infinity
root        81     0  0 06:41 pts/0    00:00:00 bash
root        96    81  0 06:42 pts/0    00:00:00 ps -ef
挂载/run/system后可以控制宿主机服务
[root@kata-deploy-6r86s /]# systemctl status containerd
● containerd.service - containerd container runtime
   Loaded: loaded (/usr/lib/systemd/system/containerd.service; enabled; vendor preset: disabled)
   Active: active (running) since Thu 2022-04-07 06:47:47 UTC; 6min ago
     Docs: https://containerd.io
  Process: 3033 ExecStartPre=/sbin/modprobe overlay (code=exited, status=0/SUCCESS)
 Main PID: 3035
    Tasks: 585
   Memory: 535.9M
   CGroup: /system.slice/containerd.service
```

# 部署后containerd配置追加
```bash
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
  runtime_type = "io.containerd.kata.v2"
  privileged_without_host_devices = true
  pod_annotations = ["io.katacontainers.*"]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata.options]
    ConfigPath = "/opt/kata/share/defaults/kata-containers/configuration.toml"
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-fc]
  runtime_type = "io.containerd.kata-fc.v2"
  privileged_without_host_devices = true
  pod_annotations = ["io.katacontainers.*"]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-fc.options]
    ConfigPath = "/opt/kata/share/defaults/kata-containers/configuration-fc.toml"
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-qemu]
  runtime_type = "io.containerd.kata-qemu.v2"
  privileged_without_host_devices = true
  pod_annotations = ["io.katacontainers.*"]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-qemu.options]
    ConfigPath = "/opt/kata/share/defaults/kata-containers/configuration-qemu.toml"
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-clh]
  runtime_type = "io.containerd.kata-clh.v2"
  privileged_without_host_devices = true
  pod_annotations = ["io.katacontainers.*"]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-clh.options]
    ConfigPath = "/opt/kata/share/defaults/kata-containers/configuration-clh.toml"
```

# 部署日志：
```bash
copying kata artifacts onto host
warning: containerd-shim-kata-fc-v2 already exists
#!/usr/bin/env bash
KATA_CONF_FILE=/opt/kata/share/defaults/kata-containers/configuration-fc.toml /opt/kata/bin/containerd-shim-kata-v2 "$@"
warning: containerd-shim-kata-qemu-v2 already exists
#!/usr/bin/env bash
KATA_CONF_FILE=/opt/kata/share/defaults/kata-containers/configuration-qemu.toml /opt/kata/bin/containerd-shim-kata-v2 "$@"
Creating the default shim-v2 binary
warning: containerd-shim-kata-clh-v2 already exists
#!/usr/bin/env bash
KATA_CONF_FILE=/opt/kata/share/defaults/kata-containers/configuration-clh.toml /opt/kata/bin/containerd-shim-kata-v2 "$@"
Add Kata Containers as a supported runtime for containerd
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
  runtime_type = "io.containerd.kata.v2"
  privileged_without_host_devices = true
  pod_annotations = ["io.katacontainers.*"]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata.options]
    ConfigPath = "/opt/kata/share/defaults/kata-containers/configuration.toml"
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-fc]
  runtime_type = "io.containerd.kata-fc.v2"
  privileged_without_host_devices = true
  pod_annotations = ["io.katacontainers.*"]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-fc.options]
    ConfigPath = "/opt/kata/share/defaults/kata-containers/configuration-fc.toml"
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-qemu]
  runtime_type = "io.containerd.kata-qemu.v2"
  privileged_without_host_devices = true
  pod_annotations = ["io.katacontainers.*"]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-qemu.options]
    ConfigPath = "/opt/kata/share/defaults/kata-containers/configuration-qemu.toml"
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-clh]
  runtime_type = "io.containerd.kata-clh.v2"
  privileged_without_host_devices = true
  pod_annotations = ["io.katacontainers.*"]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata-clh.options]
    ConfigPath = "/opt/kata/share/defaults/kata-containers/configuration-clh.toml"
node/localhost.localdomain unlabeled
```
