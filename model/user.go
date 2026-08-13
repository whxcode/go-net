package model

import (
	"sync"
	"time"
)

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50" json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updateAt"`
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

func (db *userDb) GetUserByUserID(userID uint) (*User, error) {
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
