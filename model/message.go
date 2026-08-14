package model

import "time"

type MessageType uint

var (
	MessageTypeText  MessageType = 0
	MessageTypeImage MessageType = 0
)

type Element struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content,omitempty"`
	Hash    string      `json:"hash,omitempty"`
}

type Message struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	MsgID      string    `gorm:"column:msg_id;uniqueIndex;size:36" json:"msgId"`
	SenderID   uint      `gorm:"column:sender_id" json:"senderId"`
	ReceiverID uint      `gorm:"column:receiver_id" json:"receiverId"`
	Elements   []Element `gorm:"type:json" json:"elements"`
	Status     int       `gorm:"column:status;default:0" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (Message) TableName() string {
	return "messages"
}

type messageDB struct{}

var MessageDB = &messageDB{}

func (*messageDB) Save(messages *Message) error {
	return DB.Create(messages).Error
}

func (*messageDB) SaveBatch(messages *[]Message) error {
	return DB.Create(messages).Error
}

func (*messageDB) GetMessagesBySenderAndReceiver(senderID, receiverID uint) ([]*Message, error) {
	var messages []*Message

	err := DB.Where("sender_id = ? AND receiver_id = ?", senderID, receiverID).Find(&messages).Error

	return messages, err
}

func (*messageDB) GetMessagesByReceiverAndSender(receiverID, senderID uint) ([]*Message, error) {
	return MessageDB.GetMessagesBySenderAndReceiver(senderID, receiverID)
}
