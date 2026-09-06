package dbUser

import (
	"fmt"

	"go-net/db"
	"go-net/model"
)

type userDB struct{}

var UserDB = &userDB{}

func (u *userDB) AddUser(user *model.User) error {
	return db.DB.Create(user).Error
}

func (u *userDB) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := db.DB.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userDB) GetUserByUserID(userID model.UserID) *model.User {
	user := &model.User{}
	err := db.DB.Where("id = ?", userID).First(user).Error

	fmt.Println("err:", err)
	if err != nil {
		panic(err)
	}

	return user
}

func (u *userDB) UserExists(username string) bool {
	result := db.DB.Where("username = ?", username).First(&model.User{})

	return result.Error == nil
}

/*
* 模糊匹配用户名称
*
* */
func (u *userDB) GetUsers(username string) (result []*model.User) {
	err := db.DB.Where("username LIKE ?", "%"+username+"%").Find(&result).Error
	if err != nil {
		panic(err)
	}

	return
}

func (u *userDB) UpdateUser(user *model.User) error {
	err := db.DB.Table("users").Updates(user).Where("id = ?", user.ID).Error
	if err != nil {
		panic(err)
	}

	return nil
}

func (u *userDB) UpdateUserAvatar(userID model.UserID, avatar string) error {
	err := db.DB.Table("users").Where("id = ?", userID).
		Update("avatar", avatar).
		Error
	if err != nil {
		panic(err)
	}

	return nil
}

func (u *userDB) UpdateUserNickname(userID model.UserID, nickname string) error {
	err := db.DB.Table("users").Where("id = ?", userID).
		Update("nickname", nickname).
		Error
	if err != nil {
		panic(err)
	}

	return nil
}
