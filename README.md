# 分布式缓存

学习来源：https://geektutu.com/post/geecache.html

任务介绍
GeeCache 基本上模仿了 groupcache 的实现，为了将代码量限制在 500 行左右（groupcache 约 3000 行），裁剪了部分功能。 但总体实现上，还是与 groupcache 非常接近的。支持特性有：

* 单机缓存和基于 HTTP 的分布式缓存

* 最近最少访问(Least Recently Used, LRU) 缓存策略

* 使用 Go 锁机制防止缓存击穿

* 使用一致性哈希选择节点，实现负载均衡

* 使用 protobuf 优化节点间二进制通信

* …

[相关笔记](https://gezi-lzq.github.io/wiki/#%E5%AE%9E%E7%8E%B0%E5%88%86%E5%B8%83%E5%BC%8F%E7%BC%93%E5%AD%98:%E5%AE%9E%E7%8E%B0%E5%88%86%E5%B8%83%E5%BC%8F%E7%BC%93%E5%AD%98%20Index%20%E7%AC%94%E8%AE%B0)

![绘制](/assert/design.png)
