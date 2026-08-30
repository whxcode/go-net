package controller

import (
	"net/http"

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
// @Router /groups/{id} [get]
func (*groupController) GroupID(c *gin.Context) *utils.KResponse {
	id := utils.StringToUInt(c.Param("id"))
	result := model.GroupDB.GroupID(id)

	return utils.MakeResponse(result)
}

// @Summary 修改群信息
// @Tags 群
// @Param id path int true "群ID"
// @Param group body model.GroupChat true "群信息-传递 ID 字段无效"
// @Success 200 {object} utils.KResponse{data=model.GroupChat} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /groups/{id} [put]
func (*groupController) PutGroupID(c *gin.Context) *utils.KResponse {
	id := utils.StringToUInt(c.Param("id"))
	group := utils.ShouldBindBodyWithJSON[*model.GroupChat](c)

	group.ID = id

	result := model.GroupDB.PutGroup(group)

	return utils.MakeResponse(result)
}

type ReseutPostGroupMember struct {
	MemberIDs []model.UserID `json:"memberIDs" swaggertype:"array,integer"`
}

// @Summary 给该群添加群成员
// @Tags 群
// @Param id path int true "群ID"
// @Param request body  ReseutPostGroupMember true "成员ID列表(用户 ID)"
// @Success 200 {object} utils.KResponse{data=model.GroupChat} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /groups/:id/members [post]
func (*groupController) PostGroupMember(c *gin.Context) *utils.KResponse {
	id := utils.StringToUInt(c.Param("id"))
	data := utils.ShouldBindBodyWithJSON[*ReseutPostGroupMember](c)

	model.GroupDB.PostGroupMembers(id, data.MemberIDs)

	return utils.MakeResponse(http.StatusOK)
}

type RequstPutGroupMember struct {
	MemberIDs []model.UserID `json:"memberIDs" swaggertype:"array,integer"`
}

// @Summary 修改当前用户在该群的信息(如：是否静音、是否加密、是否退群等)
// @Description 注意： 所属群主无法修改 Role、Status 状态
// @Tags 群
// @Param id path int true "群ID"
// @Param request body  model.GroupMember true "需要修改的数据；注意：UserID 和 ID 字段无效，系统会自动获取当前用户的 ID"
// @Success 200 {object} utils.KResponse{data=model.GroupMember} "成功-修改成功后的数据"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /groups/:id/putMember [put]
func (*groupController) PutGroupMember(c *gin.Context) *utils.KResponse {
	groupID := utils.StringToUInt(c.Param("id"))
	userID := utils.GetUserID(c)
	m := utils.ShouldBindBodyWithJSON[*model.GroupMember](c)
	m.UserID = uint(userID)

	return utils.MakeResponse(model.GroupDB.PutGroupMember(groupID, m))
}

var GroupController = &groupController{}
