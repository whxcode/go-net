package controller

import (
	"net/http"
	"strconv"

	dbUser "go-net/db/user"
	"go-net/model"
	UserRedis "go-net/redis/user"
	"go-net/response"
	"go-net/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userControll struct{}

var UserControll = &userControll{}

// @Summary 获取用户信息
// @Tags 用户
// @Success 200 {object} response.KResponse{data=model.UserResponse} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /users/get [get]
func (*userControll) GetUser(c *gin.Context) *response.KResponse {
	token := utils.GetToken(c)
	u := utils.GetUserID(c)

	user := dbUser.UserDB.GetUserByUserID(u)

	return response.MakeResponse(&model.UserResponse{
		User:  *user,
		Token: token,
	})
}

// @Summary  根据用户ID获取用户信息
// @Tags 用户
// @Param id path int true "用户ID"
// @Success 200 {object} response.KResponse{data=model.User} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /users/:id [get]
func (*userControll) GetUserByID(c *gin.Context) *response.KResponse {
	u, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	user := dbUser.UserDB.GetUserByUserID(model.UserID(u))

	return response.MakeResponse(user)
}

// @Summary 根据 search 关键字 获取用户列表
// @Tags 用户
// @Param search query string false "搜索关键字" default("")
// @Success 200 {object} response.KResponse{data=[]model.User} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /users/users [get]
func (*userControll) GetUsers(c *gin.Context) *response.KResponse {
	search := c.Query("search")

	users := dbUser.UserDB.GetUsers(search)

	if users == nil {
		return response.MakeResponseWidthCode("获取用户信息失败", http.StatusInternalServerError)
	}

	return response.MakeResponse(users)
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary 用户注册
// @Tags 用户
// @Param request body RegisterRequest true "注册请求"
// @Success 200 {object} response.KResponse{data=[]model.User} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /users/register [post]
func (*userControll) Register(c *gin.Context) *response.KResponse {
	user := &RegisterRequest{}

	if err := c.ShouldBindBodyWithJSON(user); err != nil {
		return response.MakeResponseWidthCode("无效的参数", http.StatusBadRequest)
	}

	if dbUser.UserDB.UserExists(user.Username) {
		return response.MakeResponseWidthCode("用户已存在", http.StatusBadRequest)
	}

	hashed := utils.GenerateFromPasswordString(user.Password)
	user.Password = hashed

	dbUser.UserDB.AddUser(&model.User{
		Username: user.Username,
		Password: user.Password,
	})

	user.Password = ""

	return response.MakeResponse(user)
}

// @Summary 用户登录
// @Tags 用户
// @Param request body RegisterRequest true "登录请求"
// @Success 200 {object} response.KResponse{data=[]model.UserResponse} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /users/login [post]
func (*userControll) Login(c *gin.Context) *response.KResponse {
	reqUser := &RegisterRequest{}

	if err := c.ShouldBindBodyWithJSON(reqUser); err != nil {
		return response.MakeResponseWidthCode("无效的参数", http.StatusBadRequest)
	}

	user, err := dbUser.UserDB.GetUserByUsername(reqUser.Username)
	if err != nil {
		return response.MakeResponseWidthCode("用户不存在", http.StatusBadRequest)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password))
	if err != nil {
		return response.MakeResponseWidthCode("密码错误", http.StatusBadRequest)
	}

	token, err := utils.GenerateToken(uint(user.ID))
	if err != nil {
		return response.MakeResponseWidthCode("生成token失败", http.StatusInternalServerError)
	}

	UserRedis.SetToken(c, token, user.ID)
	UserRedis.SetUserInfo(c, user.ID, user)

	return response.MakeResponse(&model.UserResponse{
		User:  *user,
		Token: token,
	})
}

// @Summary 退出登录
// @Tags 用户
// @Success 200 {object} response.KResponse{data=string} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /users/logout [get]
func (*userControll) Logout(c *gin.Context) *response.KResponse {
	token := utils.GetToken(c)
	err := UserRedis.DelToken(c, token)
	if err != nil {
		return response.MakeResponseWidthCode("退出失败", http.StatusInternalServerError)
	}

	return response.MakeResponse("退出成功")
}

type UserPutPasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// @Summary 修改用户密码;
// @Description 成功后；服务端将清除当前用户token的登录状态
// @Tags 用户
// @Param request body UserPutPasswordRequest true "修改密码请求"
// @Success 200 {object} response.KResponse{data=string} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /users/password [put]
func (*userControll) UserPutPassword(c *gin.Context) *response.KResponse {
	userID := utils.GetUserID(c)

	data := &UserPutPasswordRequest{}

	if err := c.ShouldBindBodyWithJSON(data); err != nil {
		return response.MakeResponseWidthCode("无效的参数", http.StatusBadRequest)
	}

	user := dbUser.UserDB.GetUserByUserID(userID)

	if !utils.CompareHashAndPassword(user.Password, data.OldPassword) {
		return response.MakeResponseWidthCode("旧密码错误", http.StatusBadRequest)
	}

	user.Password = utils.GenerateFromPasswordString(data.NewPassword)
	user.Nickname = "修改 passwowrd"

	dbUser.UserDB.UpdateUser(user)

	return response.MakeResponse(user)
}

type UserPutAvatarRequest struct {
	Avatar string `json:"avatar" binding:"required"`
}

// @Summary 用改用户头像;注意是存hash地址；而不是 url 地址;
// @Tags 用户
// @Param request body UserPutAvatarRequest true "修改头像请求"
// @Success 200 {object} response.KResponse{data=model.User} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /users/avatar [put]
func (*userControll) UserPutAvatar(c *gin.Context) *response.KResponse {
	userID := utils.GetUserID(c)
	data := utils.ShouldBindBodyWithJSON[*UserPutAvatarRequest](c)
	user := dbUser.UserDB.GetUserByUserID(userID)

	dbUser.UserDB.UpdateUserAvatar(userID, data.Avatar)

	user.Avatar = data.Avatar

	return response.MakeResponse(user)
}

type UserPutNicknameRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

// @Summary 用改用户中文名称;
// @Tags 用户
// @Param request body UserPutNicknameRequest true "修改中文名称请求"
// @Success 200 {object} response.KResponse{data=model.User} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /users/nickname [put]
func (*userControll) UserPutNickname(c *gin.Context) *response.KResponse {
	userID := utils.GetUserID(c)
	data := utils.ShouldBindBodyWithJSON[*UserPutNicknameRequest](c)
	user := dbUser.UserDB.GetUserByUserID(userID)

	dbUser.UserDB.UpdateUserNickname(userID, data.Nickname)

	user.Nickname = data.Nickname

	return response.MakeResponse(user)
}
