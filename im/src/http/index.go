package httpServer

import (
	"log"

	"go-net/im/src/http/controll"

	"github.com/gin-gonic/gin"
)

// 自定义响应中间件
func responseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 执行后续 handler

		// 获取 handler 返回的数据（通过 c.Set 传递）
		data, exists := c.Get("response")
		if !exists {
			return
		}

		// 统一格式化
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": data,
		})
	}
}

func Start() {
	r := gin.Default()

	r.Use(responseMiddleware()) // 使用自定义响应中间件

	for k, v := range controll.GetControllMap {
		r.GET(k, func(c *gin.Context) {
			// 调用 handler 并获取返回值
			result := v(c)
			// 将返回值设置到上下文中
			c.Set("response", result)
		})
	}

	for k, v := range controll.PostControllMap {
		r.POST(k, func(c *gin.Context) {
			// 调用 handler 并获取返回值
			result := v(c)
			// 将返回值设置到上下文中
			c.Set("response", result)
		})
	}

	r.Static("/src", "./src")

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
