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
