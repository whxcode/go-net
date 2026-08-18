package httpServer

import (
	"log"

	"go-net/config"
	"go-net/http/controller"
	"go-net/logs"
	"go-net/middleware"
	"go-net/oss"
	"go-net/wss"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
}

func execute(handle controller.KResponseHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := handle(c)
		c.Set("response", result)
	}
}

func Start() {
	config.ConfigData.Dump()
	oss.Init()

	r := gin.New()
	r.LoadHTMLGlob(config.ConfigData.Server.TemplatePath)

	api := r.Group("/api")
	api.Use(logs.LoggerMiddleware()) // 自定义日志中间件
	api.Use(gin.Recovery())
	// 使用默认CORS中间件，允许所有跨域请求
	api.Use(cors.Default())

	api.Use(middleware.ResponseMiddleware()) // 使用自定义响应中间件

	// r.Static("/asset", "/home/whx/study/go-net/im/pages/asset/")

	// 文件上传和下载路由
	fileRouter := api.Group("/file")

	fileRouter.POST(controller.KUpload, execute(controller.FileController.Upload))
	fileRouter.POST(controller.KGetfile, execute(controller.FileController.GetFile))

	fileDowloadRouter := api.Group(controller.KDowloadfile)
	fileDowloadRouter.Use(controller.FileController.DownloadMiddleware()) // 使用自定义响应中间件
	fileDowloadRouter.GET("", controller.FileController.DowloadFile)

	{

		userRouterPrivate := api.Group("/user")
		userRouterPrivate.Use(middleware.AuthorizationMiddleware())
		userRouterPrivate.GET(controller.K_User_GetUser, execute(controller.UserControll.GetUser))
		userRouterPrivate.GET(controller.K_User_Logout, execute(controller.UserControll.Logout))
		userRouterPrivate.GET(controller.K_User_GetUsers, execute(controller.UserControll.GetUsers))

		userRouterPublic := api.Group("/user")
		userRouterPublic.POST(controller.K_User_Register, execute(controller.UserControll.Register))
		userRouterPublic.POST(controller.K_User_Login, execute(controller.UserControll.Login))
	}

	{
		wssRouter := api.Group("/ws")
		wssRouter.Use(middleware.BeatMiddleware())
		wssRouter.GET("/im", middleware.BeatExecute(wss.RequestWsHandle, wss.ChannelHandle)) // WebSocket路由
	}

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(config.ConfigData.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
