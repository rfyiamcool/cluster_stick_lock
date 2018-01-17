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
