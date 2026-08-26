package controller

import (
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

type groupController struct{}

// @Summary 获取群列表
// @Tags 群
// @Success 200 {object} utils.KResponse{data=model.UserResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /groups [get]
func (*groupController) Groups(c *gin.Context) *utils.KResponse {
	return utils.MakeResponse("获取成功")
}

var GroupController = &groupController{}
