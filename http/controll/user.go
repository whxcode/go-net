package controll

import "github.com/gin-gonic/gin"

type userControll struct{}

var UserControll = &userControll{}

func (*userControll) GetUser(c *gin.Context) *KResponse {
	return MakeResponse([]string{"get user success"})
}
