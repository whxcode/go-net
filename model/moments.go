package model

import "time"

/*
create table if not exists moment_likes (

	id bigint not null comment "点赞id" primary key auto_increment,
	moment_id bigint not null comment "朋友圈id",
	user_id bigint not null comment "点赞者id",
	created_at datetime default current_timestamp comment "创建时间",
	UNIQUE KEY uk_moment_user (moment_id, user_id)

) ENGINE=InnoDB default charset=utf8mb4 comment "朋友圈点赞表;记录谁给谁的朋友圈点赞";
*/

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
	Visible int `gorm:"column:visible;default:0" json:"visible"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// ===== 朋友圈点赞结构 ======
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

/**
create table if not exists moments_visible (

	id bigint not null comment "朋友圈id" primary key auto_increment,
	moment_id bigint not null comment "朋友圈id",
	user_id bigint not null comment "用户 id",
	visible int default 0 comment "0 该好友可见,1 该好友不可见;需要配合coments.visible = 3 的情况",

	UNIQUE KEY uk_moment_user (moment_id, user_id)

) ENGINE=InnoDB default charset=utf8mb4 comment "该条记录谁可见；谁不可见表";
* */

// ===== 朋友圈可见性结构 ======
type MomentVisible struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	MomentID uint `gorm:"column:moment_id" json:"momentId"`
	UserID   uint `gorm:"column:user_id" json:"userId"`
	Visible  int  `gorm:"column:visible;default:0" json:"visible"`
}

/**
*
朋友圈评论表
create table if not exists moments_comments (
	id bigint not null comment "评论id" primary key auto_increment,
	moment_id bigint not null comment "朋友圈id",
	user_id bigint not null comment "评论者id",
	content text not null comment "评论内容",
	status tinyint default 0 comment "评论状态，0-正常，1-删除",
	created_at datetime default current_timestamp comment "创建时间",
	updated_at datetime default current_timestamp on update current_timestamp comment "更新时间"

) ENGINE=InnoDB default charset=utf8mb4 comment "朋友圈评论表";
*
* */

var (
	// 评论状态，0-正常，1-删除
	MomentCommentStatusNormal  = 0
	MomentCommentStatusDeleted = 1
)

// ===== 朋友圈可见性结构 ======
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

/*
*
*
create table if not exists moment_privacy (

	id bigint not null primary key auto_increment comment "隐私设置id",
	user_id bigint not null comment "用户id",
	target_id bigint not null comment "目标用户id",
	hide_their tinyint default 0 comment "我不看TA的朋友圈，0-不屏蔽，1-屏蔽",
	hide_mine tinyint default 0 comment "不让TA看我的朋友圈，0-不屏蔽，1-屏蔽",
	UNIQUE KEY uk_user_target (user_id, target_id)

) ENGINE=InnoDB default charset=utf8mb4 comment "朋友圈隐私表;记录谁屏蔽谁的朋友圈";
* */

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
