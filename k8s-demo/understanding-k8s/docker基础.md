[TOC]

补充一些docker基础知识点


# 虚悬镜像
镜像列表中，可能会看到一个特殊的镜像，这个镜像既没有仓库名，也没有标签，均为
```bash
<none>               <none>              00285df0df87        5 days ago          342 MB
```
原因有可能是官方镜像维护，发布了新版本后，原镜像名被转移到了新下载的镜像身上，而旧的镜像上的这个名称则被取消，从而成为了 <none>。除了 docker pull 可能导致这种情况，docker build 也同样可以导致这种现象。由于新旧镜像同名，旧镜像名称被取消，从而出现仓库名、标签均为 <none> 的镜像。这类无标签镜像也被称为 虚悬镜像(dangling image) ，可以用下面的命令专门显示这类镜像：
```bash
docker image ls -f dangling=true
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
<none>              <none>              00285df0df87        5 days ago          342 MB
```

一般来说，虚悬镜像已经失去了存在的价值，是可以随意删除的，可以用下面的命令删除。
```bash
docker image prune
```
??实测无法全部删完


# 中间层镜像
默认的 docker image ls 列表中只会显示顶层镜像，如果希望显示包括中间层镜像在内的所有镜像的话，需要加 -a 参数。
```
docker image ls -a
```

# 一些命令技巧
- 只输出镜像ID和仓库名：
```
$ docker image ls --format "{{.ID}}:{{.Repository}}"
5f515359c7f8: redis
05a60462f8ba: nginx
fe9198c04d62: mongo
00285df0df87: <none>
329ed837d508: ubuntu
329ed837d508: ubuntu
```
- 只输出id列
```
docker image ls -q
```

- 删除所有仓库名为 redis 的镜像
```bash
$ docker image rm $(docker image ls -q redis)
```

