[TOC]

裸机 vs runc容器 vs kata容器磁盘IO对比

# 测试环境与配置：

说明：
- 限制容器request/limit 1C2G
- kata设置debug_console_enabled=true（虚拟机开销占用业务开销）
- kata设置debug_console_enabled=false（虚拟机开销不限制）

# dd 
使用dd对磁盘做性能测试：(4k块大小1G文件)
```bash
dd if=/dev/zero of=/test/host.txt bs=4096 count=1024000 conv=fsync oflag=direct
```
|      | count=1024000 | count=10240000 | 
| ---- | ------------- | -------------- | 
| host | 128~129MB/s   | 111 MB/s       | 
| runc | 126~129MB/s   | 114.9MB/s      | 
| kata | 130~135 MB/s  | 128.0MB/s      | 


# fio 
使用fio 对磁盘做性能测试：(4k块大小1G文件)
```bash

```

## 测试过程及测试结果



## 测试数据分析
|          |裸机|runc|kata（true）|kata（false）|
|----------|-----------|-------------|-------------|-------------|
|顺序读     | 
|随机写     | 
|顺序写     |  
|混合随机读写| 

