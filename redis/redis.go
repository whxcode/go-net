package RediusDB

import (
	"context"

	"go-net/config"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	_Ctx        = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.ConfigData.Redis.Addr,
	})

	if err := RedisClient.Ping(_Ctx).Err(); err != nil {
		panic("Redis 连接失败: " + err.Error())
	}
	println("Redis 连接成功")
}
