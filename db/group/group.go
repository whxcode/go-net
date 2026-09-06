package dbGroup

import (
	"fmt"

	"go-net/db"
	"go-net/m_tools"
	"go-net/model"
)

type groupDB struct{}

var GroupDB = &groupDB{}

func (g *groupDB) groups(userID *[]model.UserID, groupID *[]uint) []*model.GroupChatResponse {
	var result []*model.GroupChatResponse
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

	query := db.DB.Raw(sql, parmas...)

	err := query.Find(&result).Error
	if err != nil {
		panic(err)
	}

	if result == nil || len(result) == 0 {
		return nil
	}

	groupIDs := make([]uint, len(result))

	groupMap := make(map[uint]*model.GroupChatResponse)

	for i, group := range result {
		groupIDs[i] = group.ID
		groupMap[group.ID] = group
		group.Members = []*model.GroupMemberResponse{}
	}

	var members []*model.GroupMemberResponse

	err = db.DB.Table("group_members as m").
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

func (g *groupDB) Groups(userID model.UserID) []*model.GroupChatResponse {
	return g.groups(&[]model.UserID{userID}, nil)
}

func (g *groupDB) GroupID(groupID uint) *model.GroupChatResponse {
	result := g.groups(nil, &[]uint{groupID})

	if result == nil || len(result) == 0 {
		return nil
	}

	return result[0]
}

func (g *groupDB) GroupMembers(groupID uint) []uint {
	var result []uint
	err := db.DB.Debug().Table("group_members m").
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

func (g *groupDB) PutGroup(group *model.GroupChat) *model.GroupChatResponse {
	err := db.DB.Model(&model.GroupChat{}).
		Where("id = ?", group.ID).
		Updates(group).Error
	if err != nil {
		panic(err)
	}

	return g.GroupID(group.ID)
}

func (g *groupDB) PostGroupMembers(groupID uint, MemberIDs []model.UserID) {
	members := make([]*model.GroupMember, len(MemberIDs))

	for i, memberID := range MemberIDs {
		members[i] = &model.GroupMember{
			GroupID: groupID,
			UserID:  uint(memberID),
		}
	}

	err := db.DB.Debug().Save(members).Error
	if err != nil {
		panic(err)
	}
}

func (g *groupDB) PutGroupMember(groupID uint, groupMember *model.GroupMember) *model.GroupMember {
	var result *model.GroupMember
	var group *model.GroupChat

	err := db.DB.Table("group_members").
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

	err = db.DB.Table("group_chats").
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

	err = db.DB.Save(result).Error
	if err != nil {
		panic(err)
	}

	return result
}
