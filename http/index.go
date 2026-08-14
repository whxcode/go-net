package httpServer

import (
	"log"

	"go-net/config"
	"go-net/http/controll"
	"go-net/logs"
	"go-net/middleware"
	"go-net/oss"
	"go-net/wss"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
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

	r.Use(middleware.ResponseMiddleware()) // 使用自定义响应中间件

	// r.Static("/asset", "/home/whx/study/go-net/im/pages/asset/")

	// 文件上传和下载路由
	fileRouter := r.Group("/file")

	fileRouter.POST(controll.KUpload, execute(controll.FileController.Upload))
	fileRouter.POST(controll.KGetfile, execute(controll.FileController.GetFile))

	fileDowloadRouter := r.Group(controll.KDowloadfile)
	fileDowloadRouter.Use(controll.FileController.DownloadMiddleware()) // 使用自定义响应中间件
	fileDowloadRouter.GET("", controll.FileController.DowloadFile)

	{

		userRouterPrivate := r.Group("/user")
		userRouterPrivate.Use(middleware.AuthorizationMiddleware())
		userRouterPrivate.GET(controll.K_User_GetUser, execute(controll.UserControll.GetUser))
		userRouterPrivate.GET(controll.K_User_Logout, execute(controll.UserControll.Logout))

		userRouterPublic := r.Group("/user")
		userRouterPublic.POST(controll.K_User_Register, execute(controll.UserControll.Register))
		userRouterPublic.POST(controll.K_User_Login, execute(controll.UserControll.Login))
	}

	r.GET("/im", wss.HandleWSS) // WebSocket路由

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(config.ConfigData.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
