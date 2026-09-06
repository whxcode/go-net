package dbFriend

import (
	"errors"

	"gorm.io/gorm"

	"go-net/db"
	"go-net/model"
)

type friendDB struct{}

var FriendDB = &friendDB{}

// 根据用户id 查询好用列表
func (f *friendDB) Firends(userId model.UserID) []*model.FriendResponse {
	var result []*model.FriendResponse
	db.DB.Exec("SET @user_id = ?", userId)
	db.DB.Exec("SET @status = ?", model.FriendStatusAccepted)

	err := db.DB.Raw(
		`
		select
		f.id,
		case when f.user_id = @user_id then f.user_id else f.friend_id end as user_id,
		case when f.user_id = @user_id then f.friend_id else f.user_id end as friend_id,
		u.username,
		u.nickname,
		u.avatar
		from friends f
		left join users u on u.id = (case when f.user_id = @user_id then f.friend_id else f.user_id end)
		where (f.user_id = @user_id or f.friend_id = @user_id) and f.status = @status
	`, userId, model.FriendStatusAccepted).Scan(&result).Error
	if err != nil {
		panic(err)
	}

	if result == nil {
		result = make([]*model.FriendResponse, 0)
	}

	return result
}

func (f *friendDB) Requests(userID model.UserID) (result []*model.FriendResponse) {
	err := db.DB.Table("friends f").
		Select("f.id, f.user_id, f.friend_id, u.username, u.nickname, u.avatar").
		Joins("left join users u on u.id = (case when f.user_id = ? then f.friend_id else f.user_id end)", userID).
		Where("(f.user_id = ? or f.friend_id = ?) and f.status in (?)", userID, userID, []model.RequestFriendStatus{model.FriendStatusPending, model.FriendStatusRejected}).
		Find(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

/*
* 好友申请
* user_id -> 我
* friendId -> 对方
*
* */
func (*friendDB) Request(userId, friendId model.UserID, remark string) *model.Friend {
	f := &model.Friend{
		UserID:   userId,
		FriendID: friendId,
		Status:   model.FriendStatusPending,
		Remark:   remark,
	}
	/*
		// First
		err := DB.First(&user, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
		    // 查不到
		} else if err != nil {
		    // 报错
		}

		// Find
		err := DB.Find(&users).Error
		if err != nil {
		    // 报错
		} else if len(users) == 0 {
		    // 查不到
		}
	*/

	var existingFriend *model.Friend

	err := db.DB.Table("friends").Where("user_id = ? AND friend_id = ?", userId, friendId).First(&existingFriend).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = db.DB.Create(f).Error
		if err != nil {
			panic(err)
		}

		return f
	}

	if err != nil {
		panic(err)
	}

	switch existingFriend.Status {
	case model.FriendStatusAccepted, model.FriendStatusPending:
	case model.FriendStatusRejected, model.FriendStatusDeleted:

		existingFriend.Status = model.FriendStatusPending
		existingFriend.Remark = remark
		err = db.DB.Save(existingFriend).Error
		if err != nil {
			panic(err)
		}

	}

	return existingFriend
}

func (*friendDB) PutRequestFriendStatus(id uint, status model.RequestFriendStatus) {
	var f *model.Friend

	err := db.DB.Table("friends").Where("id = ?", id).Find(&f).Error
	if err != nil {
		panic(err)
	}

	if f == nil {
		panic("好友申请不存在")
	}

	if status == model.FriendStatusDeleted {
		err = db.DB.Table("friends").
			Where("user_id = ? AND friend_id = ?", f.FriendID, f.UserID).
			Update("status", status).Error
		if err != nil {
			panic(err)
		}

	}

	// 查询时;需要交互一下 查询条件条件
	err = db.DB.Table("friends").Where("id = ?", id).
		Update("status", status).Error
	if err != nil {
		panic(err)
	}
}
