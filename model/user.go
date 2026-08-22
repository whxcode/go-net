package model

import (
	"strconv"
	"sync"
	"time"
)

type UserID uint

const InvalidUserID UserID = 0

func (id UserID) String() string {
	return string(id)
}

func (id UserID) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(id), 10)), nil
}

type User struct {
	// 用户 ID
	ID UserID `gorm:"primarykey" json:"id" example:"1" validate:"required"`
	// 用户账户；创建时使用；后期不可修改。
	Username string `gorm:"uniqueIndex;size:50" json:"username" example:"whx" validate:"required"`
	// 用户头像、创建时；为空字符串
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

type userDb struct {
	users []*User
	mutex sync.Mutex
}

var UserDb = &userDb{}

func (db *userDb) AddUser(user *User) error {
	return DB.Create(user).Error
}

func (db *userDb) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := DB.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *userDb) GetUserByUserID(userID UserID) (*User, error) {
	user := &User{}
	err := DB.Where("id = ?", userID).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *userDb) UserExists(username string) bool {
	result := DB.Where("username = ?", username).First(&User{})

	return result.Error == nil
}

/*
* 模糊匹配用户名称
*
* */
func (db *userDb) GetUsers(username string) (result []*User) {
	err := DB.Where("username LIKE ?", "%"+username+"%").Find(&result).Error
	if err != nil {
		panic(err)
	}

	return
}
