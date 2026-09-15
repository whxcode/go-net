package controller

import (
	RediusDB "go-net/redis"

	"github.com/gin-gonic/gin"
)

func GetRedis(c *gin.Context) {
	// 1. 获取所有 key
	keys, err := RediusDB.RedisClient.Keys(c, "*").Result()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 2. 遍历获取每个 key 的值
	result := make(map[string]interface{})
	for _, key := range keys {
		// 获取值（字符串类型）
		val, err := RediusDB.RedisClient.Get(c, key).Result()
		if err != nil {
			// 如果 key 不存在或类型不对，跳过
			continue
		}
		result[key] = val
	}

	// 3. 返回 JSON
	c.JSON(200, gin.H{
		"total": len(result),
		"data":  result,
	})
}
