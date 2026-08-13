package main

import (
	httpServer "go-net/http"
	"go-net/model"
)

func main() {
	model.InitDB()
	model.InitRedis()

	httpServer.Start()
	// view.Test()
	// Create a Gin router with default middleware (logger and recovery)
}
