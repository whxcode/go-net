package redis

import (
	"fmt"
	"time"
)

type message struct{}

var Message *message = &message{}

func (m *message) SaveOfflineMessage(userID uint, msg []byte) {
	key := fmt.Sprintf("offlineMessage:%d", userID)

	if err := RedisClient.RPush(Ctx, key, msg).Err(); err != nil {
		panic(err)
	}

	// 设置过期时间为 0，表示永不过期
	if err := RedisClient.Expire(Ctx, key, 7*24*time.Hour).Err(); err != nil {
		panic(err)
	}
}

func (m *message) GetOfflineMessage(userID uint) [][]byte {
	key := fmt.Sprintf("offlineMessage:%d", userID)

	result, err := RedisClient.LRange(Ctx, key, 0, -1).Result()
	if err != nil {
		return nil
	}

	RedisClient.Del(Ctx, key) // 获取后删除离线消息

	messages := make([][]byte, len(result))

	for i, v := range result {
		messages[i] = []byte(v)
	}

	return messages
}
