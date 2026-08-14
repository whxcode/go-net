package redis

import (
	"fmt"
	"time"
)

func SaveOfflineMessage(userID uint, msg []byte) error {
	key := fmt.Sprintf("offlineMessage:%d", userID)

	if err := RedisClient.RPush(Ctx, key, msg).Err(); err != nil {
		return err
	}

	return RedisClient.Expire(Ctx, key, 7*24*time.Hour).Err() // 设置过期时间为 0，表示永不过期
}

func GetOfflineMessage(userID uint) ([][]byte, error) {
	key := fmt.Sprintf("offlineMessage:%d", userID)

	result, err := RedisClient.LRange(Ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	RedisClient.Del(Ctx, key) // 获取后删除离线消息

	messages := make([][]byte, len(result))

	for i, v := range result {
		messages[i] = []byte(v)
	}

	return messages, nil
}
