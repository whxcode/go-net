package redis

import (
	"fmt"
	"net/http"
	"time"

	"go-net/model"
	"go-net/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
	if err := RedisClient.Set(Ctx, token, userID, 7*24*time.Hour).Err(); err != nil {
		panic(err)
	}

	return nil
}

/**
* token: userID
*
* */
func (u *user) GetToken(c *gin.Context, token string) (string, error) {
	r, err := RedisClient.Get(c, token).Result()
	if err != nil {
		if err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MakeResponseWidthCode("登录已过期", http.StatusUnauthorized))
			return "", err
		}

		panic("GetToken error" + err.Error())
	}

	return r, nil
}

func (u *user) DelToken(token string) error {
	err := RedisClient.Del(Ctx, token).Err()
	if err != nil {
		panic(err)
	}
	return nil
}
