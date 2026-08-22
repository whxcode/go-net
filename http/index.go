package httpServer

import (
	"log"

	"go-net/config"
	"go-net/http/controller"
	"go-net/logs"
	"go-net/middleware"
	"go-net/oss"
	"go-net/redis"
	"go-net/wss"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "go-net/docs" // ✅ 打开这行，你的项目名替换成 go.mod 里的 module 名
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

	// Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 获取所有redis key:value 及json格式化的形式返回给前端调式
	// 获取所有 Redis key:value，JSON 格式返回
	r.GET("/redis", func(c *gin.Context) {
		// 1. 获取所有 key
		keys, err := redis.RedisClient.Keys(c, "*").Result()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 2. 遍历获取每个 key 的值
		result := make(map[string]interface{})
		for _, key := range keys {
			// 获取值（字符串类型）
			val, err := redis.RedisClient.Get(c, key).Result()
			if err != nil {
				// 如果 key 不存在或类型不对，跳过
				continue
			}
			result[key] = val
		}

		// 3. 返回 JSON
		c.JSON(200, gin.H{
			"total": len(result),
			"data":  result,
		})
	})

	api := r.Group("/api")
	api.Use(logs.LoggerMiddleware()) // 自定义日志中间件
	// ✅ 自定义 Recovery 中间件
	api.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录日志
				// 返回统一格式
				c.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": err,
				})
			}
		}()
		c.Next()
	})
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

		userRouterPrivate := api.Group("/users")
		userRouterPrivate.Use(middleware.AuthorizationMiddleware())
		userRouterPrivate.GET(controller.KUserGetUser, execute(controller.UserControll.GetUser))
		userRouterPrivate.GET(controller.KUserGetUserByID, execute(controller.UserControll.GetUserByID))
		userRouterPrivate.GET(controller.KUserLogout, execute(controller.UserControll.Logout))
		userRouterPrivate.GET(controller.KUserGetUsers, execute(controller.UserControll.GetUsers))
		userRouterPrivate.PUT(controller.UserPassword, execute(controller.UserControll.UserPutPassword))

		userRouterPublic := api.Group("/users")
		userRouterPublic.POST(controller.KUserRegister, execute(controller.UserControll.Register))
		userRouterPublic.POST(controller.KUserLogin, execute(controller.UserControll.Login))
	}

	{

		friendRouterPrivate := api.Group("/friend")
		friendRouterPrivate.Use(middleware.AuthorizationMiddleware())
		friendRouterPrivate.GET(controller.KFriends, execute(controller.FriendController.Firends))
		friendRouterPrivate.POST(controller.KFriendRequest, execute(controller.FriendController.Request))
		friendRouterPrivate.GET(controller.KFriendRequests, execute(controller.FriendController.Requests))
		friendRouterPrivate.PUT(controller.KFriendPutRequest, execute(controller.FriendController.PutRequesetRequest))

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
