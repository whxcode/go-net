package utils

import (
	"fmt"
	"strconv"

	"go-net/model"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) model.UserID {
	userID, exists := c.Get("userID")

	if !exists {
		fmt.Printf("userID is not of type model.UserID, got %T\n", userID)
		panic(fmt.Sprintf("userID is not of type model.UserID, got %T", userID))
	}

	id, ok := userID.(model.UserID)

	if !ok {
		fmt.Printf("userID is not of type model.UserID, got %T\n", userID)
		panic(fmt.Sprintf("userID is not of type model.UserID, got %T", userID))
	}

	return id
}

func GetToken(c *gin.Context) string {
	return c.Request.Header.Get("token")
}

func GetTokenFromQuery(c *gin.Context) string {
	return c.Query("token")
}

func SetUserID(c *gin.Context, userID model.UserID) {
	c.Set("userID", model.UserID(userID))
}

func ShouldBindBodyWithJSON[T any](c *gin.Context) T {
	var data T

	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		panic(fmt.Sprintf("Failed to bind JSON: %v", err))
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

func StringToUInt(s string) uint {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return uint(i)
}

/*
* 解析分页查询必备参数
*
* */
func ParsePageQuery(c *gin.Context) (limit int, offset int) {
	limitS := c.DefaultQuery("limit", "20")
	offsetS := c.DefaultQuery("offset", "0")

	return StringToInt(limitS), StringToInt(offsetS)
}
