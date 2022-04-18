
(1) 移出该节点：etcdctl member remove xx（可以通过member list获取XX）
(2) 停止etcd服务
(3) 需要将配置中的 cluster_state改为：existing，因为是加入已有集群，不能用 new
修复机器问题，删除旧的数据目录，重新启动 etcd 服务
(4) 加入 memeber： etcdctl member add xxx –peer-urls=https://x.x.x.x:2380
(5) 验证：etcdctl endpoint status

