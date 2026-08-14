package middleware

import (
	"net/http"
	"strconv"

	"go-net/model/redis"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MakeResponseWidthCode("未登录", http.StatusUnauthorized))
			return
		}

		userIdStr, err := redis.RedisClient.Get(c, token).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MakeResponseWidthCode("登录已过期", http.StatusUnauthorized))
			return
		}

		userId, err := strconv.ParseUint(userIdStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MakeResponseWidthCode("无效的用户ID", http.StatusUnauthorized))
			return
		}

		c.Set("userId", uint(userId))
		c.Set("token", token)

		c.Next()
	}
}
