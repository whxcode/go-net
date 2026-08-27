package model

import (
	"fmt"
	"time"

	"go-net/m_tools"
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

type groupDB struct{}

var GroupDB = &groupDB{}

func (db *groupDB) groups(userID *[]UserID, groupID *[]uint) []*GroupChatResponse {
	var result []*GroupChatResponse

	// 2 次查询
	query := DB.Debug().Table("group_chats as c").
		Select("c.*")

	if userID != nil {
		query = query.
			Select("c.*,m.user_id").
			Joins("left join group_members as m on c.id = m.group_id").
			Where("user_id IN ?", *userID)
	}

	if groupID != nil {
		query = query.Where("c.id IN ?", *groupID)
	}

	err := query.Find(&result).Error
	if err != nil {
		panic(err)
	}

	groupIDs := make([]uint, len(result))

	groupMap := make(map[uint]*GroupChatResponse)

	for i, group := range result {
		groupIDs[i] = group.ID
		groupMap[group.ID] = group
		group.Members = []*GroupMemberResponse{}
	}

	var members []*GroupMemberResponse

	err = DB.Table("group_members as m").
		Select("m.*,u.username,u.nickname,u.avatar").
		Joins("left join users as u on m.user_id = u.id").
		Where("m.group_id IN ?", m_tools.UniqueSlice(groupIDs)).Find(&members).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("members--", result, groupIDs, members)

	for _, member := range members {
		if group, ok := groupMap[member.GroupID]; ok {
			group.Members = append(group.Members, member)
		}
	}

	return result
}

func (db *groupDB) Groups(userID UserID) []*GroupChatResponse {
	return db.groups(&[]UserID{userID, 2}, nil)
}

func (db *groupDB) GroupID(groupID uint) *GroupChatResponse {
	return db.groups(nil, &[]uint{groupID})[0]
}
