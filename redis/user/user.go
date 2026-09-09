package UserRedis

import (
	"encoding/json"
	"time"

	"go-net/model"

	RediusDB "go-net/redis"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const SevenDay time.Duration = 7 * 24 * time.Hour

/**
* token: userID
*
* */
func SetToken(c *gin.Context, token string, userID model.UserID) error {
	if err := RediusDB.RedisClient.Set(c.Request.Context(), token, userID, SevenDay).Err(); err != nil {
		panic(err)
	}

	return nil
}

/**
* token: userID
*
* */
func GetToken(c *gin.Context, token string) (string, error) {
	r, err := RediusDB.RedisClient.Get(c, token).Result()
	if err != nil {
		if err == redis.Nil {
			return "", err
		}

		panic("GetToken error" + err.Error())
	}

	return r, nil
}

func DelToken(c *gin.Context, token string) error {
	err := RediusDB.RedisClient.Del(c.Request.Context(), token).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

/**
* userID -> userInfo
* */
func SetUserInfo(c *gin.Context, userID model.UserID, user *model.User) {
	data, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	err = RediusDB.RedisClient.Set(c.Request.Context(),
		userID.String(),
		data,
		SevenDay).
		Err()
	if err != nil {
		panic(err)
	}
}

/**
* userID -> userInfo
* */
func GetUserInfo(c *gin.Context, userID model.UserID) *model.User {
	r, err := RediusDB.RedisClient.Get(c.Request.Context(), userID.String()).Result()
	if err != nil {
		panic(err)
	}

	var user *model.User // = &model.User{}
	err = json.Unmarshal([]byte(r), &user)
	if err != nil {
		panic(err)
	}

	return user
}
