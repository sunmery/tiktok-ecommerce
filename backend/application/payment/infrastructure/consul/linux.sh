#!/usr/bin/env bash
# 启用 POSIX 模式并设置严格的错误处理机制
set -o posix errexit -o pipefail

# 8300 TCP协议，用于Consul集群中各个节点相互连结通信的端口
# 8301 TCP或者UDP协议，用于Consul节点之间相互使用Gossip协议健康检查等交互
# 8302 TCP或者UDP协议，用于单个或多个数据中心之间的服务器节点的信息同步
# 8500 HTTP协议，用于API接口或者我们上述的网页管理界面访问
# 8600 TCP或者UDP协议，作为DNS服务器，用于通过节点名查询节点信息
#
# 所以如果是在服务器上面部署，记得配置好防火墙放行上述端口。在Spring Cloud模块集成Consul服务发现时，需要配置8500端口。
# 除此之外，我们来看一下命令最后的几个参数：
#
# agent 表示启动一个Agent进程
# -server 表示该节点类型为Server节点（下面会讲解集群中的节点类型）
# -ui 开启网页可视化管理界面
# -node 指定该节点名称，注意每个节点的名称必须唯一不能重复！上面指定了第一台服务器节点的名称为n1，那么别的节点就得用其它名称
# -bootstrap-expect 最少集群的Server节点数量，少于这个值则集群失效，这个选项必须指定，由于这里是单机部署，因此设定为1即可
# -advertise 这里要指定本节点外网地址，用于在集群时告诉其它节点自己的地址，如果是在自己电脑上或者是内网搭建单节点/集群则不需要带上这个参数
# -client 指定可以外部连接的地址，0.0.0.0表示外网全部可以连接
#
# 除此之外，还可以加上-datacenter参数自定义一个数据中心名，同一个数据中心的节点数据中心名应当指定为一样！
mkdir -p /home/docker/consul/data
rm -rf /home/docker/consul/data/*
docker stop consul
docker rm consul
docker run -itd \
 --name=consul \
 -p 8300:8300 \
 -p 8301:8301 \
 -p 8302:8302 \
 -p 8500:8500 \
 -p 8600:8600 \
 -v /home/docker/consul/data:/consul/data \
 hashicorp/consul agent -server -ui -node=n1 -bootstrap-expect=1 -client=0.0.0.0
docker ps
docker logs -f consul
