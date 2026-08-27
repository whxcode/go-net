package controller

import (
	"go-net/model"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

type groupController struct{}

// @Summary 获取群列表
// @Tags 群
// @Success 200 {object} utils.KResponse{data=[]model.GroupChatResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /groups [get]
func (*groupController) Groups(c *gin.Context) *utils.KResponse {
	u := utils.GetUserID(c)
	result := model.GroupDB.Groups(u)

	return utils.MakeResponse(result)
}

// @Summary 获取群详情
// @Tags 群
// @Param id path int true "群ID"
// @Success 200 {object} utils.KResponse{data=model.GroupChatResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /groups/:id [get]
func (*groupController) GroupID(c *gin.Context) *utils.KResponse {
	id := utils.StringToUInt(c.Param("id"))
	result := model.GroupDB.GroupID(id)

	return utils.MakeResponse(result)
}

var GroupController = &groupController{}
