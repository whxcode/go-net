// model/redis.go
package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		panic("Redis 连接失败: " + err.Error())
	}
	println("Redis 连接成功")
}
