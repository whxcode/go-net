package controller

import (
	"strconv"

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
*
*
* 好友接口
| GET | /api/friends/friends | 数组 `[{id, userId, friendId, status, remark}]`（前端取 friendId 当好友 id） | ✅ |
| GET | /api/friends | 数组 `[{id, username, nickname, avatar, is_online, auto_delete}]`（前端 store 的 Friend 结构，多个页面在用） | ⬜ |
| GET | /api/friends/requests | 申请数组（前端原样 setRequests） | ⬜ |
| POST | /api/friends/request | `{already_friends?: boolean}`（body: {friend_id, message}） | ⬜ |
| POST | /api/friends/accept | 任意（body: {friend_id}） | ⬜ |
| POST | /api/friends/auto-delete | 任意（body: {friend_id, auto_delete}） | ⬜ |
| PUT | /api/friends/remark | 任意（body: {friend_id, remark}） | ⬜ |
| POST | /api/tags/:tagId/assign | 任意（body: {friend_ids}） | ⬜ |
| POST | /api/tags/:tagId/unassign | 任意 | ⬜ |
| DELETE | /api/tags/:id | 任意 | ⬜ |
*
*/

type firendsControll struct{}

var FriendController *firendsControll = &firendsControll{}

// @Summary 获取好友列表
// @Tags 好友模块
// @Success 200 {object} utils.KResponse{data=[]model.FriendResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /friends/friends [get]
func (*firendsControll) Firends(c *gin.Context) *utils.KResponse {
	u := utils.GetUserID(c)

	return utils.MakeResponse(model.FriendDB.Firends(u))
}

type RequestResponse struct {
	// 好友的 用户ID (类型 uint)
	FriendID model.UserID `json:"friendID" swaggertype:"integer"`
	// 备注信息
	Remark string `json:"remark"`
}

// @Summary 发起好友申请
// @Description 好友表没有关系则新增一条记录 status = 0
// @Description 好友表中已经存在；且 status != 1 或 != 0时；则重新更新 status = 0 和 remark(如果有传递)
// @Description 好友表中已经存在；且 status == 1 或 = 0时；不会有任何变化；返回成功
// @Tags 好友模块
// @Param request body RequestResponse true "发起好友申请"
// @Success 200 {object} utils.KResponse{data=RequestResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /friends/request [post]
func (*firendsControll) Request(c *gin.Context) *utils.KResponse {
	UserID := utils.GetUserID(c)
	param := utils.ShouldBindBodyWithJSON[*RequestResponse](c)

	model.FriendDB.Request(UserID, param.FriendID, param.Remark)

	return utils.MakeResponse(param)
}

// @Summary 获取申请列表
// @Description user_id = 当前登录的代表主动发起申请(等待别人反馈)；否则；是其他用户(friend_id)向当前用户发起生气 [同意，拒绝]
// @Tags 好友模块
// @Success 200 {object} utils.KResponse{data=[]model.FriendResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /friends/requests [get]
func (*firendsControll) Requests(c *gin.Context) *utils.KResponse {
	u := utils.GetUserID(c)

	return utils.MakeResponse(model.FriendDB.Requests(u))
}

type PutRequesetRequestResponse struct {
	// 状态 `status` tinyint(1) DEFAULT '0' COMMENT '0-待确认, 1-已确认, 2-已拒绝, 3-已删除',
	Status model.RequestFriendStatus `json:"status" swaggertype:"integer"`
}

// @Summary 修改好友申请状态,可用于 同意(1),拒绝(2)、删除(3)
// @Description 如果将好友移除；同时也会移除对方的好友关系
// @Tags 好友模块
// @Param id path int true "修改好友申请状态;好友列表（表的 ID）"
// @Param request body PutRequesetRequestResponse true "修改好友申请状态"
// @Success 200 {object} utils.KResponse{data=PutRequesetRequestResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /friends/{id} [put]
func (*firendsControll) PutRequesetRequest(c *gin.Context) *utils.KResponse {
	friendTableID := c.Param("id")
	req := utils.ShouldBindBodyWithJSON[*PutRequesetRequestResponse](c)
	id, _ := strconv.ParseUint(friendTableID, 10, 64)

	model.FriendDB.PutRequestFriendStatus(uint(id), req.Status)

	return utils.MakeResponse(req)
}
