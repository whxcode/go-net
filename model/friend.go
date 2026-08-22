package model

import (
	"fmt"
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

type friendDB struct{}

var FriendDB = &friendDB{}

// 根据用户id 查询好用列表
func (*friendDB) firends(userId UserID, status []RequestFriendStatus, extar string) []*FriendResponse {
	var result []*FriendResponse
	var err error

	if extar == "" {
		err = DB. //
				Table("friends f").
				Select("f.*,u.username,u.nickname,u.avatar").
				Joins("left join users u on u.id = f.friend_id").
				Where("(user_id = ? AND status IN (?))", userId, status).
				Find(&result).Error
	} else {
		err = DB.Table("friends f").
			Select("f.*,u.username,u.nickname,u.avatar").
			Joins("left join users u on u.id = f.friend_id").
			Where("(user_id = ? AND status IN (?))"+extar, userId, status, userId, status).
			Find(&result).Error
	}

	if err != nil {
		panic(err)
	}

	return result
}

// 根据用户id 查询好用列表
func (f *friendDB) Firends(userId UserID) []*FriendResponse {
	return f.firends(userId, []RequestFriendStatus{FriendStatusAccepted}, "")
}

func (f *friendDB) Requests(userID UserID) (result []*FriendResponse) {
	return f.firends(userID, []RequestFriendStatus{FriendStatusPending, FriendStatusRejected}, "OR (friend_id = ? AND status IN (?))")
}

/*
* 好友申请
* user_id -> 我
* friendId -> 对方
*
* */
func (*friendDB) Request(userId, friendId UserID, remark string) *Friend {
	f := &Friend{
		UserID:   userId,
		FriendID: friendId,
		Status:   FriendStatusPending,
		Remark:   remark,
	}

	var existingFriend *Friend

	err := DB.Table("friends").Where("user_id = ? AND friend_id = ?", userId, friendId).Find(&existingFriend).Error

	if err == nil {
		switch existingFriend.Status {
		case FriendStatusAccepted, FriendStatusPending:
		case FriendStatusRejected, FriendStatusDeleted:

			existingFriend.Status = FriendStatusPending
			existingFriend.Remark = remark
			err = DB.Save(existingFriend).Error
			if err != nil {
				panic(err)
			}

			// 也要重置对方队列中的状态
			result := DB.Table("friends").
				Where("user_id = ? and friend_id = ?", friendId, userId).
				Update("status", FriendStatusPending)

			if result.Error != nil {
				panic(err)
			}

			if result.RowsAffected == 0 {
				// 查不到，什么都不做
				fmt.Println("没有记录需要更新")
			}

		}

		return existingFriend
	}

	err = DB.Create(f).Error
	if err != nil {
		panic(err)
	}

	return f
}

func (*friendDB) PutRequestFriendStatus(id uint, status RequestFriendStatus) {
	var f *Friend

	err := DB.Table("friends").Where("id = ?", id).Find(&f).Error
	if err != nil {
		panic(err)
	}

	if f == nil {
		panic("好友申请不存在")
	}

	if status == FriendStatusDeleted {
		err = DB.Table("friends").
			Where("user_id = ? AND friend_id = ?", f.FriendID, f.UserID).
			Update("status", status).Error
		if err != nil {
			panic(err)
		}

	}

	// 查询时;需要交互一下 查询条件条件
	err = DB.Table("friends").Where("id = ?", id).
		Update("status", status).Error
	if err != nil {
		panic(err)
	}
}
