package redis

import (
	"fmt"
	"time"

	"go-net/model"

	"github.com/gin-gonic/gin"
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

type user struct{}

var User *user = &user{}

/**
* token: userID
*
* */
func (u *user) SetToken(token string, userID model.UserID) error {
	return RedisClient.Set(Ctx, token, userID, 7*24*time.Hour).Err()
}

/**
* token: userID
*
* */
func (u *user) GetToken(c *gin.Context, token string) (string, error) {
	return RedisClient.Get(c, token).Result()
}

func (u *user) DelToken(token string) error {
	return RedisClient.Del(Ctx, token).Err()
}
