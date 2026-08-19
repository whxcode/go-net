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
	KDowloadfile            = "/file/download/:hash"

	KUserGetUser  = "/get"
	KUserLogout   = "/logout"
	KUserRegister = "/register"
	KUserLogin    = "/login"
	KUserGetUsers = "/users"

	KFriends          = "/friends"
	KFriendRequest    = "/request"
	KFriendRequests   = "/requests"
	KFriendPutRequest = "/request/:id"
)

type KResponseHandle = func(c *gin.Context) *utils.KResponse
