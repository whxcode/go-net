package model

import "time"

// == 时间线表
type TimeLine struct {
	// 时间线ID
	ID uint `gorm:"column:id;primaryKey;autoIncrement;comment:时间线ID" json:"id"`
	// 时间线发布者id
	OwnerID uint `gorm:"column:owner_id;not null;comment:时间线发布者id;index:idx_owner_id" json:"owner_id"`
	// 时间线内容
	Elements ElementList `gorm:"column:elements;type:json;not null;comment:时间线内容" json:"elements"`
	// 点赞数
	LikeCount int `gorm:"column:like_count;default:0;comment:点赞数" json:"like_count"`
	// 时间线状态，0-正常，1-删除,2-被举报中,3-举报成功,4-举报失败
	Status int8 `gorm:"column:status;default:0;comment:时间线状态，0-正常，1-删除,2-被举报中,3-举报成功,4-举报失败" json:"status"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	// 更新时间
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (TimeLine) TableName() string {
	return "time_lines"
}

// == 时间线评论表
type TimeComment struct {
	// 评论id
	ID uint `gorm:"column:id;primaryKey;autoIncrement;comment:评论id" json:"id"`
	// 时间线id
	TimeLineID uint `gorm:"column:time_line_id;not null;comment:时间线id;index:idx_time_line_id" json:"time_line_id"`
	// 评论者id
	OwnerID uint `gorm:"column:owner_id;not null;comment:评论者id;index:idx_owner_id" json:"owner_id"`
	// 评论内容
	Elements ElementList `gorm:"column:elements;type:json;not null;comment:评论内容" json:"elements"`
	// 评论状态，0-正常，1-删除
	Status int8 `gorm:"column:status;default:0;comment:评论状态，0-正常，1-删除" json:"status"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	// 更新时间
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (TimeComment) TableName() string {
	return "time_comments"
}

// == 时间线点赞表
type TimeLike struct {
	// 点赞id
	ID uint `gorm:"column:id;primaryKey;autoIncrement;comment:点赞id" json:"id"`
	// 时间线id
	TimeLineID uint `gorm:"column:time_line_id;not null;comment:时间线id;index:idx_time_line_id;uniqueIndex:uk_time_line_owner" json:"time_line_id"`
	// 点赞者id
	OwnerID uint `gorm:"column:owner_id;not null;comment:点赞者id;index:idx_owner_id;uniqueIndex:uk_time_line_owner" json:"owner_id"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
}

// TableName 指定表名
func (TimeLike) TableName() string {
	return "time_likes"
}
