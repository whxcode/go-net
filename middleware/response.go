package middleware

import (
	"net/http"

	"go-net/response"

	"github.com/gin-gonic/gin"
)

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 执行后续 handler

		// 获取 handler 返回的数据（通过 c.Set 传递）
		data, exists := c.Get("response")
		if !exists {
			return
		}

		s, _ := data.(response.IResponse)

		// 统一格式化
		c.JSON(http.StatusOK, s)
	}
}
