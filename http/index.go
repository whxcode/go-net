package httpServer

import (
	"log"

	"go-net/config"
	"go-net/http/controller"
	timeLineController "go-net/http/controller/time_line"
	"go-net/logs"
	"go-net/middleware"
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

	r := gin.New()

	// Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/redis", controller.GetRedis)

	api := r.Group("/api")
	api.Use(logs.LoggerMiddleware()) // 自定义日志中间件

	api.Use(middleware.RecoverMiddleware())
	api.Use(cors.Default())
	api.Use(middleware.ResponseMiddleware())

	{

		// 文件接口 /api/file/*：路径和常量见 controller/index.go 里的 KFile* 常量
		// 三类接口的返回体不一样，读代码时先看这里：
		//   1. upload / signurls 返回 JSON，走 execute + ResponseMiddleware 统一包装成 KResponse
		//   2. preview / preview meta 直接返回文件本体（图片、视频、meta.json），没有包装
		//   3. download 返回文件本体，但必须先通过临时地址上的签名校验
		fileRouter := api.Group("/file")

		fileRouter.POST(controller.KFileUpload, execute(controller.FileController.Upload))
		fileRouter.POST(controller.KFileSignURLs, execute(controller.FileController.SignURLs))

		// 预览：按 hash 直出，无鉴权；前端 <img src>、<video src> 直接引用这两个地址
		fileRouter.GET(controller.KFilePreview, controller.FileController.PreviewFile)
		fileRouter.GET(controller.KFilePreviewMeta, controller.FileController.PreviewMetaFile)

		// 下载：需要临时地址上的 expred + signature，所以只给这一条路由挂签名校验中间件
		fileRouter.GET(controller.KFileDownload,
			controller.FileController.DownloadMiddleware(),
			controller.FileController.DownloadFile)
	}

	{

		userRouterPrivate := api.Group("/users")
		userRouterPrivate.Use(middleware.AuthorizationMiddleware())
		userRouterPrivate.GET(controller.KUserGetUser, execute(controller.UserControll.GetUser))
		userRouterPrivate.GET(controller.KUserGetUserByID, execute(controller.UserControll.GetUserByID))
		userRouterPrivate.GET(controller.KUserLogout, execute(controller.UserControll.Logout))
		userRouterPrivate.GET(controller.KUserGetUsers, execute(controller.UserControll.GetUsers))
		userRouterPrivate.PUT(controller.UserPassword, execute(controller.UserControll.UserPutPassword))
		userRouterPrivate.PUT(controller.UserAvatar, execute(controller.UserControll.UserPutAvatar))
		userRouterPrivate.PUT(controller.UserNickname, execute(controller.UserControll.UserPutNickname))

		userRouterPublic := api.Group("/users")
		userRouterPublic.POST(controller.KUserRegister, execute(controller.UserControll.Register))
		userRouterPublic.POST(controller.KUserLogin, execute(controller.UserControll.Login))
	}

	{

		friendRouterPrivate := api.Group("/friends")
		friendRouterPrivate.Use(middleware.AuthorizationMiddleware())
		friendRouterPrivate.GET(controller.KFriends, execute(controller.FriendController.Firends))
		friendRouterPrivate.GET(controller.KFriendRequests, execute(controller.FriendController.Requests))
		friendRouterPrivate.POST(controller.KFriendRequest, execute(controller.FriendController.Request))
		friendRouterPrivate.PUT(controller.KFriendPutRequest, execute(controller.FriendController.PutRequesetRequest))

	}

	{
		messageRouter := api.Group("/messages")
		messageRouter.Use(middleware.AuthorizationMiddleware())
		messageRouter.GET(controller.KMessageFriend, execute(controller.MessageController.GetFrinedMessages))
		messageRouter.GET(controller.KMessageGroup, execute(controller.MessageController.GetGroupMessages))

	}

	{
		messageRouter := api.Group("/groups")
		messageRouter.Use(middleware.AuthorizationMiddleware())
		messageRouter.GET(controller.KGroupGetGroups, execute(controller.GroupController.Groups))
		messageRouter.GET(controller.KGroupGetGroup, execute(controller.GroupController.GroupID))
		messageRouter.PUT(controller.KGroupGetGroup, execute(controller.GroupController.PutGroupID))
		messageRouter.POST(controller.KPostGroupMember, execute(controller.GroupController.PostGroupMember))
		messageRouter.PUT(controller.KGroupPutMember, execute(controller.GroupController.PutGroupMember))
	}

	{
		momentsRouter := api.Group("/moments")
		momentsRouter.Use(middleware.AuthorizationMiddleware())

		momentsRouter.GET(controller.KMoments, execute(controller.MomentController.Moments))
		momentsRouter.POST(controller.KMoments, execute(controller.MomentController.PostMoments))
		momentsRouter.DELETE(controller.KMomentsID, execute(controller.MomentController.MomentDelete))

		// 某个好友的朋友圈记录
		momentsRouter.GET(controller.KMomentsUser, execute(controller.MomentController.MomentUserID))

		// 隐私
		momentsRouter.GET(controller.KMomentPrivacy, execute(controller.MomentController.MomentPrivacyTargetID))
		momentsRouter.POST(controller.KMomentPrivacy, execute(controller.MomentController.PostMomentPrivacy))

		// 点赞
		momentsRouter.GET(controller.KMomentsLikes, execute(controller.MomentController.MomentLikes))
		momentsRouter.POST(controller.KMomentsLikes, execute(controller.MomentController.MomentsIdLike))
		momentsRouter.DELETE(controller.KMomentsLikes, execute(controller.MomentController.MomentsIdUnLike))

		// 评论
		momentsRouter.GET(controller.KMomentComments, execute(controller.MomentController.MomentComments))
		momentsRouter.POST(controller.KMomentComments, execute(controller.MomentController.PostMomentComment))
		momentsRouter.PUT(controller.KMomentComments, execute(controller.MomentController.PutMomentComment))
		momentsRouter.DELETE(controller.KMomentComments, execute(controller.MomentController.DeleteMomentComment))

	}

	{
		momentsRouter := api.Group("/timeline")
		momentsRouter.Use(middleware.AuthorizationMiddleware())

		momentsRouter.GET(controller.KTimeline, execute(timeLineController.GetTimeLines))
		momentsRouter.POST(controller.KTimeline, execute(timeLineController.PostTimeLine))
		momentsRouter.PUT(controller.KTimelineID, execute(timeLineController.PutTimeLine))
		momentsRouter.DELETE(controller.KTimelineID, execute(timeLineController.DeleteTimeLine))

		momentsRouter.GET(controller.KMomentComments, execute(timeLineController.Comments))
		momentsRouter.POST(controller.KMomentComments, execute(timeLineController.PostComment))
		momentsRouter.PUT(controller.KMomentComments, execute(timeLineController.PutComment))
		momentsRouter.DELETE(controller.KMomentComments, execute(timeLineController.DeleteComment))

		momentsRouter.GET(controller.KTimelineLikes, execute(timeLineController.Likes))
		momentsRouter.POST(controller.KTimelineLikes, execute(timeLineController.PostLike))
		momentsRouter.DELETE(controller.KTimelineLikes, execute(timeLineController.DeleteLike))

	}

	{
		wssRouter := api.Group("/ws")
		wssRouter.GET("/im", wss.IM)
	}

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(config.ConfigData.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
