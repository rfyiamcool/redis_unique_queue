package unique_queue

import (
	"github.com/garyburd/redigo/redis"
)

const (
	SCRIPT_PUSH = `
local q = KEYS[1]
local q_set = KEYS[1] .. "_set"
local v = redis.call("SADD", q_set, ARGV[1])
if v == 1
then
	return redis.call("RPUSH", q, ARGV[1]) and 1
else
	return 0
end
`

	SCRIPT_POP = `
local q = KEYS[1]
local q_set = KEYS[1] .. "_set"
local v = redis.call("LPOP", q)
if v ~= ""
then
	redis.call("SREM", q_set, v)
end
return v
`
)


type UniqueQueue struct{
	pool *redis.Pool
}

func NewUniqueQueue(r *redis.Pool) *UniqueQueue {
	return &UniqueQueue{
		pool: r,
	}
}

// return 1, push success
// return 0, not push, set filter
func (u *UniqueQueue) UniquePush(q string, body string) (int, error) {
	rc := u.pool.Get()
	defer rc.Close()

	script := redis.NewScript(1, SCRIPT_PUSH)
	resp, err := redis.Int(script.Do(rc, q, body))
	if err == redis.ErrNil {
		err = nil
	}
	return resp, err
}

// redis lua script not support brpop cmd, because brpop is block redis signle thread.
func (u *UniqueQueue) UniquePop(q string) (resp string, err error) {
	rc := u.pool.Get()
	defer rc.Close()

	script := redis.NewScript(1, SCRIPT_POP)
	resp, err = redis.String(script.Do(rc, q))
	if err == redis.ErrNil {
		err = nil
	}
	return resp, err
}

func (u *UniqueQueue) Length(q string) (resp int, err error) {
	rc := u.pool.Get()
	defer rc.Close()

	resp, err = redis.Int(rc.Do("LLEN", q))
	return resp, err
}

func (u *UniqueQueue) Clear(q string) (resp int, err error) {
	rc := u.pool.Get()
	defer rc.Close()

	resp, err = redis.Int(rc.Do("DEL", q))
	if err != nil{
		return resp, err
	}
	resp, err = redis.Int(rc.Do("DEL", u.getQueueSet(q)))
	return resp, err
}

func (u *UniqueQueue) getQueueSet(q string) (string) {
	return q + "_set"
}
