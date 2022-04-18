[TOC]


# docker网络
docker0：Docker服务会在它所在的机器上创建一个名为docker0的网桥。
docker0的IP地址为172.17.0.1，子网掩码为255.255.0.0。


[root@rqy-k8s-2 kbuser]# docker network ls
NETWORK ID          NAME                DRIVER              SCOPE
e29e51f3f12a        bridge              bridge              local
b70c41f7568e        host                host                local
3a890954cce6        none                null                local

Docker容器在启动时，如果没有显式指定加入任何网络，就会默认加入到名为bridge的网络。而这个bridge网络就是基于docker0实现的。
加入host网络的容器，可以实现和Docker daemon守护进程（也就是Docker服务）所在的宿主机网络环境进行直接通信
none网络，则表示容器在启动时不带任何网络设备。


# k8s网络
## docker0:
1. 同一Pod内的容器间通信:
    因为pause容器提供pod内网络共享，所以容器直接可以使用localhost(lo)访问其他容器
2. 各Pod彼此之间的通信(两个pod在一台主机上面, 两个pod分布在不同主机之上)
   1)两个pod在一台主机上面: 通过docker默认的docker网桥互连容器(docker0), ifconfig 查了docker0
   2)两个pod分布在不同主机之上: 通过CNI插件实现，eg: flannel,calico

3. Pod与Service间的通信
   Service分配的ip叫cluster ip是一个虚拟ip（相对固定，除非删除service），这个ip只能在k8s集群内部使用，
   如果service需要对外提供，只能使用Nodeport方式映射到主机上，使用主机的ip和端口对外提供服务。
   节点上面有个kube-proxy进程，这个进程从master apiserver获取信息，感知service和endpoint的创建，然后做两个事：
    1.为每个service 在集群中每个节点上面创建一个随机端口，任何该端口上面的连接会代理到相应的pod
    2.集群中每个节点安装iptables/ipvs规则，用于clusterip + port路由到上一步定义的随机端口上面，
    所以集群中每个node上面都有service的转发规则:iptables -L -n -t filter

## calico网络
### Calico主要由Felix、Orchestrator Plugin、etcd、BIRD和BGP Router Reflector等组件组成。
* Felix：Calico Agent，运行于每个节点。
* Orchestrator Plugi：编排系统（如 Kubernetes 、 OpenStack 等）以将 Calico整合进系统中的插件，例如Kubernetes的CNI。
* etcd：持久存储Calico数据的存储管理系统。
* BIRD：用于分发路由信息的BGP客户端。
* BGP Route Reflector：BGP路由反射器，可选组件，用于较大规模的网络场景。

### BGP: Pod 1跨节点访问Pod 2大致流程如下：
1. 数据包从容器1出到达Veth Pair另一端（宿主机上，以cali前缀开头）；
2. 宿主机根据路由规则，将数据包转发给下一跳（网关）；
3. 到达Node2，根据路由规则将数据包转发给cali设备，从而到达容器2。
这里最核心的“下一跳”路由规则，就是由Calico的Felix进程负责维护的。这些路由规则信息，则是通过BGP Client也就是BIRD组件，使用BGP协议传输而来的。

### Calico Route Reflector模式（RR）（路由器反射）
Calico维护的网络在默认是（Node-to-Node Mesh）全互联模式（mesh）

### IPIP: Pod 1访问Pod 2大致流程如下：
1. 数据包从容器1出到达Veth Pair另一端（宿主机上，以cali前缀开头）；
2. 进入IP隧道设备（tunl0），由Linux内核IPIP驱动封装在宿主机网络的IP包中（新的IP包目的地之是原IP包的下一跳地址，即192.168.31.63），这样，就成了Node1到Node2的数据包；
3. 数据包经过路由器三层转发到Node2；
4. Node2收到数据包后，网络协议栈会使用IPIP驱动进行解包，从中拿到原始IP包；
5. 然后根据路由规则，根据路由规则将数据包转发给cali设备，从而到达容器2。


