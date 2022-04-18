[TOC]

# qcow2和raw
* raw格式是原始镜像，会直接当作一个块设备给虚拟机来使用，至于文件里面的空洞，则是由宿主机的文件系统来管理的，Linux下的文件系统可以很好的支持空洞的特性，所以，如果你创建了一个100G的raw格式的文件，ls看的时候，可以看到这个文件是100G的，但是用du 来看，这个文件会很小。
* qcow2是kvm支持的磁盘镜像格式，我们创建一个100G的qcow2磁盘之后，无论用ls来看，还是du来看，都是很小的。这说明了，qcow2本身会记录一些内部块分配的信息的。

无论哪种格式，磁盘的利用率来说，都是一样的，因为实际占用的块数量都是一样的。但是raw的虚拟机会比qcow2的虚拟机IO效率高一些，实际测试的时候会比qcow2高25%，这个性能的差异还是不小的，所以追求性能建议选raw。

raw唯一的缺点在于，ls看起来很大，在scp的时候，这会消耗很多的网络IO，而tar这么大的文件，也是很耗时间跟CPU的，解决方法是把raw转换成qcow2的格式，对空间压缩就很大了。而且速度很快。



## 镜像转换
### img转qcow2
#### 源文件disk.img 目标文件disk.qcow2
qemu-img convert -f raw -O qcow2 disk.img disk.qcow2

### qcow2转docker镜像
FROM registry:5500/kubevirt/container-disk-v1alpha
ADD ./disk.qcow2 /disk

#### 1.上传基础镜像：docker load -i container-disk-v1alpha.tar（dockerhub镜像仓下载也可https://hub.docker.com/r/kubevirt/container-disk-v1alpha）
#### 2.源文件disk.qcow2   生成镜像registry:5500/weiyun-kubevirt/test-qcow:v1.0
docker build -t registry:5500/weiyun-kubevirt/test-qcow:v1.0 .


### iso转qcow2/img/docker镜像
参考https://blog.csdn.net/ximenjianxue/article/details/103782462
1）把ISO文件copy到linux的机器上，并确保硬盘有足够的空间
2）用qemu命令创建qcow2镜像磁盘（用于安装suse镜像），例
qemu-img create -f qcow2 /root/hff/iso_2_qcow2/hfftest.qcow2 50G
 
3）启动Kvm，安装操作系统
virt-install --name hfftest --ram 2048 --cdrom=/root/hff/iso_2_qcow2/VCloud-1.1.32-09270826.iso --disk path=/root/hff/iso_2_qcow2/hfftest.qcow2,format=qcow2 --graphics vnc,listen=0.0.0.0 --noautoconsole --os-type=linux --os-variant=rhel7 --check all=off

可能会出现KVM cannot access storage file (as uid:107, gid:107)permission denied问题。
解决方法
https://www.xiexianbin.cn/linux/kvm/2017-03-14-kvm-cannot-access-storage-file-as-uid-107-gid-107-permission-denied/index.html
Changing /etc/libvirt/qemu.conf make working things. Uncomment user/group to work as root.
然后重启服务：
systemctl restart libvirtd.service
 
KVM客户机的配置文件放置在**/etc/libvirt/qemu**下。使用vi可以查看虚拟机的xml配置文件。

3）查看vnc端口
virsh vncdisplay hfftest
 
用vnc去连接安装系统