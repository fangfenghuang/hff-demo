[TOC]

# Docker 镜像存储原理
docker支持多种graphDriver，包括vfs、devicemapper、overlay、overlay2、aufs等等


## overlay2

Linux 内核在 4.0 版本对overlay做了很多必要的改进，此时的 OverlayFS 被称之为overlay2
- 要想使用overlay2，Docker 版本必须高于 17.06.02。
- 如果你的操作系统是 RHEL 或 CentOS，Linux 内核版本必须使用 3.10.0-514 或者更高版本，其他 Linux 发行版的内核版本必须高于 4.0（例如 Ubuntu 或 Debian）
- overlay2最好搭配 xfs 文件系统使用，并且使用 xfs 作为底层文件系统时，d_type必须开启，可以使用以下命令验证 d_type 是否开启：
```bash
$ xfs_info /var/lib/docker | grep ftype
naming   =version 2              bsize=4096   ascii-ci=0 ftype=1
```

**verlay2 是这样储存文件的：** overlay2将镜像层和容器层都放在单独的目录，并且有唯一 ID，每一层仅存储发生变化的文件，最终使用联合挂载技术将容器层和镜像层的所有文件统一挂载到容器中，使得容器中看到完整的系统文件。


```bash
[root@localhost ~]# docker info | grep Storage
 Storage Driver: overlay2
```

Overlay驱动的镜像只有两层：一个upper文件系统和一个lower文件系统，分别代表Docker的镜像层和容器层

merged层是lower层和upper层的统一视图

当需要修改一个文件时，使用CoW将文件从只读的lower复制到可写的upper进行修改，结果也保存在upper层。

在Docker中，底下的只读层就是image，可写层就是Container。voerlay2驱动存储位置为/var/lib/docker/（/app/docker）

```bash
[root@tztest kbuser]# ls /app/docker/
builder  buildkit  containerd  containers  image  kubelet  network  overlay2  plugins  runtimes  swarm  tmp  trust  volumes
[root@tztest kbuser]# ls /app/docker/overlay2/7e32a486c4dc7a52d0685b7fb3e4f6fb4eb7b9ae963f0820285c3659e3f361ca
committed  diff  link  lower  work
[root@tztest kbuser]# ls /app/docker/image/overlay2
distribution  imagedb  layerdb  repositories.json

```
- overlay2目录，存储镜像和容器的层信息

- image目录，存储镜像元相关信息

- repositories.json就是存储镜像信息，主要是name和image id的对应，digest和image id的对应



### 新增一个镜像
```bash
[root@tztest kbuser]# docker images | grep web-demo
web-demo_c7cac94                    dev-1648826415118                          df2b573855e0        12 days ago         17.7MB
[root@tztest kbuser]# ls /app/docker/overlay2/d1142002c8d33580ab1213fff5746027a174d763e37885039df98b7108d41eb0
committed  diff  link  lower  work

（可能新增多个镜像层目录文件）


//允许一个容器
[root@tztest kbuser]# docker run --name hfftest -it  df2b573855e0 sh
[root@tztest kbuser]#  docker ps | grep hfftest
1e73601f56d9        df2b573855e0                                          "sh"                     34 seconds ago      Up 32 seconds       80/tcp, 8080/tcp         hfftest


[root@tztest kbuser]# ls /app/docker/overlay2/ -tl
总用量 3900
drwx------ 5 root root  4096 4月  14 16:42 4695d98d682387009f6ab5c84d86599fe35460dae36aea7d00e1561042d630e0
drwx------ 4 root root  4096 4月  14 16:42 4695d98d682387009f6ab5c84d86599fe35460dae36aea7d00e1561042d630e0-init
drwx------ 2 root root 69632 4月  14 16:42 l
drwx------ 4 root root  4096 4月   1 23:22 d1142002c8d33580ab1213fff5746027a174d763e37885039df98b7108d41eb0
//多了两个目录（4695d98...108d41eb0）和一个l目录

l目录是一堆软连接，把一些较短的随机串软连到镜像层的 diff 文件夹下，这样做是为了避免达到mount命令参数的长度限制。


[root@tztest kbuser]# docker inspect df2b573855e0 | grep Dir

"LowerDir": "/app/docker/overlay2/d90eefdee868034a568195793cd1c269664a8e459f773cef20cefe5fc81f9e1b/diff:/app/docker/overlay2/b1789038611adaa88789bc8600ff21ebf28d0f4437ae720bfea03b34d50cc689/diff:/app/docker/overlay2/9c3ce321c132b4e638a2187c58d6af6dff58f74c3f49780df8ee4eb80b0dc5f2/diff:/app/docker/overlay2/198313ba0d321e9efa1205820408780b4edb07027456cfcd89cc660bbfa92953/diff",
"MergedDir": "/app/docker/overlay2/d1142002c8d33580ab1213fff5746027a174d763e37885039df98b7108d41eb0/merged",
"UpperDir": "/app/docker/overlay2/d1142002c8d33580ab1213fff5746027a174d763e37885039df98b7108d41eb0/diff",
"WorkDir": "/app/docker/overlay2/d1142002c8d33580ab1213fff5746027a174d763e37885039df98b7108d41eb0/work"

LowerDir：对应的是容器的只读镜像层，在新生成目录2c...74-init下；
UpperDir：为容器的可读写层，在新生成目录2c...74下；
MergedDir：为底层只读镜像层和上层可读写容器层的统一视图

//写入一个文件：
[root@tztest kbuser]# docker run --name hfftest -it  df2b573855e0 sh
/app # touch hfftest0414
/app # echo 111 > hfftest0414

[root@tztest kbuser]# find / -name  hfftest0414
/app/docker/overlay2/4695d98d682387009f6ab5c84d86599fe35460dae36aea7d00e1561042d630e0/diff/app/hfftest0414
/app/docker/overlay2/4695d98d682387009f6ab5c84d86599fe35460dae36aea7d00e1561042d630e0/merged/app/hfftest0414
//写入文件在UpperDir和MergedDir中

//停止容器后，创建的文件仍然存在（diff,merge不在了）。当容器被删除后，两个新增目录及其相关文件被删除

```

### 写时复制
overlay2 对文件的修改采用的是写时复制的工作机制：
- 第一次修改文件：当我们第一次在容器中修改某个文件时，overlay2 会触发写时复制操作，overlay2 首先从镜像层复制文件到容器层，然后在容器层执行对应的文件修改操作。
- 删除文件或目录：当文件或目录被删除时，overlay2 并不会真正从镜像中删除它，因为镜像层是只读的，overlay2 会创建一个特殊的文件或目录，这种特殊的文件或目录会阻止容器的访问。

## 查看docker镜像、容器、数据卷所占用的空间
```bash
[root@tztest kbuser]# docker system df
TYPE                TOTAL               ACTIVE              SIZE                RECLAIMABLE
Images              117                 51                  18.52GB             11.19GB (60%)
Containers          203                 100                 2.842GB             1.446GB (50%)
Local Volumes       13                  12                  30.42MB             0B (0%)
Build Cache         0                   0                   0B                  0B
```

# 参考资料
https://www.modb.pro/db/127388
