# Go-BlockChain 的 P2P 包

适配 [BlockChain-CN](https://github.com/Blockchain-CN) 的 P2P 库

## 项目介绍

以 [BlockChain-CN](https://github.com/Blockchain-CN) 为基础实现跨域，项目本身使用的 [P2P](github.com/Blockchain-CN/pheromones) 库并不支持跨域，即跨局域网连接，只是在单主机上，已知分配的端口的条件下，简单的 TCP 端口连接，而且其本身也不是个对等网络，每个结点的地位不同。现在需要实现多个局域网之间的 P2P 网络，进而实现跨局域网的区块链网络，故重新实现自己的 P2P 库 My_P2P。同时还需要兼容上层[BlockChain-CN](https://github.com/Blockchain-CN)的使用

## 具体工作

重新构造工程，之前整个逻辑写的很烂，居然是浪费这么多地址通信
原来的根本不是 P2P，只是个像 P2P 的垃圾

保留 P2P 功能的协议部分，协议写的还行

## 项目成员

- 凌智
- 陈诚
- 唐昊哲

## TODO

1. 面向对象重构
2. 更换 SHA256 为国标
3. 前端

## 开源协议
MIT License
