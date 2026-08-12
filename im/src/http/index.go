package httpServer

import (
	"fmt"
	"log"

	config "go-net/im/src"
	"go-net/im/src/http/controll"
	"go-net/im/src/logs"
	"go-net/im/src/oss"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
}

// 自定义响应中间件
func responseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 执行后续 handler

		// 获取 handler 返回的数据（通过 c.Set 传递）
		data, exists := c.Get("response")
		if !exists {
			return
		}

		s, ok := data.(*controll.KResponse)

		if ok {
			fmt.Println(s) // hello
		}

		// 统一格式化
		c.JSON(200, gin.H{
			"code": s.Code,
			"msg":  "success",
			"data": s.Data,
		})
	}
}

func Start() {
	config.ConfigData.Dump()
	oss.Init()

	r := gin.New()
	r.Use(logs.LoggerMiddleware()) // 自定义日志中间件
	r.Use(gin.Recovery())

	r.LoadHTMLGlob(config.ConfigData.Server.TemplatePath)
	// 使用默认CORS中间件，允许所有跨域请求
	r.Use(cors.Default())

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

	r.Static("/asset", "/home/whx/study/go-net/im/pages/asset/")

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(config.ConfigData.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
