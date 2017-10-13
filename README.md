# redis_unique_queue

使用redis lua script 操作list + set数据结构, 构建redis的去重队列, 这样既能保证FIFO，又能保证去重.


`example:`

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


	u := unique_queue.NewUniqueQueue(redis_client)
	for i := 0; i < 100; i++ {
		u.UniquePush("qqq", "body...")
		u.UniquePop("qqq")
	}
	fmt.Println("end")
}
```
