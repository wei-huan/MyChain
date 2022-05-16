# MyChain

## 介绍

支持 2 个节点的跨域组网，可以正常发布区块并同步

## 说明

p2p_node 里是区块链节点运行的代码

peer_server 里是打洞服务器运行的代码

## 使用

推荐使用最新的 Golang 版本, 最低要求 Go 1.17

使用时需要关闭防火墙或者开放端口



### 编译

go build


### 运行

#### 第一个节点

`./MyChain -server=:10024`



#### 第二个节点

`./MyChain -server=:10025`



#### http 请求节点1发布区块

`curl -H 'content-type: application/json' -X POST -d '{"name":"luda","data":"first block"}' http://113.54.228.13:10024/blockchain/create`



#### http 请求节点2发布区块

`curl -H 'content-type: application/json' -X POST -d '{"name":"luda","data":"second block"}' http://113.54.228.13:10025/blockchain/create`



#### 查看节点1的链

`curl -H 'content-type: application/json' -X POST -d '{"chain":true, "peer":true}' http://113.54.228.13:10024/blockchain/show >> node1.json`



#### 查看节点2的链

`curl -H 'content-type: application/json' -X POST -d '{"chain":true, "peer":true}' http://113.54.228.13:10025/blockchain/show >> node2.json`
