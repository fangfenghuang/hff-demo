
[TOC]
# sysbench
ctr -n k8s.io run  -t --rm docker.io/ljishen/sysbench:latest sysbench cpu run


# fio
ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/xridge/fio:latest hfftest sh

