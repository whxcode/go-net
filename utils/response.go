package utils

import (
	"net/http"

	"go-net/model"

	"github.com/gin-gonic/gin"
)

type KResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Code    int    `json:"code"`
}

func MakeResponse(data any) *KResponse {
	return &KResponse{
		Data: data,
		Code: http.StatusOK,
	}
}

func MakeResponseWidthCode(data any, code int) *KResponse {
	return &KResponse{
		Data: data,
		Code: code,
	}
}

func GetUserID(c *gin.Context) model.UserID {
	userID, exists := c.Get("userID")

	if !exists {
		return model.InvalidUserID
	}

	id, ok := userID.(model.UserID)

	if !ok {
		return model.InvalidUserID
	}

	return id
}

func GetToken(c *gin.Context) string {
	userID, exists := c.Get("token")

	if !exists {
		return ""
	}

	id, ok := userID.(string)

	if !ok {
		return ""
	}

	return id
}

func SetToken(c *gin.Context, token string) {
	c.Set("token", token)
}

func SetUserID(c *gin.Context, userID model.UserID) {
	c.Set("userID", model.UserID(userID))
}
