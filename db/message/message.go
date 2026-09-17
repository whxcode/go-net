package dbMessage

import (
	"fmt"
	"sync"
	"time"

	"go-net/db"
	"go-net/model"
)

type messageDB struct {
	mutex sync.RWMutex
	msgs  []*model.Message
}

var MessageDB *messageDB

func init() {
	MessageDB = &messageDB{
		msgs: make([]*model.Message, 0),
	}

	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				{
					if len(MessageDB.msgs) == 0 {
						continue
					}
					MessageDB.mutex.RLock()
					db.DB.Create(MessageDB.msgs)
					MessageDB.msgs = make([]*model.Message, 0)
					MessageDB.mutex.RUnlock()
					fmt.Printf("---------------------将数据写入数据库--------------")

				}
			}
		}
	}()
}

// 保存单条
func (m *messageDB) Save(msg *model.Message) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.msgs = append(m.msgs, msg)

	if len(m.msgs) > 1000 {
		db.DB.Create(m.msgs)
		m.msgs = make([]*model.Message, 0)
	}
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
