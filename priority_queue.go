package unique_queue

import (
	"github.com/garyburd/redigo/redis"
)

const (
	DEFAULT_PRIORITY = 1

	SCRIPT_PRIORITY_PUSH = `
local q = KEYS[1]
local q_set = KEYS[1] .. "_set"
local v = redis.call("SADD", q_set, ARGV[1])
if v == 1
then
	local q_pri = KEYS[1] .. "_" .. ARGV[2]
	return redis.call("RPUSH", q_pri, ARGV[1]) and 1
else
	return 0
end

`
	SCRIPT_PRIORITY_POP = `
local q_set = KEYS[1] .. "_set"
local loop_num = tonumber(ARGV[1])
local i = 1

while (i <= loop_num) do
	local q = KEYS[1] .. "_" ..i
	local v = redis.call("LPOP", q)
	-- when not res, is bool
	if type(v) == "boolean"
	then
		local x = KEYS[1] .. "_" ..i
	else
		redis.call("SREM", q_set, v)
		return v
	end
	i = i + 1
end
`
)


type PriorityQueue struct{
	pool *redis.Pool
	Priority int
	Unique bool
}

func NewPriorityQueue(priority int, unique bool, r *redis.Pool) *PriorityQueue {
	if priority < 0 {
		priority = DEFAULT_PRIORITY
	}
	return &PriorityQueue{
		pool: r,
		Priority: priority,
		Unique: unique,
	}
}

// return 1, push success
// return 0, not push, set filter
func (u *PriorityQueue) Push(q string, body string, priority int) (int, error) {
	rc := u.pool.Get()
	defer rc.Close()

	script := redis.NewScript(1, SCRIPT_PRIORITY_PUSH)
	resp, err := redis.Int(script.Do(rc, q, body, priority))
	if err == redis.ErrNil {
		err = nil
	}
	return resp, err
}

// redis lua script not support brpop cmd, because brpop is block redis signle thread.
func (u *PriorityQueue) Pop(q string) (resp string, err error) {
	rc := u.pool.Get()
	defer rc.Close()

	script := redis.NewScript(1, SCRIPT_PRIORITY_POP)
	resp, err = redis.String(script.Do(rc, q, u.Priority))
	if err == redis.ErrNil {
		err = nil
	}
	return resp, err
}
