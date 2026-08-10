// Package view 视图类型
package view

import (
	"fmt"

	"go-net/src/model"
)

func Test() {
	user := model.User{}
	user.Name = "王恒星"

	u1 := &user

	u1.Name = "ww"
	u1.Pet.Name = "whx"
	u1.Type = 1

	fmt.Println(user.SayName())
}
