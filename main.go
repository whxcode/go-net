package main

/**
*application/json	c.ShouldBindJSON(&req)	解析 JSON body
application/x-www-form-urlencoded	c.ShouldBind(&req) 或 c.PostForm("key")	表单格式
multipart/form-data	c.ShouldBind(&req) 或 c.FormFile("file")	文件上传
text/plain	c.GetRawData()	原始 body
XML	c.ShouldBindXML(&req)	XML body
GraphQL	c.GetRawData()	解析 query 字段
* */

// @title go-net
// @version 1.0
// @description 一个 IM 服务器
// @description 支持用户注册、登录、好友管理、群组聊天
// @description 基于 Go + Gin + GORM + Redis
// @description 目前访问服务器任何接口: localhost:8080/api 前端可用先使用代理设置
// @description 使用文件 hash 获取文件地址可用先固定: localhost:8080/api/file/:hash 前端可用看作做;方便后期迁移
// @host localhost:8080
// @BasePath /api
import (
	httpServer "go-net/http"
	"go-net/model"
	"go-net/redis"
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
