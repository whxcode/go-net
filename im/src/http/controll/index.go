package controll

import "github.com/gin-gonic/gin"

const (
	kGET  = "GET"
	kPOST = "POST"
)

type RouterPath string

const (
	KRegister    RouterPath = "/register"
	KUpload                 = "/upload"
	KGetfile                = "/getfile"
	KDowloadfile            = "/download/:hash"
)

type KResponseHandle = func(c *gin.Context) *KResponse
