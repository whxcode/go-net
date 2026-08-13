package httpServer

import (
	"fmt"
	"log"

	"go-net/config"
	"go-net/http/controll"
	"go-net/http/middleware"
	"go-net/logs"
	"go-net/oss"

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

	// r.Static("/asset", "/home/whx/study/go-net/im/pages/asset/")

	// 文件上传和下载路由
	fileRouter := r.Group("/file")

	fileRouter.POST(controll.KUpload, execute(controll.FileController.Upload))
	fileRouter.POST(controll.KGetfile, execute(controll.FileController.GetFile))

	fileDowloadRouter := r.Group(controll.KDowloadfile)
	fileDowloadRouter.Use(controll.FileController.DownloadMiddleware()) // 使用自定义响应中间件
	fileDowloadRouter.GET("", controll.FileController.DowloadFile)

	userRouter := r.Group("/user")
	userRouter.Use(middleware.AuthorizationMiddleware())
	userRouter.GET(controll.K_User_GetUser, execute(controll.UserControll.GetUser))

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(config.ConfigData.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
