package model

import (
	"time"
)

/***
* 好友管理DB
* 数据表设计的是存 2 份数据
*   user_id,friend_id
* - 获取好友列表
* 	select * from friends where user_id = userID and status = FriendStatusAccepted
* - 获取好用申请列表
* 	select * from friends where user_id = userID and status != FriendStatusAccepted and status !=  FriendStatusDeleted
*
* - 发起添加好友申请
* - 同意好友申请
* - 拒绝好友申请
* - 删除好友
*
 */

type RequestFriendStatus = uint8

// `status` tinyint(1) DEFAULT '0' COMMENT '0-待确认, 1-已确认, 2-已拒绝, 3-已删除',
const (
	// 待确认
	FriendStatusPending RequestFriendStatus = iota // 待确认
	// 已确认
	FriendStatusAccepted // 已确认
	// 已拒绝
	FriendStatusRejected // 已拒绝
	// 已删除
	FriendStatusDeleted // 已删除
)

type Friend struct {
	// 唯一 ID
	ID uint `gorm:"primarykey" json:"id"`
	// 自身的 ID 和当前登录的 用户 ID 一直。 (类型 uint)
	UserID UserID `gorm:"column:user_id;index;uniqueIndex:uk_user_friend" json:"userId" swaggertype:"integer"`
	// 好友的 用户ID (类型 uint)
	FriendID UserID `gorm:"column:friend_id;index;uniqueIndex:uk_user_friend" json:"friendId" swaggertype:"integer"`
	// 状态:  `status` tinyint(1) DEFAULT '0' COMMENT '0-待确认, 1-已确认, 2-已拒绝, 3-已删除',
	Status RequestFriendStatus `gorm:"column:status;default:0" json:"status"`
	// 发起好友申请时备注时的备注
	Remark    string    `gorm:"column:remark;size:50" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

type FriendResponse struct {
	*Friend
	// 用户账号
	Username string `gorm:"column:username" json:"username"`
	// 用户中文名称
	Nickname string `gorm:"column:nickname" json:"nickname"`
	// 用户头像
	Avatar string `gorm:"column:avatar" json:"avatar"`
	// 在线状态
	IsOnline bool `gorm:"column:is_online" json:"isOnline"`
}

func (*Friend) TableName() string {
	return "friends"
}
