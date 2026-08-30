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

type groupDB struct{}

var GroupDB = &groupDB{}

func (db *groupDB) groups(userID *[]UserID, groupID *[]uint) []*GroupChatResponse {
	var result []*GroupChatResponse
	var parmas []interface{} = make([]interface{}, 0)

	sql := `WITH  m as (select distinctrow group_id from group_members where status = 0`

	if userID != nil {

		usrUint := make([]uint, len(*userID))
		for i, id := range *userID {
			usrUint[i] = uint(id)
		}
		sql += ` and user_id in (?)) `
		parmas = append(parmas, usrUint)
	} else {
		sql += `) `
	}

	sql += `select c.* from m left join group_chats c on m.group_id = c.id`

	if groupID != nil {
		sql += ` where c.id in (?)`
		parmas = append(parmas, *groupID)
	}

	query := DB.Raw(sql, parmas...)

	err := query.Find(&result).Error
	if err != nil {
		panic(err)
	}

	if result == nil || len(result) == 0 {
		return nil
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
		Where("m.group_id IN ?  and m.status = ?", m_tools.UniqueSlice(groupIDs), 0).Find(&members).Error
	if err != nil {
		panic(err)
	}

	for _, member := range members {
		if group, ok := groupMap[member.GroupID]; ok {
			group.Members = append(group.Members, member)
		}
	}

	return result
}

func (db *groupDB) Groups(userID UserID) []*GroupChatResponse {
	return db.groups(&[]UserID{userID}, nil)
}

func (db *groupDB) GroupID(groupID uint) *GroupChatResponse {
	result := db.groups(nil, &[]uint{groupID})

	if result == nil || len(result) == 0 {
		return nil
	}

	return result[0]
}

func (db *groupDB) GroupMembers(groupID uint) []uint {
	var result []uint
	err := DB.Debug().Table("group_members m").
		Select("m.user_id").
		Where("m.group_id = ? and m.status = 0", groupID).
		Find(&result).
		Error
	if err != nil {
		fmt.Println("GroupMembers error:", err)
		panic(err)
	}

	return result
}

func (db *groupDB) PutGroup(group *GroupChat) *GroupChatResponse {
	err := DB.Model(&GroupChat{}).
		Where("id = ?", group.ID).
		Updates(group).Error
	if err != nil {
		panic(err)
	}

	return db.GroupID(group.ID)
}

func (db *groupDB) PostGroupMembers(groupID uint, MemberIDs []UserID) {
	members := make([]*GroupMember, len(MemberIDs))

	for i, memberID := range MemberIDs {
		members[i] = &GroupMember{
			GroupID: groupID,
			UserID:  uint(memberID),
		}
	}

	err := DB.Debug().Save(members).Error
	if err != nil {
		panic(err)
	}
}

func (db *groupDB) PutGroupMember(groupID uint, groupMember *GroupMember) *GroupMember {
	var result *GroupMember
	var group *GroupChat

	err := DB.Table("group_members").
		Where("group_id = ? AND user_id = ?", groupID, groupMember.UserID).
		Find(&result).
		Error
		// err := DB.Raw("update table group_members", groupID, groupMember.UserID).Delete(&GroupMember{}).Error
	if err != nil {
		panic(err)
	}

	if result == nil {
		panic(fmt.Sprintf("group member not found for groupID: %d, userID: %d", groupID, groupMember.UserID))
	}

	err = DB.Table("group_chats").
		Where("id = ?", groupID).Find(&group).Error
	if err != nil {
		panic(err)
	}

	if group == nil {
		panic(fmt.Sprintf("group not found for groupID: %d", groupID))
	}

	result.IsMuted = groupMember.IsMuted
	result.IsNotifyDisabled = groupMember.IsNotifyDisabled

	if group.OwnerID != groupMember.UserID {
		result.Status = groupMember.Status
		result.Role = groupMember.Role
	}

	err = DB.Save(result).Error
	if err != nil {
		panic(err)
	}

	return result
}
