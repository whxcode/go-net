package model

import (
	"strconv"
	"time"
)

// UserID 用户ID
type UserID uint

const InvalidUserID UserID = 0

func (id UserID) String() string {
	return strconv.FormatUint(uint64(id), 10) // ✅ "4"
}

func (id UserID) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(id), 10)), nil
}

type User struct {
	// 用户 ID;其类型是一个 uint 类型
	ID UserID `gorm:"primarykey" json:"id" validate:"required" swaggertype:"integer"`
	// 用户账户；创建时使用；后期不可修改。
	Username string `gorm:"uniqueIndex;size:50" json:"username" example:"whx" validate:"required"`
	// 用户头像、创建时；为空字符串;注意只是保存 文件的hash 地址；而不是 URL 地址
	Avatar string `gorm:"column:avatar" json:"avatar" example:"''" validate:"required"`
	// 用户中文名称，默认为 ‘’，可后期通过修改用户信息设置
	Nickname  string    `gorm:"column:nickname" json:"nickname"`
	Password  string    `json:"-" comment:"密码（不返回给前端）"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createAt" comment:"创建时间"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updateAt" comment:"更新时间"`
}

type UserResponse struct {
	User
	// 认证令牌
	Token string `json:"token" validate:"required"`
}
