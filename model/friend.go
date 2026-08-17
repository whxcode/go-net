package model

import (
	"time"
)

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
func (*friendDB) GetFirends(userId UserID) []*Friend {
	var result []*Friend

	err := DB.Table("friends"). //
					Select("friends.friend_id as user_id,users.username,friends.remark").
					Joins("left join users ON users.id = friends.friend_id").
					Where("friends.user_id = ? AND friends.status = ?", userId, 1).
					Find(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

// 根据用户 id；添加好友
func (*friendDB) Requesets(userId, friendId UserID, remark string) *Friend {
	f := &Friend{
		UserID:   userId,
		FriendID: friendId,
		Status:   FriendStatusPending,
		Remark:   remark,
	}

	DB.Create(f)

	return nil
}

// 同意
func (*friendDB) AcceptRequest(userID, friendID UserID) error {
	DB.Model(&Friend{}).Where("user_id = ? AND friend_id = ?", userID, friendID).Update("status", FriendStatusAccepted)

	return nil
}
