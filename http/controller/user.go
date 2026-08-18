package controller

import (
	"fmt"
	"net/http"
	"time"

	"go-net/model"
	"go-net/model/redis"
	"go-net/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userControll struct{}

var UserControll = &userControll{}

func (*userControll) GetUser(c *gin.Context) *utils.KResponse {
	userId, _ := c.Get("userId")
	token, _ := c.Get("token")

	u, o := userId.(uint)

	if !o {
		return utils.MakeResponseWidthCode("获取用户信息失败", http.StatusInternalServerError)
	}

	user, err := model.UserDb.GetUserByUserID(u)
	if err != nil {
		return utils.MakeResponseWidthCode("获取用户信息失败", http.StatusInternalServerError)
	}

	return utils.MakeResponse(gin.H{
		"id":       user.ID,
		"username": user.Username,
		"token":    token,
	})
}

func (*userControll) Register(c *gin.Context) *utils.KResponse {
	user := &model.User{}

	if err := c.ShouldBindBodyWithJSON(user); err != nil {
		return utils.MakeResponseWidthCode("无效的参数", http.StatusBadRequest)
	}

	if model.UserDb.UserExists(user.Username) {
		return utils.MakeResponseWidthCode("用户已存在", http.StatusBadRequest)
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)

	model.UserDb.AddUser(user)

	user.Password = ""

	return utils.MakeResponse(user)
}

func (*userControll) Login(c *gin.Context) *utils.KResponse {
	reqUser := &model.User{}

	if err := c.ShouldBindBodyWithJSON(reqUser); err != nil {
		return utils.MakeResponseWidthCode("无效的参数", http.StatusBadRequest)
	}

	user, err := model.UserDb.GetUserByUsername(reqUser.Username)
	if err != nil {
		return utils.MakeResponseWidthCode("用户不存在", http.StatusBadRequest)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password))
	if err != nil {
		return utils.MakeResponseWidthCode("密码错误", http.StatusBadRequest)
	}

	token, err := utils.GenerateToken(uint(user.ID))
	if err != nil {
		return utils.MakeResponseWidthCode("生成token失败", http.StatusInternalServerError)
	}

	err = redis.RedisClient.Set(redis.Ctx, token, uint(user.ID), 10*time.Minute).Err()
	if err != nil {
		fmt.Println(err)
		return utils.MakeResponseWidthCode("保存token失败", http.StatusInternalServerError)
	}

	return utils.MakeResponse(gin.H{
		"id":       user.ID,
		"username": user.Username,
		"token":    token,
	})
}

func (*userControll) Logout(c *gin.Context) *utils.KResponse {
	token, _ := c.Get("token")

	err := redis.RedisClient.Del(redis.Ctx, token.(string)).Err()
	if err != nil {
		return utils.MakeResponseWidthCode("退出失败", http.StatusInternalServerError)
	}

	return utils.MakeResponse("退出成功")
}
