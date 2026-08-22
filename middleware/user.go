package middleware

import (
	"net/http"
	"strconv"

	"go-net/model"
	"go-net/redis"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := utils.GetToken(c)

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MakeResponseWidthCode("未登录", http.StatusUnauthorized))
			return
		}

		userIdStr, err := redis.User.GetToken(c, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MakeResponseWidthCode("登录已过期", http.StatusUnauthorized))
			return
		}

		userId, err := strconv.ParseUint(userIdStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MakeResponseWidthCode("无效的用户ID", http.StatusUnauthorized))
			return
		}

		utils.SetUserID(c, model.UserID(userId))
		c.Set("token", token)

		c.Next()
	}
}
