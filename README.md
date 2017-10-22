# redis_unique_queue

## 功能:

使用redis lua script 封装的去重及优先级队列方法, 达到了组合命令的原子性和节省来往的io请求的目的.

`去重队列:`

不仅能保证FIFO, 而且去重.

`优先级去重队列:`

按照优先级获取任务, 并且去重.


## Usage:

PriorityQueue

```
NewPriorityQueue(priority int, unique bool, r *redis.Pool)

Push(q string, body string, pri int) (int, error)

Pop(q string) (resp string, err error)
```

UniqueQueue

```
NewUniqueQueue(r *redis.Pool) *UniqueQueue

UniquePush(q string, body string) (int, error)

UniquePop(q string) (resp string, err error)

more..
```

`to do:`

* 加入批量操作

`example: 优先级+去重队列`

```
package main

// xiaorui.cc

import (
	"fmt"

	"github.com/rfyiamcool/redis_unique_queue"
)

func main() {
	fmt.Println("start")
	redis_client_config := unique_queue.RedisConfType{
		RedisPw:          "",
		RedisHost:        "127.0.0.1:6379",
		RedisDb:          0,
		RedisMaxActive:   100,
		RedisMaxIdle:     100,
		RedisIdleTimeOut: 1000,
	}
	redis_client := unique_queue.NewRedisPool(redis_client_config)

	qname := "xiaorui.cc"
	body := "message from xiaorui.cc"

	u := unique_queue.NewPriorityQueue(3, true, redis_client)
	// 3: 3个优先级,从1-3级
	// true: 开启unique set

	u.Push(qname, body, 2)
	// 2, 优先级

	fmt.Println(u.Pop(qname))
}
```

`example: 去重队列`

```
package main

import (
	"fmt"

	"github.com/rfyiamcool/redis_unique_queue"
)

func main() {
	fmt.Println("start")
	redis_client_config := unique_queue.RedisConfType{
		RedisPw:          "",
		RedisHost:        "127.0.0.1:6379",
		RedisDb:          0,
		RedisMaxActive:   100,
		RedisMaxIdle:     100,
		RedisIdleTimeOut: 1000,
	}
	redis_client := unique_queue.NewRedisPool(redis_client_config)


	qname := "xiaorui.cc"
	u := unique_queue.NewUniqueQueue(redis_client)
	for i := 0; i < 100; i++ {
		u.UniquePush(qname, "body...")
	}

	fmt.Println(u.Length(qname))

	for i := 0; i < 100; i++ {
		u.UniquePop(qname)
	}

	fmt.Println(u.Length(qname))


	fmt.Println("end")
}
```

