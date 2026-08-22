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
)

type KResponseHandle = func(c *gin.Context) *utils.KResponse
