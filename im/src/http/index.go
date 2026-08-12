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

func execute(handle controll.KResponseHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := handle(c)
		c.Set("response", result)
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

	fileRouter := r.Group("/file")
	fileRouter.Use(controll.DownloadMiddleware()) // 使用自定义响应中间件

	fileRouter.POST(controll.KUpload, execute(controll.Upload))
	fileRouter.POST(controll.KGetfile, execute(controll.GetFile))
	fileRouter.GET(controll.KDowloadfile, controll.DowloadFile)

	r.Static("/asset", "/home/whx/study/go-net/im/pages/asset/")

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(config.ConfigData.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
