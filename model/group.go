package model

import (
	"time"
)

// GroupChat 群组表
type GroupChat struct {
	// 群ID
	ID uint `gorm:"primarykey" json:"id"`
	// 群主ID
	OwnerID uint `gorm:"column:owner_id;index" json:"ownerId"`
	// 群名称
	Name string `gorm:"column:name;size:255;default:''" json:"name"`
	// 群头像
	Avatar string `gorm:"column:avatar;size:255;default:''" json:"avatar"`
	// 群公告
	Notice string `gorm:"column:notice;size:255;default:''" json:"notice"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	// 是否加密 0-否 1-是
	Encrypted int8 `gorm:"column:encrypted;default:0" json:"encrypted"`
	// 是否全员禁言 0-否 1-是
	IsMuted int8 `gorm:"column:is_muted;default:0" json:"isMuted"`
}

func (*GroupChat) TableName() string {
	return "group_chats"
}

// GroupMember 群成员表
type GroupMember struct {
	// 成员ID
	ID uint `gorm:"primarykey" json:"id"`
	// 群ID
	GroupID uint `gorm:"column:group_id;index;uniqueIndex:uk_group_user" json:"groupId"`
	// 用户ID
	UserID uint `gorm:"column:user_id;index;uniqueIndex:uk_group_user" json:"userId"`
	// 角色 0-成员 1-管理员 2-群主
	Role int8 `gorm:"column:role;default:0" json:"role"`
	// 加入时间
	JoinedAt time.Time `gorm:"column:joined_at;autoCreateTime" json:"joinedAt"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	// 该用户是否被禁言 0-否 1-是
	IsMuted int8 `gorm:"column:is_muted;default:0" json:"isMuted"`
	// 该用户是否关闭该群通知 0-否 1-是
	IsNotifyDisabled int8 `gorm:"column:is_notify_disabled;default:0" json:"isNotifyDisabled"`

	// 该群员状态 0-正常 1-退群(自己退、被踢出)
	Status int8 `gorm:"column:status;default:0" json:"status"`
}

type GroupMemberResponse struct {
	*GroupMember
	// 用户信息；注意：该字段是通过关联查询获取的，而不是直接存储在数据库中的

	// 用户账户；创建时使用；后期不可修改。
	Username string `gorm:"uniqueIndex;size:50" json:"username" example:"whx" validate:"required"`
	// 用户头像、创建时；为空字符串;注意只是保存 文件的hash 地址；而不是 URL 地址
	Avatar string `gorm:"column:avatar" json:"avatar" example:"''" validate:"required"`
	// 用户中文名称，默认为 ‘’，可后期通过修改用户信息设置
	Nickname string `gorm:"column:nickname" json:"nickname"`
}

type GroupChatResponse struct {
	*GroupChat
	// 用户 ID；如果 ownerId == userId，则表示该用户是群主
	// UserID UserID `gorm:"column:user_id;index;uniqueIndex:uk_group_user" json:"userId" swaggertype:"integer"`
	// 该群的所有成员信息；注意：该字段是通过关联查询获取的，而不是直接存储在数据库中的
	Members []*GroupMemberResponse `gorm:"foreignKey:GroupID;references:ID" json:"members"`
}

func (*GroupMember) TableName() string {
	return "group_members"
}
