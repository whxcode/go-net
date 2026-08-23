package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ====== 消息元素类型 ======
type MessageType uint

const (
	MsgTypeText  MessageType = 0
	MsgTypeImage MessageType = 1
	MsgTypeVideo MessageType = 2
	MsgTypeFile  MessageType = 3
)

// ====== 单条元素 ======
type Element struct {
	// 消息类型（0：文本，1：图片，2：视频，3：文件）
	Type MessageType `json:"type"`
	// 消息内容（文本内容或文件 HASH）
	Content string `json:"content,omitempty"`
	// 文件信息（仅在 Type 为图片/视频/文件时有效）
	Hash   string `json:"hash,omitempty"`
	Url    string `json:"url,omitempty"`
	Name   string `json:"name,omitempty"`
	Size   int64  `json:"size,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// ====== 元素切片（支持 JSON 序列化）======
type ElementList []*Element

func (e *ElementList) Scan(value interface{}) error {
	if value == nil {
		*e = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, e)
}

func (e ElementList) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	return json.Marshal(e)
}

// ====== 通道消息类型（不存数据库）======
type ChannelType uint8

const (
	ChannelTypeNORMAL ChannelType = iota
	ChannelTypePING
	ChannelTypePONG
)

// ====== 消息主结构 ======
type Message struct {
	// 消息类型 0：普通消息 1：ping 2：pong
	Type ChannelType `json:"type" gorm:"-"`
	// 消息 数据库索引
	ID uint `gorm:"primarykey" json:"id"`
	// 消息唯一 ID（客户端生成，保证幂等性）
	MsgID string `gorm:"column:msg_id;uniqueIndex;size:64" json:"msgId"` // varchar(64)
	// 消息发送者 ID
	SenderID uint `gorm:"column:sender_id;index" json:"senderId"`
	// 消息接收者 ID
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
	return DB.Create(msg).Error
}

// 批量保存
func (*messageDB) SaveBatch(msgs []*Message) error {
	if len(msgs) == 0 {
		return nil
	}
	return DB.Create(msgs).Error
}

// 查询两人聊天记录
func (*messageDB) GetHistory(userId, friendId UserID, limit, offset int) ([]*Message, int) {
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
