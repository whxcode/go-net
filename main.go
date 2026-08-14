package main

import (
	httpServer "go-net/http"
	"go-net/model"
	"go-net/utils"
)

func main() {
	// 初始化雪花算法（节点ID 1，每台机器不同）
	if err := utils.InitSnowflake(1); err != nil {
		panic("雪花法初始化失败: " + err.Error())
	}

	model.InitDB()
	model.InitRedis()

	httpServer.Start()
	// view.Test()
	// Create a Gin router with default middleware (logger and recovery)
}
