var (
	// params KEYS[1]: key
	// params ARGV[1]: value
	// params ARGV[2]: 过期秒数
	// return 1: 拿到锁.
	// return 0: 未拿到锁.
	// desc: 为避免分布式锁超时, unlocker飘逸浮动.
	SCRIPTS = `
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
`
)
