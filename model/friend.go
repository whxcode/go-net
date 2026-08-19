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
	FriendStatusPending  RequestFriendStatus = iota // 待确认
	FriendStatusAccepted                            // 已确认
	FriendStatusRejected                            // 已拒绝
	FriendStatusDeleted                             // 已删除
)

type Friend struct {
	ID        uint                `gorm:"primarykey" json:"id"`
	UserID    UserID              `gorm:"column:user_id;index;uniqueIndex:uk_user_friend" json:"userId"`
	FriendID  UserID              `gorm:"column:friend_id;index;uniqueIndex:uk_user_friend" json:"friendId"`
	Status    RequestFriendStatus `gorm:"column:status;default:0" json:"status"`
	Remark    string              `gorm:"column:remark;size:50" json:"remark"`
	CreatedAt time.Time           `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time           `gorm:"column:updated_at" json:"updatedAt"`
}

func (*Friend) TableName() string {
	return "friends"
}

type friendDB struct{}

var FriendDB = &friendDB{}

// 根据用户id 查询好用列表
func (*friendDB) Firends(userId UserID) []*Friend {
	var result []*Friend

	err := DB.Table("friends"). //
					Select("friends.friend_id as user_id,users.username,friends.remark").
					Joins("left join users ON users.id = friends.friend_id").
					Where("friends.user_id = ? AND friends.status = ?", userId, 1).
					Find(&result).Error
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	return result
}

/**
* 获取好友申请列表
* */

func (*friendDB) Requests(userID UserID) (result []*Friend) {
	// 我主动发起的申请
	// 其他人向我发起的申请
	err := DB.Table("friends").Where("(user_id = ? AND status IN (?)) OR (friend_id = ? AND status IN (?))", userID, []RequestFriendStatus{FriendStatusPending, FriendStatusRejected}, userID, []RequestFriendStatus{FriendStatusPending, FriendStatusRejected}).Find(&result).Error
	if err != nil {
		panic(err)
	}

	return
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

	err := DB.Create(f).Error
	if err != nil {
		panic(err)
	}

	return f
}

/**
* 同意好友申请
* user_id -> 我
* friendId -> 对方
*
* */
func (*friendDB) AcceptRequeset(userId, friendId UserID) *Friend {
	// 查询时;需要交互一下 查询条件条件
	err := DB.Table("friends").Where("user_id = ? AND friend_id = ?", friendId, userId).
		Update("status", FriendStatusAccepted).Error
	if err != nil {
		panic(err)
	}

	f := &Friend{
		UserID:   userId,
		FriendID: friendId,
		Status:   FriendStatusAccepted,
	}

	err = DB.Create(f).Error
	if err != nil {
		panic(err)
	}

	return nil
}

/**
* 拒绝好友申请
*
* */
func (*friendDB) RejectedRequeset(userId, friendId UserID) error {
	return DB.Table("friends").
		Where("user_id = ? AND friend_id = ?", friendId, userId).
		Update("status", FriendStatusRejected).Error
}

/**
* 删除好友
* */
func (*friendDB) DeleteFriend(userId, friendId UserID) error {
	err := DB.Table("friends").
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userId, friendId, friendId, userId).
		Update("status", FriendStatusDeleted).Error
	if err != nil {
		panic(err)
	}

	return nil
}
