package model

import (
	"database/sql/driver"
	"encoding/json"
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
