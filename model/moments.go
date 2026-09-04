package model

import (
	"time"
)

/*
create table if not exists moment_likes (

	id bigint not null comment "点赞id" primary key auto_increment,
	moment_id bigint not null comment "朋友圈id",
	user_id bigint not null comment "点赞者id",
	created_at datetime default current_timestamp comment "创建时间",
	UNIQUE KEY uk_moment_user (moment_id, user_id)

) ENGINE=InnoDB default charset=utf8mb4 comment "朋友圈点赞表;记录谁给谁的朋友圈点赞";
*/

// 可见性，0-公开，1-好友可见，2-仅自己可见 3-部分好友可见
type MomentVisbileStatus int

const (
	MomentVisibleStatusPublic         MomentVisbileStatus = 0
	MomentVisibleStatusFriendsOnly    MomentVisbileStatus = 1
	MomentVisibleStatusSelfOnly       MomentVisbileStatus = 2
	MomentVisibleStatusPartialFriends MomentVisbileStatus = 3
)

// ====== 朋友圈主结构 ======
type Moment struct {
	// 朋友圈 ID
	ID uint `gorm:"primaryKey" json:"id"`
	// 朋友圈发布者 ID
	OwnerID uint `gorm:"column:owner_id;index" json:"ownerId"`
	// 朋友圈元素列表（JSON 序列化存储）
	Elements ElementList ` gorm:"type:json" json:"elements"`
	// 朋友圈状态，0-正常，1-删除
	Status int `gorm:"column:status;default:0" json:"status"`
	// 点赞数
	LikeCount int `gorm:"column:like_count;default:0" json:"likeCount"`
	// 可见性，0-公开，1-好友可见，2-仅自己可见 3-部分好友可见
	Visible MomentVisbileStatus `gorm:"column:visible;default:0" json:"visible"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *Moment) TableName() string {
	return "moments"
}

// swagger:model MomentPrivacy
type MomentLike struct {
	// 点赞 ID
	ID uint `gorm:"primaryKey" json:"id"`
	// 朋友圈记录 ID
	MomentID uint `gorm:"column:moment_id" json:"momentId"`
	// 点赞者 ID
	UserID uint `gorm:"column:user_id" json:"userId"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (m *MomentLike) TableName() string {
	return "moments_likes"
}

type MomentUser struct {
	// 用户中文名称，默认为 ‘’，可后期通过修改用户信息设置
	Nickname string `gorm:"column:nickname" json:"nickname"`

	// 用户头像、创建时；为空字符串;注意只是保存 文件的hash 地址；而不是 URL 地址
	Avatar string `gorm:"column:avatar" json:"avatar" example:"''" validate:"required"`
}

// swagger:model MomentPrivacy
type MomentLikeResponse struct {
	*MomentLike
	*MomentUser
}

type MomentVisible struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	MomentID uint `gorm:"column:moment_id" json:"momentId"`
	UserID   uint `gorm:"column:user_id" json:"userId"`
	Visible  int  `gorm:"column:visible;default:0" json:"visible"`
}

func (m *MomentVisible) TableName() string {
	return "moments_visible"
}

var (
	// 评论状态，0-正常，1-删除
	MomentCommentStatusNormal  = 0
	MomentCommentStatusDeleted = 1
)

// swagger:model MomentComments
type MomentComments struct {
	// 评论ID
	ID uint `gorm:"primaryKey" json:"id"`
	// 朋友圈记录ID
	MomentID uint `gorm:"column:moment_id" json:"momentId"`
	// 评论者用户ID
	UserID uint `gorm:"column:user_id" json:"userId"`
	// 评论内容
	Content string `gorm:"column:visible;default:0" json:"visible"`
	// 评论状态，0-正常，1-删除
	Status uint `gorm:"column:status;default:0" json:"status"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *MomentComments) TableName() string {
	return "moments_comments"
}

// swagger:model MomentCommentsResponse
type MomentCommentsResponse struct {
	*MomentComments
	*MomentUser
}

// swagger:model MomentPrivacy
type MomentPrivacy struct {
	// 隐私设置ID
	ID uint `gorm:"primaryKey" json:"id"`
	// 用户ID
	UserID uint `gorm:"column:user_id" json:"userId"`
	// 目标用户ID
	TargetID uint `gorm:"column:target_id" json:"targetId"`
	// 我不看TA的朋友圈，0-不屏蔽，1-屏蔽
	HideTheir bool `gorm:"column:hide_their;default:false" json:"hideTheir"`
	// 不让TA看我的朋友圈，0-不屏蔽，1-屏蔽
	HideMine bool `gorm:"column:hide_mine;default:false" json:"hideMine"`
}

func (m *MomentPrivacy) TableName() string {
	return "moments_privacy"
}
