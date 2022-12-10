# 分布式缓存

学习来源：https://geektutu.com/post/geerpc.html

任务介绍 GeeRPC 选择从零实现 Go 语言官方的标准库 net/rpc，并在此基础上，新增了协议交换(protocol exchange)、注册中心(registry)、服务发现(service discovery)、负载均衡(load balance)、超时处理(timeout processing)等特性。

对应任务：

- 服务端与消息编码
- 支持并发与异步的客户端
- 服务注册
- 超时处理
- 支持 HTTP 协议
- 负载均衡
- 服务发现与注册中心

[相关笔记](https://gezi-lzq.github.io/wiki/#%E5%AE%9E%E7%8E%B0%E5%88%86%E5%B8%83%E5%BC%8F%E7%BC%93%E5%AD%98:%E5%AE%9E%E7%8E%B0%E5%88%86%E5%B8%83%E5%BC%8F%E7%BC%93%E5%AD%98%20Index%20%E7%AC%94%E8%AE%B0)
