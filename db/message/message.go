package dbMessage

import (
	"go-net/db"
	"go-net/model"
)

type messageDB struct{}

var MessageDB = &messageDB{}

// 保存单条
func (*messageDB) Save(msg *model.Message) error {
	err := db.DB.Create(msg).Error
	if err != nil {
		panic(err)
	}

	return nil
}

// 批量保存
func (*messageDB) SaveBatch(msgs []*model.Message) error {
	if len(msgs) == 0 {
		return nil
	}
	return db.DB.Create(msgs).Error
}

// 查询两人聊天记录
func (*messageDB) GetFriendsHistory(userId, friendId model.UserID, limit, offset int) ([]*model.Message, int) {
	var messages []*model.Message
	var total int64 = 0
	// 1. 先查总数
	err := db.DB.Model(&model.Message{}).
		Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userId, friendId, friendId, userId,
		).
		Count(&total).Error

	err = db.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userId, friendId, friendId, userId,
	).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	if err != nil {
		panic(err)
	}

	return messages, int(total)
}

// 查询未读消息
func (*messageDB) GetUnread(userID uint) ([]model.Message, error) {
	var messages []model.Message
	err := db.DB.Where("receiver_id = ? AND status = ?", userID, 0).
		Where("type != ?", model.ChannelTypePING).
		Where("type != ?", model.ChannelTypePONG).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

// 标记已读
func (*messageDB) MarkAsRead(msgID string) error {
	return db.DB.Model(&model.Message{}).Where("msg_id = ?", msgID).Update("status", 1).Error
}

// 查询两人聊天记录
func (*messageDB) GetGroupHistory(groupID model.UserID, limit, offset int) ([]*model.Message, int) {
	var messages []*model.Message
	var total int64 = 0
	// 1. 先查总数
	err := db.DB.Model(&model.Message{}).
		Where(
			"(receiver_id = ?)",
			groupID).
		Count(&total).Error

	err = db.DB.Where(
		"(receiver_id = ?)", groupID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	if err != nil {
		panic(err)
	}

	return messages, int(total)
}
