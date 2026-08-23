package controller

import (
	"go-net/model"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

type messageController struct{}

/*
## 三、消息（核心）

| 方法 | 路径 | 返回 (data) | 状态 |
|---|---|---|---|
| GET | /api/messages/private/:friendId?limit=50&offset=0 | **直接返回后端 Message 数组**（前端转换层自动映射，字段见下） | ⬜ |
| GET | /api/messages/group/:groupId?limit=50&offset=0 | 同上（群聊，暂未启用） | ⬜ |
| POST | /api/upload | `{url, key}`（FormData file 字段；url 用于消息里图片/文件显示） | ⬜ |
*/

func parseHistoryQuery(c *gin.Context) (model.UserID, int, int) {
	friendID := c.Param("friendID")
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	return utils.StringToUserID(friendID), utils.StringToInt(limit), utils.StringToInt(offset)
}

type GetPrivateMessagesResponse struct {
	Data  []*model.Message `json:"data"`
	Total int              `json:"total"`
	Size  int              `json:"size"`
}

// @Summary 获取与好友的历史记录
// @Tags 好友消息
// @Param limit query int false "限制条数" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} utils.KResponse{data=GetPrivateMessagesResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /messages/private/:friendID [get]
func (*messageController) GetPrivateMessages(c *gin.Context) *utils.KResponse {
	userID := utils.GetUserID(c)
	friendID, limit, offset := parseHistoryQuery(c)

	m, t := model.MessageDB.GetHistory(userID, friendID, limit, offset)

	return utils.MakeResponse(&GetPrivateMessagesResponse{
		Data:  m,
		Total: t,
	})
}

var MessageController = &messageController{}
