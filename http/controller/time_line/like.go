package timeLineController

import (
	dbTimeLine "go-net/db/time_line"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

// @Summary 获取该时间线的点赞列表
// @Tags 时间线/点赞
// @Param id path int true "时间线ID"
// @Success 200 {object} utils.KResponse{data=[]model.TimeLike} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /timeline/:id/likes [get]
func Likes(c *gin.Context) *utils.KResponse {
	timeID := utils.StringToUInt(c.Param("id"))

	return utils.MakeResponse(dbTimeLine.Likes(timeID))
}

// @Summary 给某个时间线点赞
// @Tags 时间线/点赞
// @Param id path int true "时间线ID"
// @Success 200 {object} utils.KResponse{data=model.Moment} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /timeline/:id/likes [post]
func PostLike(c *gin.Context) *utils.KResponse {
	timeID := utils.StringToUInt(c.Param("id"))
	userID := uint(utils.GetUserID(c))

	return utils.MakeResponse(dbTimeLine.PostLikes(timeID, userID))
}

// @Summary 取消点赞
// @Tags 时间线/点赞
// @Parma id path int true "点赞 ID"
// @Success 200 {object} utils.KResponse{data=model.Moment} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /timeline/:id/likes [delete]
func DeleteLike(c *gin.Context) *utils.KResponse {
	timeID := utils.StringToUInt(c.Param("id"))
	userID := uint(utils.GetUserID(c))

	return utils.MakeResponse(dbTimeLine.DeleteLikes(timeID, userID))
}
