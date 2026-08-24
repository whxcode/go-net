package utils

import (
	"fmt"
	"net/http"
	"strconv"

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
		panic(fmt.Sprintf("userID is not of type model.UserID, got %T", userID))
		return model.InvalidUserID
	}

	id, ok := userID.(model.UserID)

	if !ok {
		panic(fmt.Sprintf("userID is not of type model.UserID, got %T", userID))
	}

	return id
}

func GetToken(c *gin.Context) string {
	return c.Request.Header.Get("token")
}

func SetUserID(c *gin.Context, userID model.UserID) {
	c.Set("userID", model.UserID(userID))
}

func ShouldBindBodyWithJSON[T any](c *gin.Context) T {
	var data T

	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		panic(err)
	}

	return data
}

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func StringToUserID(s string) model.UserID {
	i, err := strconv.Atoi(s)
	fmt.Println("StringToUserID:", s, i, err)
	if err != nil {
		panic(err)
	}

	return model.UserID(i)
}
