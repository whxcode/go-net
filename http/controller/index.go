package controller

import (
	"go-net/response"

	"github.com/gin-gonic/gin"
)

// 路由路径常量：路径都是各分组下的相对路径，分组前缀写在 http/index.go 里
const (
	// 文件模块：挂在 /api/file 分组下（见 http/index.go）
	// 预览/下载都按文件 hash 定位，hash 是文件内容的 sha256，也是落盘文件名。
	KFileUpload      = "/upload"             // 上传文件（multipart，字段名 files）
	KFileSignURLs    = "/signurls"           // 批量把 hash 换成带签名的临时下载地址
	KFilePreview     = "/preview/:hash"      // 按 hash 直出文件内容（无鉴权）
	KFilePreviewMeta = "/preview/:hash/meta" // 按 hash 直出文件元信息 JSON（无鉴权）
	KFileDownload    = "/download/:hash"     // 带签名下载（expred + signature）

	KUserGetUser     = "/get"
	KUserGetUserByID = "/:id"
	KUserLogout      = "/logout"
	KUserRegister    = "/register"
	KUserLogin       = "/login"
	KUserGetUsers    = "/users"
	UserPassword     = "/password"
	UserAvatar       = "/avatar"
	UserNickname     = "/nickname"

	KFriends          = "/friends"
	KFriendRequest    = "/request"
	KFriendRequests   = "/requests"
	KFriendPutRequest = "/:id"

	KMessageFriend = "/friend/:friendID"
	KMessageGroup  = "/group/:groupID"

	KGroupGetGroups  = "/"
	KGroupGetGroup   = "/:id"
	KPostGroupMember = "/:id/members"
	KGroupPutMember  = "/:id/putMember"

	KMoments   = "/"
	KMomentsID = "/:id"

	KMomentsUser   = "/user/:userID"
	KMomentPrivacy = "/privacy/:userID"

	KMomentsLikes = "/likes/:id"

	// 朋友圈评论相关
	KMomentComments = "/:id/comments"

	// 时间线相关
	KTimeline   = "/"
	KTimelineID = "/:id"

	KTimelineComments = "/:id/comments"
	KTimelineLikes    = "/:id/likes"
)

type KResponseHandle = func(c *gin.Context) *response.KResponse
