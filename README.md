# cluster_stick_lock

使用redis实现分布式锁已经是个很场景的需求，当使用< setnx key value ex 10 nx >创建锁之后，如果我们想续下ttl怎么办? 集群环境下，多个节点都在setnx, 当超时发生时，每一个节点都有可能拿到锁. 另外，get key 和 expire 组合操作会有小概率误操作.

