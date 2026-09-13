package timeLineController

import (
	"net/http"

	dbTimeLine "go-net/db/time_line"
	"go-net/model"
	"go-net/response"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

// @Summary 获取该时间线的评论列表
// @Tags 时间线/评论
// @Param id path int true "时间线ID"
// @Success 200 {object} response.KResponse{data=[]model.TimeComment} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /timeline/:id/comments [get]
func Comments(c *gin.Context) *response.KResponse {
	id := utils.StringToUInt(c.Param("id"))

	return response.MakeResponse(dbTimeLine.Comments(id))
}

// @Summary 给某个时间线添加一条评论 (只取 Elements 字段)
// @Tags 时间线/评论
// @Param id path int true "时间线ID"
// @Param request body model.TimeComment true "评论体"
// @Success 200 {object} response.KResponse{data=model.TimeComment} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /timeline/:id/comments [post]
func PostComment(c *gin.Context) *response.KResponse {
	data := utils.ShouldBindBodyWithJSON[*model.TimeComment](c)
	data.OwnerID = uint(utils.GetUserID(c))
	data.TimeLineID = utils.StringToUInt(c.Param("id"))
	t := dbTimeLine.PostComments(data)

	return response.MakeResponseWithUserInfo(t)
}

// @Summary 修改一条评论 (只取 Elements 字段)
// @Tags 时间线/评论
// @Param id path int true "评论ID"
// @Param request body model.TimeComment true "评论体"
// @Success 200 {object} response.KResponse{data=model.TimeComment} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /timeline/:id/comments [put]
func PutComment(c *gin.Context) *response.KResponse {
	data := utils.ShouldBindBodyWithJSON[*model.TimeComment](c)
	id := utils.StringToUInt(c.Param("id"))
	data.ID = id

	return response.MakeResponse(dbTimeLine.PutComment(data))
}

// @Summary 删除一条评论 (只取 Elements 字段)
// @Tags 时间线/评论
// @Param id path int true "评论ID"
// @Success 200 {object} response.KResponse{data=int} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /timeline/:id/comments [delete]
func DeleteComment(c *gin.Context) *response.KResponse {
	id := utils.StringToUInt(c.Param("id"))

	dbTimeLine.DeleteComment(id)

	return response.MakeResponse(http.StatusOK)
}
