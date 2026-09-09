package timeLineController

import (
	dbTimeLine "go-net/db/time_line"
	"go-net/model"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

// @Summary 获取时间线列表 (所有可见的)
// @Tags 时间线
// @Param limit query int false "限制条数" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} utils.KResponse{data=[]model.MomentResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /timeline [get]
func GetTimeLines(c *gin.Context) *utils.KResponse {
	return utils.MakeResponse(dbTimeLine.GetTimeLines(utils.ParsePageQuery(c)))
}

// @Summary 添加一个时间线 (仅取 Elements 字段)
// @Tags 时间线
// @Param request body model.TimeLine true "内容主体"
// @Success 200 {object} utils.KResponse{data=[]model.MomentResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /timeline [post]
func PostTimeLine(c *gin.Context) *utils.KResponse {
	data := utils.ShouldBindBodyWithJSON[*model.TimeLine](c)
	data.OwnerID = uint(utils.GetUserID(c))

	return utils.MakeResponse(dbTimeLine.PostTimeLines(data))
}

// @Summary 修改一个时间线 (仅取 Elements 字段)
// @Tags 时间线
// @Param id path int true "时间线ID"
// @Param request body model.TimeLine true "内容主体"
// @Success 200 {object} utils.KResponse{data=[]model.MomentResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /timeline/:id [put]
func PutTimeLine(c *gin.Context) *utils.KResponse {
	data := utils.ShouldBindBodyWithJSON[*model.TimeLine](c)
	data.ID = uint(utils.StringToUserID(c.Param("id")))

	return utils.MakeResponse(dbTimeLine.PutTimeLines(data))
}

// @Summary 设置时间线的状态 只取 status 字段 (status: 0-正常，1-删除,2-被举报中,3-举报成功,4-举报失败)
// @Tags 时间线
// @Param id path int true "时间线ID"
// @Param request body model.TimeLine true "内容主体"
// @Success 200 {object} utils.KResponse{data=[]model.MomentResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /timeline/:id [delete]
func DeleteTimeLine(c *gin.Context) *utils.KResponse {
	data := utils.ShouldBindBodyWithJSON[*model.TimeLine](c)
	data.ID = uint(utils.StringToUInt(c.Param("id")))

	return utils.MakeResponse(dbTimeLine.DeleteTimeLines(data.ID, data.Status))
}
