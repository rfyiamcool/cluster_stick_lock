# cluster_stick_lock

`一句话, 原子解决分布式锁的续约问题`

使用redis实现分布式锁已经是个很场景的需求，当使用 `setnx key value ex 10 nx` 创建锁之后...如果我们想续下ttl怎么办? 当然你会说用 expire 啊... 你怎么确定这个key是你创建的， 当然你会说 放在value做标记呀...

集群环境下，多个节点都在setnx, 当超时发生时，每一个节点都有可能拿到锁. 另外，get key 和 expire 组合操作会有小概率误操作.

该脚本只是用lua把setnx、get、expire组合起来而已.

```
// params KEYS[1]: key
// params ARGV[1]: value
// params ARGV[2]: 过期秒数
// return 1: 拿到锁.
// return 0: 未拿到锁.
// desc: 为避免分布式锁超时, unlocker飘逸浮动.

if redis.call("set", KEYS[1], ARGV[1], "ex", ARGV[2], "nx")
then
	redis.call("set", "lua_debug", "setnx success")
    return 1
end

if redis.call("get", KEYS[1]) == ARGV[1]
then
	redis.call("expire", KEYS[1], ARGV[2])
	return 1
else
	return 0
end
```
