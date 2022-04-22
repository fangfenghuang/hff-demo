
[TOC]
# sysbench
ctr -n k8s.io run  -t --rm docker.io/ljishen/sysbench:latest sysbench cpu run

ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/dotnetdr/sysbench:0.5 hfftest sh

# fio
ctr -n k8s.io run --runtime io.containerd.kata.v2 -t --rm docker.io/xridge/fio:latest hfftest sh
