package controll

import "github.com/gin-gonic/gin"

const (
	kGET  = "GET"
	kPOST = "POST"
)

const (
	kRegister = "/register"
	kUpload   = "/upload"
)

type KResponseHandle = func(c *gin.Context) *KResponse

var GetControllMap map[string]KResponseHandle = map[string]KResponseHandle{
	kRegister: Register,
}

var PostControllMap map[string]KResponseHandle = map[string]KResponseHandle{
	kRegister: Register,
	kUpload:   Upload,
}
