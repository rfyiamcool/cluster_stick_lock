package primaryCtl

// 使用样例... 自己看吧

import (
	"sync/atomic"
	"time"

	"github.com/garyburd/redigo/redis"

	"task_dispatcher/conf"
	"task_dispatcher/core/mq"
)

const EXPIRE = 10
const SLEEP_INTERVAL = 1
const JUDGE_INTERVAL = 1
const PRIMARY_REDIS_KEY = "primary_scaner_flag"

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

type Primary struct {
	mark   int32
	lastTs time.Time
}

func NewPrimary() *Primary {
	return &Primary{
		mark:   0,
		lastTs: time.Now(),
	}
}

func (p *Primary) setPrimaryMark() {
	atomic.StoreInt32(&p.mark, 1)
}

func (p *Primary) unsetPrimaryMark() {
	atomic.StoreInt32(&p.mark, 0)
}

func (p *Primary) TrySetPrimary() bool {
	rc := mq.RedisEngineClient.Get()
	defer rc.Close()

	js := redis.NewScript(1, SCRIPTS)
	res, err := redis.Int(js.Do(rc, PRIMARY_REDIS_KEY, conf.HOSTNAME, EXPIRE))
	if err != nil {
		p.unsetPrimaryMark()
		return false
	}
	if res == 1 {
		p.setPrimaryMark()
		return true
	}

	p.unsetPrimaryMark()
	return false
}

func (p *Primary) SetNxPrimary() (bool, error) {
	rc := mq.RedisEngineClient.Get()
	defer rc.Close()

	_, err := redis.String(rc.Do("SET", PRIMARY_REDIS_KEY, conf.HOSTNAME, "EX", EXPIRE, "NX"))
	if err == redis.ErrNil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	atomic.StoreInt32(&p.mark, 1)
	return true, nil
}

func (p *Primary) Get() bool {
	if atomic.LoadInt32(&p.mark) == 1 {
		return true
	}
	return false
}

func (p *Primary) Check() bool {
	val := time.Now().Sub(p.lastTs)
	if val.Seconds() >= JUDGE_INTERVAL {
		return true
	}
	return false
}

func (p *Primary) Run() {
	for {
		p.TrySetPrimary()
		time.Sleep(SLEEP_INTERVAL * time.Second)
	}
}

