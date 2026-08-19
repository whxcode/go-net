package main

/**
*application/json	c.ShouldBindJSON(&req)	解析 JSON body
application/x-www-form-urlencoded	c.ShouldBind(&req) 或 c.PostForm("key")	表单格式
multipart/form-data	c.ShouldBind(&req) 或 c.FormFile("file")	文件上传
text/plain	c.GetRawData()	原始 body
XML	c.ShouldBindXML(&req)	XML body
GraphQL	c.GetRawData()	解析 query 字段
* */

import (
	httpServer "go-net/http"
	"go-net/model"
	"go-net/model/redis"
	"go-net/utils"
)

func main() {
	// 初始化雪花算法（节点ID 1，每台机器不同）
	if err := utils.InitSnowflake(1); err != nil {
		panic("雪花法初始化失败: " + err.Error())
	}

	model.InitDB()
	redis.InitRedis()

	// 	model.FriendDB.GetFirends(1)

	httpServer.Start()
	// view.Test()
	// Create a Gin router with default middleware (logger and recovery)
}
