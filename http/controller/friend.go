package controller

import (
	"net/http"

	"go-net/model"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

/*
接口	方法	路径	说明
发送好友申请	POST	/api/friends/request	添加好友
获取好友申请列表	GET	/api/friends/requests	待处理列表
同意好友申请	PUT	/api/friends/request/:id	同意
拒绝好友申请	PUT	/api/friends/request/:id	拒绝
好友列表	GET	/api/friends	✅ 已有
*/

type firendsControll struct{}

var FirendsControll *firendsControll = &firendsControll{}

func (*firendsControll) Request(c *gin.Context) *utils.KResponse {
	type Param struct {
		UserID   model.UserID `json:"userID"`
		FriendID model.UserID `json:"friendID"`
		Remark   string       `json:"remark"`
	}

	var param Param
	if err := c.ShouldBindBodyWithJSON(&param); err != nil {
		return utils.MakeResponseWidthCode(err, http.StatusInternalServerError)
	}

	return nil
}
