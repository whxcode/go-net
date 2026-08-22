package controller

import (
	"net/http"
	"strconv"

	"go-net/model"
	"go-net/redis"
	"go-net/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userControll struct{}

var UserControll = &userControll{}

// @Summary 获取用户信息
// @Tags 用户
// @Success 200 {object} utils.KResponse{data=model.UserResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /users/get [get]
func (*userControll) GetUser(c *gin.Context) *utils.KResponse {
	token := utils.GetToken(c)
	u := utils.GetUserID(c)

	user, err := model.UserDb.GetUserByUserID(u)
	if err != nil {
		return utils.MakeResponseWidthCode("获取用户信息失败", http.StatusInternalServerError)
	}

	return utils.MakeResponse(&model.UserResponse{
		User:  *user,
		Token: token,
	})
}

// @Summary  根据用户ID获取用户信息
// @Tags 用户
// @Param id path int true "用户ID"
// @Success 200 {object} utils.KResponse{data=model.UserResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /users/:id [get]
func (*userControll) GetUserByID(c *gin.Context) *utils.KResponse {
	token := utils.GetToken(c)
	u := utils.GetUserID(c)

	user, err := model.UserDb.GetUserByUserID(u)
	if err != nil {
		return utils.MakeResponseWidthCode("获取用户信息失败", http.StatusInternalServerError)
	}

	return utils.MakeResponse(&model.UserResponse{
		User:  *user,
		Token: token,
	})
}

// @Summary 根据 search 关键字 获取用户列表
// @Tags 用户
// @Param search query string false "搜索关键字" default("")
// @Success 200 {object} utils.KResponse{data=[]model.User} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /users/users [get]
func (*userControll) GetUsers(c *gin.Context) *utils.KResponse {
	search := c.Query("search")

	users := model.UserDb.GetUsers(search)

	if users == nil {
		return utils.MakeResponseWidthCode("获取用户信息失败", http.StatusInternalServerError)
	}

	return utils.MakeResponse(users)
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary 用户注册
// @Tags 用户
// @Param request body RegisterRequest true "注册请求"
// @Success 200 {object} utils.KResponse{data=[]model.User} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /users/register [post]
func (*userControll) Register(c *gin.Context) *utils.KResponse {
	user := &RegisterRequest{}

	if err := c.ShouldBindBodyWithJSON(user); err != nil {
		return utils.MakeResponseWidthCode("无效的参数", http.StatusBadRequest)
	}

	if model.UserDb.UserExists(user.Username) {
		return utils.MakeResponseWidthCode("用户已存在", http.StatusBadRequest)
	}

	hashed := utils.GenerateFromPasswordString(user.Password)
	user.Password = hashed

	model.UserDb.AddUser(&model.User{
		Username: user.Username,
		Password: user.Password,
	})

	user.Password = ""

	return utils.MakeResponse(user)
}

// @Summary 用户登录
// @Tags 用户
// @Param request body RegisterRequest true "登录请求"
// @Success 200 {object} utils.KResponse{data=[]model.UserResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /users/login [post]
func (*userControll) Login(c *gin.Context) *utils.KResponse {
	reqUser := &RegisterRequest{}

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

	err = redis.User.SetToken(token, user.ID)
	if err != nil {
		return utils.MakeResponseWidthCode("保存token失败", http.StatusInternalServerError)
	}

	return utils.MakeResponse(&model.UserResponse{
		User:  *user,
		Token: token,
	})
}

// @Summary 退出登录
// @Tags 用户
// @Success 200 {object} utils.KResponse{data=string} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /users/logout [get]
func (*userControll) Logout(c *gin.Context) *utils.KResponse {
	token := utils.GetToken(c)
	err := redis.User.DelToken(token)
	if err != nil {
		return utils.MakeResponseWidthCode("退出失败", http.StatusInternalServerError)
	}

	return utils.MakeResponse("退出成功")
}

type UserPutPasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// @Summary 修改用户密码;
// @Description 成功后；服务端将清除当前用户token的登录状态
// @Tags 用户
// @Param request body UserPutPasswordRequest true "修改密码请求"
// @Success 200 {object} utils.KResponse{data=string} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /users/password [put]
func (*userControll) UserPutPassword(c *gin.Context) *utils.KResponse {
	token := utils.GetToken(c)
	userID := utils.GetUserID(c)

	data := &UserPutPasswordRequest{}

	if err := c.ShouldBindBodyWithJSON(data); err != nil {
		return utils.MakeResponseWidthCode("无效的参数", http.StatusBadRequest)
	}

	user, err := model.UserDb.GetUserByUserID(userID)
	if err != nil {
		return utils.MakeResponseWidthCode("用户不存在", http.StatusBadRequest)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.OldPassword))
	if err != nil {
		return utils.MakeResponseWidthCode(gin.H{
			"data": data,
			"err":  err,
			"user": user,
		}, http.StatusBadRequest)
	}

	newPassword := utils.GenerateFromPasswordString(data.NewPassword)

	/*
		err := redis.RedisClient.Del(redis.Ctx, token.(string)).Err()
		if err != nil {
			return utils.MakeResponseWidthCode("退出失败", http.StatusInternalServerError)
		}
	*/

	return utils.MakeResponse(gin.H{
		"msg":         "修改成功" + token + ":" + strconv.Itoa(int(userID)),
		"data":        data,
		"newPassword": newPassword,
	})
}
