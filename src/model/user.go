// Package model 提供用户数据模型
package model

type Pet struct {
	Name string
	Type uint8
}

type User struct {
	Name string
	Pet
}

func (u *User) SayName() string {
	return u.Name
}
