package model

import (
	"time"
)

// ====== 通道消息类型（不存数据库）======
type ChannelType uint8

const (
	// 好友之间消息  SenderID(UserID) -> ReceiverID (UserID)
	ChannelTypeFriend ChannelType = iota
	// 群组消息  SenderID(UserID) -> ReceiverID (GroupID)
	ChannelTypeGroup
	// 0 PING
	ChannelTypePING
	// 1 PONG
	ChannelTypePONG
)

// ====== 消息主结构 ======
type Message struct {
	// 消息类型 (0：好友消息，1：群组消息，2：PING，3：PONG)
	Type ChannelType `gorm:"column:type;default:0" json:"type"`
	// 消息 数据库索引
	ID uint `gorm:"primarykey" json:"id"`
	// 消息发送者 ID
	SenderID uint `gorm:"column:sender_id;index" json:"senderId"`
	// 消息接收者 ID, 如果是好友消息则为好友的 UserID，如果是群组消息则为群组的 GroupID
	ReceiverID uint `gorm:"column:receiver_id;index" json:"receiverId"`
	// 消息元素列表（JSON 序列化存储）
	Elements  ElementList `gorm:"type:json" json:"elements"`
	Status    int         `gorm:"column:status;default:0" json:"status"`
	CreatedAt time.Time   `gorm:"column:created_at" json:"createdAt"`
}

func (Message) TableName() string {
	return "messages"
}

// ====== 数据库操作 ======

type messageDB struct{}

var MessageDB = &messageDB{}

// 保存单条
func (*messageDB) Save(msg *Message) error {
	err := DB.Create(msg).Error
	if err != nil {
		panic(err)
	}

	return nil
}

// 批量保存
func (*messageDB) SaveBatch(msgs []*Message) error {
	if len(msgs) == 0 {
		return nil
	}
	return DB.Create(msgs).Error
}

// 查询两人聊天记录
func (*messageDB) GetFriendsHistory(userId, friendId UserID, limit, offset int) ([]*Message, int) {
	var messages []*Message
	var total int64 = 0
	// 1. 先查总数
	err := DB.Model(&Message{}).
		Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userId, friendId, friendId, userId,
		).
		Count(&total).Error

	err = DB.Where(
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
func (*messageDB) GetUnread(userID uint) ([]Message, error) {
	var messages []Message
	err := DB.Where("receiver_id = ? AND status = ?", userID, 0).
		Where("type != ?", ChannelTypePING).
		Where("type != ?", ChannelTypePONG).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

// 标记已读
func (*messageDB) MarkAsRead(msgID string) error {
	return DB.Model(&Message{}).Where("msg_id = ?", msgID).Update("status", 1).Error
}

// 查询两人聊天记录
func (*messageDB) GetGroupHistory(groupID UserID, limit, offset int) ([]*Message, int) {
	var messages []*Message
	var total int64 = 0
	// 1. 先查总数
	err := DB.Model(&Message{}).
		Where(
			"(receiver_id = ?)",
			groupID).
		Count(&total).Error

	err = DB.Where(
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
