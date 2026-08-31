package redis

import (
	"time"

	"go-net/model"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

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
