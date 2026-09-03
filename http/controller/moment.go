/*
1. GET /api/moments —— 朋友圈列表（首页刷的流，自己+好友的）。前端无参数。
2. GET /api/moments/user/:userId?limit=3 —— 看某个人的朋友圈。前端传 userId（URL里
   ）+ limit（条数，个人主页用3）。
3. POST /api/moments —— 发朋友圈。前端传：text_content（文字）、images（图片数组）
   、video（视频对象，含 url/thumbnail/duration，可选）。
4. POST /api/moments/:id/like —— 点赞。传 id。
5. DELETE /api/moments/:id/like —— 取消点赞。传 id。
6. DELETE /api/moments/:id —— 删自己的朋友圈。传 id。
7. POST /api/moments/:id/comments —— 评论。传 id + content（评论内容）。
8. GET /api/moments/privacy/:targetId —— 查和某个人的屏蔽设置。传 targetId。

9. POST /api/moments/privacy —— 设置屏蔽。传 target_id、hide_their（我不看TA的）、h
   ide_mine（不让TA看我的）。
*/

package controller

import (
	"net/http"

	"go-net/model"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

type momentController struct{}

// @Summary 获取朋友圈列表 (首页刷的流,当前登陆用户 + 好友)
// @Tags 朋友圈
// @Param limit query int false "限制条数" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} utils.KResponse{data=model.UserResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /moments [get]
func (m *momentController) Moments(c *gin.Context) *utils.KResponse {
	limit, offset := utils.ParsePageQuery(c)

	return utils.MakeResponse(model.MomentDB.GetMoments(utils.GetUserID(c), limit, offset))
}

// @Summary 看某个人的朋友圈
// @Tags 朋友圈
// @Param limit query int false "限制条数" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} utils.KResponse{data=[]model.Moment} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /moments/:userID [get]
func (m *momentController) MomentUserID(c *gin.Context) *utils.KResponse {
	limit, offset := utils.ParsePageQuery(c)

	return utils.MakeResponse([]int{limit, offset})
}

// @Summary 发布朋友圈;仅取 Elements,Visbile 字段.
// @Tags 朋友圈
// @Param request body model.Moment true "注册请求"
// @Success 200 {object} utils.KResponse{data=model.Moment} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /moments [post]
func (m *momentController) PostMoments(c *gin.Context) *utils.KResponse {
	moment := utils.ShouldBindBodyWithJSON[*model.Moment](c)
	moment.OwnerID = uint(utils.GetUserID(c))

	return utils.MakeResponse(model.MomentDB.AddMoment(moment))
}

// @Summary  POST /api/moments/:id/like —— 点赞。传 id。
// @Tags 朋友圈
// @Param id path int true "朋友圈ID"
// @Success 200 {object} utils.KResponse{data=Moment} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /moments/:id/like [post]
func (m *momentController) MomentsIdLike(c *gin.Context) *utils.KResponse {
	id := c.Param("id")

	return utils.MakeResponse(model.MomentDB.MomentLike(utils.GetUserID(c), utils.StringToUInt(id)))
}

// @Summary  取消点赞
// @Tags 朋友圈
// @Param id path int true "朋友圈ID"
// @Success 200 {object} utils.KResponse{data=int} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /moments/:id/like [delete]
func (m *momentController) MomentsIdUnLike(c *gin.Context) *utils.KResponse {
	id := c.Param("id")
	return utils.MakeResponse(model.MomentDB.MomentUnLike(utils.GetUserID(c), utils.StringToUInt(id)))
}

// @Summary  删除朋友圈
// @Tags 朋友圈
// @Param id path int true "朋友圈ID"
// @Success 200 {object} utils.KResponse{data=int} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /moments/:id [delete]
func (m *momentController) MomentDelete(c *gin.Context) *utils.KResponse {
	id := c.Param("id")

	model.MomentDB.DeleteMoment(utils.StringToUInt(id))

	return utils.MakeResponse(http.StatusOK)
}

// @Summary  查和某个人的屏蔽设置。传 targetId。
// @Tags 朋友圈
// @Param targetId path int true "用户ID"
// @Success 200 {object} utils.KResponse{data=int} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /moments/privacy/:targetId [get]
func (m *momentController) MomentPrivacyTargetID(c *gin.Context) *utils.KResponse {
	targetId := c.Param("targetId")

	return utils.MakeResponse(targetId)
}

type RequestPrivacy struct {
	TargetId  uint `json:"target_id" binding:"required"`
	HideTheir bool `json:"hide_their" binding:"required"` // 我不看TA的
	HideMine  bool `json:"hide_mine" binding:"required"`  // 不让TA看我的
}

// @Summary  设置屏蔽。传 target_id、hide_their（我不看TA的）、hide_mine（不让TA看我的）。
// @Tags 朋友圈
// Param request body RequestPrivacy  true "请求体"
// @Success 200 {object} utils.KResponse{data=int} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /moments/privacy [post]
func (m *momentController) MomentPrivacy(c *gin.Context) *utils.KResponse {
	data := utils.ShouldBindBodyWithJSON[*RequestPrivacy](c)

	return utils.MakeResponse(data)
}

var MomentController = &momentController{}
