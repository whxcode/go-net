package controll

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

	K_User_GetUser  = "/get"
	K_User_Logout   = "/logout"
	K_User_Register = "/register"
	K_User_Login    = "/login"
)

type KResponseHandle = func(c *gin.Context) *utils.KResponse
