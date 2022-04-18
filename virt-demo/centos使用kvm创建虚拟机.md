https://blog.csdn.net/dhRainer/article/details/83411555

## 安装完KVM之后，需要配置一下网卡，增加一个桥接网卡：

```bash
cd /etc/sysconfig/network-scripts/
cp ifcfg-eno16777728 ifcfg-br0  # 拷贝当前的网卡文件
vim ifcfg-eno16777728  # 修改文件内容如下
TYPE=Ethernet
BOOTPROTO=dhcp
DEFROUTE=yes
PEERDNS=yes
PEERROUTES=yes
IPV4_FAILURE_FATAL=no
IPV6INIT=yes
IPV6_AUTOCONF=yes
IPV6_DEFROUTE=yes
IPV6_PEERDNS=yes
IPV6_PEERROUTES=yes
IPV6_FAILURE_FATAL=no
NAME=eno16777728
DEVICE=eno16777728
ONBOOT=yes
BRIDGE=br0

vim ifcfg-br0  # 修改文件内容如下
TYPE=Bridge
BOOTPROTO=dhcp
DEFROUTE=yes
PEERDNS=yes
PEERROUTES=yes
IPV4_FAILURE_FATAL=no
IPV6INIT=yes
IPV6_AUTOCONF=yes
IPV6_DEFROUTE=yes
IPV6_PEERDNS=yes
IPV6_PEERROUTES=yes
IPV6_FAILURE_FATAL=no
NAME=br0
DEVICE=br0
ONBOOT=yes

systemctl restart network  # 重启服务

systemctl start libvirtd  # 启动libvirtd服务

brctl show  # 可以看到两个网卡
```

## 创建卷
```bash
cd /root/hff/kvm
qemu-img create -f qcow2    hfftest.qcow2  10G
qemu-img info  hfftest.qcow2 
image: hfftest.qcow2
file format: qcow2
virtual size: 10G (10737418240 bytes)
disk size: 196K
cluster_size: 65536
Format specific information:
    compat: 1.1
    lazy refcounts: false
```

## 准备一个操作系统的镜像文件
a)命令行安装虚拟机
```bash
virt-install \
--virt-type=kvm \
--name=hfftest \
--vcpus=2 \
--memory=2048 \
--location=/hff/kvm/iso/CentOS-7-x86_64-Minimal-2009.iso \
--disk path=/hff/kvm/images/hfftest.qcow2,size=15,format=qcow2 \
--network bridge=virbr0 \
--graphics none \
--extra-args='console=ttyS0' \
--force
```
>命令说明：
--name 指定虚拟机的名称
--memory 指定分配给虚拟机的内存资源大小
maxmemory 指定可调节的最大内存资源大小，因为KVM支持热调整虚拟机的资源
--vcpus 指定分配给虚拟机的CPU核心数量
maxvcpus 指定可调节的最大CPU核心数量
--os-type 指定虚拟机安装的操作系统类型
--os-variant 指定系统的发行版本
--location 指定ISO镜像文件所在的路径，支持使用网络资源路径，也就是说可以使用URL
--disk path 指定虚拟硬盘所存放的路径及名称，size 则是指定该硬盘的可用大小，单位是G
--bridge 指定使用哪一个桥接网卡，也就是说使用桥接的网络模式
--graphics 指定是否开启图形
--console 定义终端的属性，target_type 则是定义终端的类型
--extra-args 定义终端额外的参数

查询虚拟机所使用的vnc端口virsh vncdisplay centos


## 虚拟机操作
```bash
[root@localhost ~]# virsh console study01  # 进入指定的虚拟机，进入的时候还需要按一下回车
[root@localhost ~]# virsh start study01  # 启动虚拟机
[root@localhost ~]# virsh shutdown study01  # 关闭虚拟机
[root@localhost ~]# virsh destroy study01  # 强制停止虚拟机
[root@localhost ~]# virsh undefine study01  # 彻底销毁虚拟机，会删除虚拟机配置文件，但不会删除虚拟磁盘
[root@localhost ~]# virsh autostart study01  # 设置宿主机开机时该虚拟机也开机
[root@localhost ~]# virsh autostart --disable study01  # 解除开机启动
[root@localhost ~]# virsh suspend study01 # 挂起虚拟机
[root@localhost ~]# virsh resume study01 # 恢复挂起的虚拟机
```

## 安装完成，配置固定IP
```bash
[root@localhost network-scripts]# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: ens3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 52:54:00:40:a4:ef brd ff:ff:ff:ff:ff:ff
    inet 192.168.122.151/24 brd 192.168.122.255 scope global noprefixroute dynamic ens3
       valid_lft 3598sec preferred_lft 3598sec
    inet6 fe80::f5ed:53d0:e383:59e6/64 scope link noprefixroute
       valid_lft forever preferred_lft forever
[root@localhost network-scripts]#
```
