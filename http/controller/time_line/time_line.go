package timeLineController

import (
	"go-net/model"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

/**
*
## 六、时间线 Timeline

| 方法 | 路径 | 返回 (data) |
|---|---|---|
| GET | /api/timeline | 数组 `[{id, text_content, is_anonymous, images, videos, created_at}]` |
| POST | /api/timeline | 任意（body: {text_content, is_anonymous, images}） |
| GET | /api/timeline/:postId | 单条详情 |
| DELETE | /api/timeline/:postId | 任意 |
| POST | /api/timeline/:postId/like | 任意 |
| DELETE | /api/timeline/:postId/like | 任意 |
| POST | /api/timeline/:postId/comments | 任意 |
*
* */

// @Summary 获取时间线列表 (所有可见的)
// @Tags 时间线
// @Param limit query int false "限制条数" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} utils.KResponse{data=[]model.MomentResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /timeline [get]
func GetTimeLines(c *gin.Context) *utils.KResponse {
	return utils.MakeResponse("--")
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

	return utils.MakeResponse(data)
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
	data.OwnerID = uint(utils.GetUserID(c))

	return utils.MakeResponse(data)
}

// @Summary 删除一条时间线
// @Tags 时间线
// @Param id path int true "时间线ID"
// @Success 200 {object} utils.KResponse{data=[]model.MomentResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /timeline/:id [delete]
func DeleteTimeLine(c *gin.Context) *utils.KResponse {
	id := uint(utils.StringToUserID(c.Param("id")))

	return utils.MakeResponse(id)
}
