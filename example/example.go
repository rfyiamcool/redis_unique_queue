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
