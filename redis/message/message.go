package MessageRedis

import (
	"fmt"
	"time"

	RediusDB "go-net/redis"

	"github.com/gin-gonic/gin"
)

func SaveOfflineMessage(c *gin.Context, userID uint, msg []byte) {
	key := fmt.Sprintf("offlineMessage:%d", userID)

	if err := RediusDB.RedisClient.RPush(c.Request.Context(), key, msg).Err(); err != nil {
		panic(err)
	}

	// 设置过期时间为 0，表示永不过期
	if err := RediusDB.RedisClient.Expire(c.Request.Context(), key, 7*24*time.Hour).Err(); err != nil {
		panic(err)
	}
}

func GetOfflineMessage(c *gin.Context, userID uint) [][]byte {
	key := fmt.Sprintf("offlineMessage:%d", userID)

	result, err := RediusDB.RedisClient.LRange(c.Request.Context(), key, 0, -1).Result()
	if err != nil {
		return nil
	}

	RediusDB.RedisClient.Del(c.Request.Context(), key) // 获取后删除离线消息

	messages := make([][]byte, len(result))

	for i, v := range result {
		messages[i] = []byte(v)
	}

	return messages
}
