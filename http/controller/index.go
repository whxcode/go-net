package controller

import (
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

const (
	kGET  = "GET"
	kPOST = "POST"
)

type RouterPath string

const (
	KRegister    RouterPath = "/register"
	KUpload                 = "/upload"
	KGetfile                = "/getfile"
	KGetfileHash            = "/getfile/:hash"
	KPreviewFile            = "/:hash"

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

	KMoments               = "/"
	KMomentsUser           = "/:userID"
	KMomentsLike           = "/:id/like"
	KMomentsLikes          = "/likes/:id"
	KMomentsID             = "/:id"
	KMomentPrivacyTargetID = "/privacy/:targetId"
	KMomentsPrivacy        = "/privacy"

	// 朋友圈评论相关
	KMomentComments = "/comments/:id"
)

type KResponseHandle = func(c *gin.Context) *utils.KResponse
