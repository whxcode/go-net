package controll

import (
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) *KResponse {
	return MakeResponse([]string{"register success"})
}
